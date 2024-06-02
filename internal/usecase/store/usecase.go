package store

import (
	"context"
	"fmt"

	appcommon "github.com/KartoonYoko/gophkeeper/internal/common"
	sfmodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/filestore"
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
	// - зашифровать данные
	// - сохранить шифрованные данные на файловом хранилище, получив ID
	// - сохранить общую информацию в БД
	// - в случае ошибки удалить запись в БД и данные на файловом хранилище
	encryptedData := uc.dataCipherHandler.Encrypt(request.Data)
	r := &sfmodel.SaveDataRequestModel{
		Data:   encryptedData,
		UserID: request.UserID,
	}
	sfr, err := uc.fstorage.SaveData(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	rsd := &smodel.SaveDataRequestModel{
		BinaryID:    sfr.ID,
		Description: request.Description,
		DataType:    smodel.DataType(request.DataType),
	}
	_, err = uc.storage.SaveData(ctx, rsd)
	if err != nil {
		err = fmt.Errorf("failed to save data: %w", err)
		// удаляем сохраненный файл
		removeRequest := &sfmodel.RemoveDataByIDRequestModel{
			UserID: request.UserID,
			ID:     sfr.ID,
		}
		removeDataErr := uc.fstorage.RemoveDataByID(ctx, removeRequest)
		if removeDataErr != nil {
			err = fmt.Errorf("failed delete not saved data: %w", err)
		}
		return err
	}
	return nil
}

func (uc *Usecase) GetDataByID(ctx context.Context, request *model.GetDataByIDRequestModel) (*model.GetDataByIDResponseModel, error) {
	// получаем метаинформацию о данных
	srequest:= &smodel.GetDataByIDRequestModel{
		UserID: request.UserID,
		ID:     request.ID,
	}
	sdr, err := uc.storage.GetDataByID(ctx, srequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	// получаем сами данные
	frequest := &sfmodel.GetDataByIDRequestModel{
		UserID: request.UserID,
		ID:     request.ID,
	}
	gdr, err := uc.fstorage.GetDataByID(ctx, frequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	// расшифровываем данные
	encryptedData, err := uc.dataCipherHandler.Decrypt(gdr.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	
	response := &model.GetDataByIDResponseModel{
		Data: encryptedData,
		Description: sdr.Description,
		DataType:    sdr.DataType,
	}

	return response, nil
}
