package clientstorage

import (
	"context"
	"encoding/json"
	"errors"
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
		if errors.Is(err, os.ErrNotExist) {
			return ``, ``, ErrTokensNotFound
		}
		return ``, ``, err
	}

	tf := &credentialsFile{at, rt, ``, ``}
	err = json.Unmarshal(b, tf)
	if err != nil {
		return ``, ``, err
	}

	return tf.AccessToken, tf.RefreshToken, nil
}

func (s *Storage) UpdateTokens(accesstoken string, refreshtoken string) (err error) {
	// считаем текущие настройки
	b, err := os.ReadFile(s.getTokensPath())
	if err != nil {
		return err
	}
	tf := &credentialsFile{}
	err = json.Unmarshal(b, tf)
	if err != nil {
		return err
	}

	// обновим токены
	tf.AccessToken = accesstoken
	tf.RefreshToken = refreshtoken

	// запишем настройки
	b, err = json.Marshal(tf)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.getTokensPath(), b, 0600)
	if err != nil {
		return err
	}

	return nil
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
	err := os.WriteFile(s.getDataPathWithName(request.ID), request.Data, os.ModePerm)
	if err != nil {
		return err
	}

	query := `INSERT INTO data_store (id, user_id, description, data_type, hash, modification_timestamp) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query,
		request.ID,
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
	rows, err := s.db.QueryContext(ctx, `SELECT id, user_id, description, data_type FROM data_store WHERE user_id = ? AND is_deleted = 0`, userID)
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

func (s *Storage) GetDataByID(ctx context.Context, id string) (*GetDataByIDResponseModel, error) {
	b, err := os.ReadFile(s.getDataPathWithName(id))
	if err != nil {
		return nil, err
	}

	res := &GetDataByIDResponseModel{
		Data: b,
	}

	query := `SELECT id, user_id, description, data_type, hash, modification_timestamp FROM data_store WHERE id = ?`
	err = s.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID,
		&res.Userid,
		&res.Description,
		&res.Datatype,
		&res.Hash,
		&res.ModificationTimestamp)

	if err != nil {
		return nil, err
	}

	return res, nil
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

func (s *Storage) UpdateData(ctx context.Context, request UpdateDataRequestModel) error {
	// - создать новый файл с названием: ID + суффикс нового файла
	// - удалить старый файл
	// - переименовать новый
	// - обновить данные в БД
	err := os.WriteFile(s.getNewDataPathWithName(request.ID), request.Data, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Remove(s.getDataPathWithName(request.ID))
	if err != nil {
		return err
	}
	err = os.Rename(s.getNewDataPathWithName(request.ID), s.getDataPathWithName(request.ID))
	if err != nil {
		return err
	}

	query := `UPDATE data_store SET hash=?, modification_timestamp=? WHERE id = ?`
	_, err = s.db.ExecContext(ctx, query, request.Hash, request.ModificationTimestamp, request.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveDataByID(ctx context.Context, request RemoveDataByIDRequestModel) error {
	query := `UPDATE data_store SET is_deleted=1, modification_timestamp=? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, request.ModificationTimestamp, request.DataID)
	if err != nil {
		return err
	}

	return os.Remove(s.getDataPathWithName(request.DataID))
}

func (s *Storage) getTokensPath() string {
	return s.rootPath + string(os.PathSeparator) + s.tokensFileName
}

func (s *Storage) getDBPath() string {
	return s.dbName
}

// getDataPathWithName возвращает путь до файла с указанным именем
func (s *Storage) getDataPathWithName(filename string) string {
	return s.rootPath + string(os.PathSeparator) + filename
}

// getNewDataPathWithName возвращает путь до файла предназначенного для замены старого с указанным именем
func (s *Storage) getNewDataPathWithName(filename string) string {
	return s.rootPath + string(os.PathSeparator) + filename + "_new"
}
