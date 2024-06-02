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

	dataCipherHandler *appcommon.DataCipherHandler
}

func New(conf Config, storage Storager, fstorager FileStorager) (*Usecase, error) {
	var err error
	uc := new(Usecase)

	uc.storage = storage
	uc.fstorage = fstorager
	uc.conf = conf
	uc.dataCipherHandler, err = appcommon.NewDataCipherHandler(conf.DataSecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret key handler: %w", err)
	}
	
	return uc, nil
}

func (uc *Usecase) SaveData(ctx context.Context, request *model.SaveDataRequestModel) error {
	encryptedData := uc.dataCipherHandler.Encrypt(request.Data)
	r := bytes.NewReader(encryptedData)
	sfr, err := uc.fstorage.SaveData(ctx, r)
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
		err = fmt.Errorf("failed to save data: %w", err)
		// удаляем сохраненный файл
		removeDataErr := uc.fstorage.RemoveDataByID(ctx, sfr.ID)
		if removeDataErr != nil {
			err = fmt.Errorf("failed delete not save data: %w", err)
		}
		return err
	}
	return nil
}

func (uc *Usecase) GetDataByID(ctx context.Context, request *model.GetDataByIDRequestModel) (*model.GetDataByIDResponseModel, error) {
	uc.fstorage.
}
