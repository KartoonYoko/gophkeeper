package store

import (
	"context"

	filestoremodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/filestore"
	storemodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/store"
)

type Storager interface {
	SaveData(ctx context.Context, request *storemodel.SaveDataRequestModel) (*storemodel.SaveDataResponseModel, error)
	UpdateData(ctx context.Context, request *storemodel.UpdateDataRequestModel) (*storemodel.UpdateDataResponseModel, error)
	GetDataByID(ctx context.Context, request *storemodel.GetDataByIDRequestModel) (*storemodel.GetDataByIDResponseModel, error)
	GetUserDataList(ctx context.Context, request *storemodel.GetUserDataListRequestModel) (*storemodel.GetUserDataListResponseModel, error)
	RemoveDataByID(ctx context.Context, request *storemodel.RemoveDataByIDRequestModel) error
}

type FileStorager interface {
	SaveData(ctx context.Context, request *filestoremodel.SaveDataRequestModel) (*filestoremodel.SaveDataResponseModel, error)
	UpdateData(ctx context.Context, request *filestoremodel.UpdateDataRequestModel) (*filestoremodel.UpdateDataResponseModel, error)
	GetDataByID(ctx context.Context, request *filestoremodel.GetDataByIDRequestModel) (*filestoremodel.GetDataByIDResponseModel, error)
	RemoveDataByID(ctx context.Context, request *filestoremodel.RemoveDataByIDRequestModel) error
}

type DataCipherHandler interface {
	Encrypt(data []byte) []byte
	Decrypt(data []byte) (encryptedname []byte, err error)
}
