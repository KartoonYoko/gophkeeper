package store

import (
	"context"
	"io"

	filestoremodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/filestore"
	storemodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/store"
)

type Storager interface {
	GetSecretKeyByUserID(ctx context.Context, userID string) (string, error)
	SaveData(ctx context.Context, request storemodel.SaveDataRequestModel) (storemodel.SaveDataResponseModel, error)
}

type FileStorager interface {
	SaveFile(ctx context.Context, reader io.Reader) (filestoremodel.FileSaveRsponseModel, error)
}
