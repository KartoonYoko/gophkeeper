package grpcserver

import (
	"context"

	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/store"
)

type StoreUsecase interface {
	SaveData(ctx context.Context, request *model.SaveDataRequestModel) (*model.SaveDataResponseModel, error)
	GetDataByID(ctx context.Context, request *model.GetDataByIDRequestModel) (*model.GetDataByIDResponseModel, error)
}
