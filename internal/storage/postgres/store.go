package postgres

import (
	"context"
	"fmt"

	model "github.com/KartoonYoko/gophkeeper/internal/storage/model/store"
)

func (s *Storage) SaveData(ctx context.Context, request *model.SaveDataRequestModel) (*model.SaveDataResponseModel, error) {
	query := `
	INSERT INTO "store"."data" (user_id, binary_id, description, data_type, id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING "id";
	`

	var id string
	err := s.pool.QueryRow(ctx, query, 
		request.UserID, 
		request.BinaryID, 
		request.Description, 
		request.DataType,
		request.ID).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to save data: %w", err)
	}

	response := new(model.SaveDataResponseModel)
	response.ID = id

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
