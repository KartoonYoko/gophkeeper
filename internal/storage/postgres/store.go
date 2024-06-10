package postgres

import (
	"context"
	"fmt"

	model "github.com/KartoonYoko/gophkeeper/internal/storage/model/store"
)

func (s *Storage) SaveData(ctx context.Context, request *model.SaveDataRequestModel) (*model.SaveDataResponseModel, error) {
	query := `
	INSERT INTO "store"."data" (user_id, binary_id, description, data_type, id, hash, modification_timestamp)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING "id";
	`

	var id string
	err := s.pool.QueryRow(ctx, query, 
		request.UserID, 
		request.BinaryID, 
		request.Description, 
		request.DataType,
		request.ID,
		request.Hash,
		request.ModificationTimestamp).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to save data: %w", err)
	}

	response := new(model.SaveDataResponseModel)
	response.ID = id

	return response, nil
}

func (s *Storage) UpdateData(ctx context.Context, request *model.UpdateDataRequestModel) (*model.UpdateDataResponseModel, error) {
	query := `
	UPDATE "store"."data" 
	SET "hash" = $1, "modification_timestamp" = $2 
	WHERE "id" = $3 AND "user_id" = $4`

	_, err := s.pool.Exec(ctx, query, request.Hash, request.ModificationTimestamp, request.ID, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to update data: %w", err)
	}

	response := new(model.UpdateDataResponseModel)

	return response, nil
}

func (s *Storage) GetDataByID(ctx context.Context, request *model.GetDataByIDRequestModel) (*model.GetDataByIDResponseModel, error) {
	var id, userID, binaryID, description, dataType string

	query := `
	SELECT * FROM "store"."data" WHERE "id" = $1 AND "user_id" = $2
	`

	err := s.pool.
		QueryRow(ctx, query, request.ID, request.UserID).
		Scan(&id, &userID, &binaryID, &description, &dataType)

	if err != nil {
		return nil, fmt.Errorf("unable get data by id: %w", err)
	}

	response := new(model.GetDataByIDResponseModel)
	response.ID = id
	response.BinaryID = binaryID
	response.DataType = dataType
	response.Description = description

	return response, nil
}

func (s *Storage) RemoveDataByID(ctx context.Context, request *model.RemoveDataByIDRequestModel) error {
	query := `
	DELETE FROM "store"."data" WHERE "id" = $1 AND "user_id" = $2
	`

	_, err := s.pool.Exec(ctx, query, request.ID, request.UserID)
	if err != nil {
		return fmt.Errorf("failed to remove data: %w", err)
	}

	return nil
}

func (s *Storage) GetUserDataList(ctx context.Context, request *model.GetUserDataListRequestModel) (*model.GetUserDataListResponseModel, error) {
	query := `
	SELECT id, description, data_type, hash, modification_timestamp, is_deleted 
	FROM "store"."data" 
	WHERE "user_id" = $1
	`

	rows, err := s.pool.Query(ctx, query, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to query to get user data list: %w", err)
	}
	defer rows.Close()

	response := new(model.GetUserDataListResponseModel)
	for rows.Next() {
		item := new(model.GetUserDataListResponseItemModel)

		rows.Scan(
			&item.ID, 
			&item.Description, 
			&item.DataType, 
			&item.Hash, 
			&item.ModificationTimestamp, 
			&item.IsDeleted)

			response.Items = append(response.Items, item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to get user data list: %w", err)
	}

	return response, nil
}
