package store

import (
	"context"
	"fmt"

	sfmodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/filestore"
	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/store"
	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/store"
)

type Usecase struct {
	conf Config

	// хранилище метаданных
	Storage Storager
	// хранилище данных
	FileStorage FileStorager

	dataCipherHandler DataCipherHandler
}

func New(conf Config, storage Storager, fstorager FileStorager, dataCipherHandler DataCipherHandler) *Usecase {
	uc := new(Usecase)

	uc.Storage = storage
	uc.FileStorage = fstorager
	uc.conf = conf
	uc.dataCipherHandler = dataCipherHandler
	return uc
}

func (uc *Usecase) SaveData(ctx context.Context, request *model.SaveDataRequestModel) (*model.SaveDataResponseModel, error) {
	// - зашифровать данные
	// - сохранить шифрованные данные на файловом хранилище, получив ID
	// - сохранить общую информацию в БД
	// - в случае ошибки удалить запись в БД и данные на файловом хранилище
	if !request.DataType.IsValid() {
		return nil, fmt.Errorf("invalid data type %s", request.DataType)
	}

	encryptedData := uc.dataCipherHandler.Encrypt(request.Data)
	r := &sfmodel.SaveDataRequestModel{
		ID:     request.ID,
		Data:   encryptedData,
		UserID: request.UserID,
	}
	sfr, err := uc.FileStorage.SaveData(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	rsd := &smodel.SaveDataRequestModel{
		ID:                    sfr.ID,
		BinaryID:              sfr.ID,
		UserID:                request.UserID,
		Description:           request.Description,
		DataType:              request.DataType.String(),
		Hash:                  request.Hash,
		ModificationTimestamp: request.ModificationTimestamp,
	}
	resSaveData, err := uc.Storage.SaveData(ctx, rsd)
	if err != nil {
		err = fmt.Errorf("failed to save data: %w", err)
		// удаляем сохраненный файл
		removeRequest := &sfmodel.RemoveDataByIDRequestModel{
			UserID: request.UserID,
			ID:     sfr.ID,
		}
		removeDataErr := uc.FileStorage.RemoveDataByID(ctx, removeRequest)
		if removeDataErr != nil {
			err = fmt.Errorf("failed delete not saved data: %w", err)
		}
		return nil, err
	}

	response := &model.SaveDataResponseModel{
		DataID: resSaveData.ID,
	}
	return response, nil
}

func (uc *Usecase) GetDataByID(ctx context.Context, request *model.GetDataByIDRequestModel) (*model.GetDataByIDResponseModel, error) {
	// получаем метаинформацию о данных
	srequest := &smodel.GetDataByIDRequestModel{
		UserID: request.UserID,
		ID:     request.ID,
	}
	sdr, err := uc.Storage.GetDataByID(ctx, srequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	// получаем сами данные
	frequest := &sfmodel.GetDataByIDRequestModel{
		UserID: request.UserID,
		ID:     sdr.BinaryID,
	}
	gdr, err := uc.FileStorage.GetDataByID(ctx, frequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	// расшифровываем данные
	encryptedData, err := uc.dataCipherHandler.Decrypt(gdr.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	response := &model.GetDataByIDResponseModel{
		Data:        encryptedData,
		Description: sdr.Description,
		DataType:    sdr.DataType,
	}

	return response, nil
}

func (uc *Usecase) UpdateData(ctx context.Context, request *model.UpdateDataRequestModel) (*model.UpdateDataResponseModel, error) {
	encryptedData := uc.dataCipherHandler.Encrypt(request.Data)

	r := &sfmodel.SaveDataRequestModel{
		ID:     request.ID,
		Data:   encryptedData,
		UserID: request.UserID,
	}
	_, err := uc.FileStorage.SaveData(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// сохранить в БД
	rs := &smodel.UpdateDataRequestModel{
		ID:                    request.ID,
		UserID:                request.UserID,
		ModificationTimestamp: request.ModificationTimestamp,
		Hash:                  request.Hash,
	}
	_, err = uc.Storage.UpdateData(ctx, rs)
	if err != nil {
		return nil, fmt.Errorf("failed to update data: %w", err)
	}

	response := new(model.UpdateDataResponseModel)

	return response, nil
}

func (uc *Usecase) RemoveDataByID(ctx context.Context, request *model.RemoveDataByIDRequestModel) (*model.RemoveDataByIDResponseModel, error) {
	err := uc.Storage.RemoveDataByID(ctx, &smodel.RemoveDataByIDRequestModel{
		ID:                    request.ID,
		UserID:                request.UserID,
		ModificationTimestamp: request.ModificationTimestamp,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove data from database: %w", err)
	}
	err = uc.FileStorage.RemoveDataByID(ctx, &sfmodel.RemoveDataByIDRequestModel{
		ID:     request.ID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove data from storage: %w", err)
	}

	return new(model.RemoveDataByIDResponseModel), nil
}

func (uc *Usecase) GetUserDataList(ctx context.Context, userID string) (*smodel.GetUserDataListResponseModel, error) {
	return uc.Storage.GetUserDataList(ctx, &smodel.GetUserDataListRequestModel{
		UserID: userID,
	})
}
