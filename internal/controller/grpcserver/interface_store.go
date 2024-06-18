package grpcserver

import (
	"context"

	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/store"
	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/store"
)

// StoreUsecase usecase для ручек хранилища
type StoreUsecase interface {
	SaveData(ctx context.Context, request *model.SaveDataRequestModel) (*model.SaveDataResponseModel, error)
	UpdateData(ctx context.Context, request *model.UpdateDataRequestModel) (*model.UpdateDataResponseModel, error)
	GetDataByID(ctx context.Context, request *model.GetDataByIDRequestModel) (*model.GetDataByIDResponseModel, error)
	GetUserDataList(ctx context.Context, userID string) (*smodel.GetUserDataListResponseModel, error)
	RemoveDataByID(ctx context.Context, request *model.RemoveDataByIDRequestModel) (*model.RemoveDataByIDResponseModel, error)
}
