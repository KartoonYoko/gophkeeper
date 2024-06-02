package store

import (
	"context"
	"io"

	filestoremodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/filestore"
	storemodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/store"
)

type Storager interface {
	SaveData(ctx context.Context, request storemodel.SaveDataRequestModel) (*storemodel.SaveDataResponseModel, error)
}

type FileStorager interface {
	SaveData(ctx context.Context, reader io.Reader) (filestoremodel.SaveDataResponseModel, error)
	GetDataByID(ctx context.Context, request filestoremodel.GetDataByIDRequestModel) (*filestoremodel.GetDataByIDResponseModel, error)
	RemoveDataByID(ctx context.Context, request filestoremodel.RemoveDataByIDRequestModel) error
}
