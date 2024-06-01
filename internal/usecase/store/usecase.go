package store

import (
	"bytes"
	"context"
	"fmt"

	appcommon "github.com/KartoonYoko/gophkeeper/internal/common"
	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/store"
	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/store"
)

type Usecase struct {
	conf Config

	storage  Storager
	fstorage FileStorager

	secretkeyHandler *appcommon.SecretKeyHandler
}

func New(conf Config, storage Storager, fstorager FileStorager) (*Usecase, error) {
	var err error
	uc := new(Usecase)

	uc.storage = storage
	uc.fstorage = fstorager
	uc.conf = conf
	uc.secretkeyHandler, err = appcommon.NewSecretKeyHandler(conf.SecretKeySecure)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret key handler: %w", err)
	}
	
	return uc, nil
}

func (uc *Usecase) SaveData(ctx context.Context, request *model.SaveDataRequestModel) error {
	// todo зашифровать данные
	sc, err := uc.storage.GetSecretKeyByUserID(ctx, request.UserID)
	if err != nil {
		return fmt.Errorf("failed to get secret key: %w", err)
	}
	r := bytes.NewReader(request.Data)
	sfr, err := uc.fstorage.SaveFile(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	rsd := smodel.SaveDataRequestModel{
		BinaryID:    sfr.ID,
		Description: request.Description,
		DataType:    smodel.DataType(request.DataType),
	}
	_, err = uc.storage.SaveData(ctx, rsd)
	if err != nil {
		// todo удалить сохраненный файл
		return fmt.Errorf("failed to save data: %w", err)
	}
	return nil
}
