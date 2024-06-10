package grpcserver

import (
	"context"

	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/store"
)

type StoreUsecase interface {
	SaveData(ctx context.Context, request *model.SaveDataRequestModel) (*model.SaveDataResponseModel, error)
	UpdateData(ctx context.Context, request *model.UpdateDataRequestModel) (*model.UpdateDataResponseModel, error)
	GetDataByID(ctx context.Context, request *model.GetDataByIDRequestModel) (*model.GetDataByIDResponseModel, error)
	RemoveDataByID(ctx context.Context, request *model.RemoveDataByIDRequestModel) (*model.RemoveDataByIDResponseModel, error)
}
