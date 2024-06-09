package clientstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"database/sql"

	_ "modernc.org/sqlite"
)

type Storage struct {
	dbName string
	db     *sql.DB

	rootPath string

	tokensFileName string
}

func New(ctx context.Context) (*Storage, error) {
	s := new(Storage)

	// инициализируем директорию для хранения данных
	hd, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable get user home directory: %w", err)
	}
	rootPath := hd + string(os.PathSeparator) + ".gophkeeper"
	err = os.MkdirAll(rootPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("unable create root directory: %w", err)
	}

	// инициализируем БД
	dbName := "data.s3db"
	s.dbName = dbName
	db, err := sql.Open("sqlite", s.getDBPath())
	if err != nil {
		return nil, err
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable ping database: %w", err)
	}
	err = migrate(db)
	if err != nil {
		return nil, fmt.Errorf("unable migrate database: %w", err)
	}

	s.db = db
	s.dbName = dbName
	s.rootPath = rootPath
	s.tokensFileName = "tokens.json"

	return s, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveCredentials(ctx context.Context, at string, rt string, sk string, userID string) error {
	b, err := json.Marshal(credentialsFile{at, rt, sk, userID})
	if err != nil {
		return err
	}
	err = os.WriteFile(s.getTokensPath(), b, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetTokens() (at string, rt string, err error) {
	b, err := os.ReadFile(s.getTokensPath())
	if err != nil {
		return ``, ``, err
	}

	tf := &credentialsFile{at, rt, ``, ``}
	err = json.Unmarshal(b, tf)
	if err != nil {
		return ``, ``, err
	}

	return tf.AccessToken, tf.RefreshToken, nil
}

func (s *Storage) RemoveTokens() error {
	return os.Remove(s.getTokensPath())
}

func (s *Storage) GetSecretKey() (sc string, err error) {
	b, err := os.ReadFile(s.getTokensPath())
	if err != nil {
		return ``, err
	}

	tf := &credentialsFile{``, ``, sc, ``}
	err = json.Unmarshal(b, tf)
	if err != nil {
		return ``, err
	}

	return tf.SecretKey, nil
}

func (s *Storage) GetUserID() (userID string, err error) {
	b, err := os.ReadFile(s.getTokensPath())
	if err != nil {
		return ``, err
	}

	tf := &credentialsFile{``, ``, ``, userID}
	err = json.Unmarshal(b, tf)
	if err != nil {
		return ``, err
	}

	return tf.UserID, nil
}

func (s *Storage) SaveData(ctx context.Context, request SaveDataRequestModel) error {
	err := os.WriteFile(s.getDataPathWithName(request.Filename), request.Data, os.ModePerm)
	if err != nil {
		return err
	}

	query := `INSERT INTO data_store (id, user_id, description, data_type, hash, modification_timestamp) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query,
		request.Filename,
		request.Userid,
		request.Description,
		request.Datatype,
		request.Hash,
		request.ModificationTimestamp)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetDataList(ctx context.Context, userID string) ([]GetDataListResponseItemModel, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, user_id, description, data_type FROM data_store WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []GetDataListResponseItemModel{}
	for rows.Next() {
		var item GetDataListResponseItemModel
		err = rows.Scan(&item.ID, &item.UserID, &item.Description, &item.Datatype)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return items, err
	}
	return items, nil
}

func (s *Storage) GetDataListToSynchronize(ctx context.Context, userID string) ([]GetDataListToSynchronizeItemModel, error) {
	query := `
	SELECT id, user_id, description, data_type, hash, modification_timestamp, is_deleted
	FROM data_store 
	WHERE user_id = ?`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []GetDataListToSynchronizeItemModel{}
	for rows.Next() {
		var item GetDataListToSynchronizeItemModel
		err = rows.Scan(&item.ID, &item.UserID, &item.Description, &item.Datatype, &item.Hash, &item.ModificationTimestamp, &item.IsDeleted)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return items, err
	}
	return items, nil
}

func (s *Storage) getTokensPath() string {
	return s.rootPath + string(os.PathSeparator) + s.tokensFileName
}

func (s *Storage) getDBPath() string {
	return s.dbName
}

func (s *Storage) getDataPathWithName(filename string) string {
	return s.rootPath + string(os.PathSeparator) + filename
}
