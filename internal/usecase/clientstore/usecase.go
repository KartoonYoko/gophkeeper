package clientstore

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/KartoonYoko/gophkeeper/internal/common"
	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
	"github.com/KartoonYoko/gophkeeper/internal/storage/clientstorage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Usecase struct {
	client pb.StoreServiceClient

	storage *clientstorage.Storage
}

func New(conn *grpc.ClientConn, store *clientstorage.Storage) *Usecase {
	uc := new(Usecase)

	uc.client = pb.NewStoreServiceClient(conn)
	uc.storage = store

	return uc
}

// CreateTextData сохраняет текстовые данные
func (uc *Usecase) CreateTextData(ctx context.Context, text string) error {
	// - шифрую данные и работаю далее только с шифрованными данными
	// - сохраняю локльно на диске (генерировать название файла, используя guid)
	// - сохраняю информацию в БД
	// - пробую сохранить удалённо

	// формируем данные для сохранения
	cipher, err := uc.getDataCipher()
	if err != nil {
		return err
	}
	id := uuid.New().String()
	encrypted := cipher.Encrypt([]byte(text))
	hash := uc.getDataHash([]byte(text))
	modts := uc.getModificationTimestamp()
	userID, err := uc.storage.GetUserID()
	if err != nil {
		return err
	}

	// пробуем сохранить локально
	r := clientstorage.SaveDataRequestModel{
		ID:                    id,
		Userid:                userID,
		Description:           "",
		Datatype:              "TEXT",
		Hash:                  hash,
		ModificationTimestamp: modts,
		Data:                  encrypted,
	}
	err = uc.storage.SaveData(ctx, r)
	if err != nil {
		return err
	}

	// сохраняем удалённо
	rr := &pb.SaveDataRequest{
		Id:                    r.ID,
		Description:           r.Description,
		Type:                  getPBDataTypeFromString(r.Datatype),
		Data:                  r.Data,
		Hash:                  r.Hash,
		ModificationTimestamp: modts,
	}
	_, err = uc.client.SaveData(ctx, rr)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) CreateBinaryData(ctx context.Context, filepath string) error {
	// TODO
	return nil
}

func (uc *Usecase) CreateBankCardData(ctx context.Context, card BankCardDataModel) error {
	// TODO
	return nil
}

func (uc *Usecase) CreateCredentialsData(ctx context.Context, credential *CredentialDataModel) error {
	// TODO
	return nil
}

func (uc *Usecase) GetDataList(ctx context.Context) ([]clientstorage.GetDataListResponseItemModel, error) {
	userid, err := uc.storage.GetUserID()
	if err != nil {
		return nil, err
	}

	return uc.storage.GetDataList(ctx, userid)
}

// Synchronize синхронизирует данные с сервером
func (uc *Usecase) Synchronize(ctx context.Context) error {
	// TODO
	// - список ID'шников, которые нужно добавить на сервер (1)
	// - список ID'шников, которые нужно обновить на сервере (2)
	// - список ID'шников, которые нужно удалить на сервере (3)
	// - список ID'шников, которые нужно добавить локльно (4)
	// - список ID'шников, которые нужно обновить локально (5)
	// - список ID'шников, которые нужно удалить локльно (6)
	//
	// - прохожусь по данным
	// 	- заполняю списки (если ID'шник попал в список, то для других списков его можно не смотреть):
	// 		- запоминаю те, которые помечены удалёнными на клиенте, но не на сервере (3)
	// 		- запоминаю те, которые помечены удалёнными на сервере, но не на клиенте (6)
	// 		- запоминаю те, которых нет на сервере (1)
	// 		- запоминаю те, которые есть на сервере, но не на клиенте (4)
	// 		- запоминаю те, у кого не совпал хеш и дата модификации меньше, чем локально (2)
	// 		- запоминаю те, у кого не совпал хеш и дата модификации больше, чем локльно (5)
	//
	// - обрабатываю полученные списки

	userID, err := uc.storage.GetUserID()
	if err != nil {
		return err
	}

	localDataDict := make(map[string]*clientstorage.GetDataListToSynchronizeItemModel)
	remoteDataDict := make(map[string]*pb.GetMetaDataListItemResponse)

	// get remote metadata list
	remoteLst, err := uc.client.GetMetaDataList(ctx, &pb.GetMetaDataListRequest{})
	if err != nil {
		return err
	}

	// get local metadata list
	localLst, err := uc.storage.GetDataListToSynchronize(ctx, userID)
	if err != nil {
		return err
	}

	for _, item := range localLst {
		localDataDict[item.ID] = &item
	}

	for _, item := range remoteLst.Items {
		remoteDataDict[item.Id] = item
	}

	idsToDeleteLocal := make([]string, 0)
	idsToDeleteRemote := make([]string, 0)
	idsToAddLocal := make([]string, 0)
	idsToAddRemote := make([]string, 0)
	idsToUpdateLocal := make([]string, 0)
	idsToUpdateRemote := make([]string, 0)

	for _, item := range localLst {
		// запоминаю те, которые помечены удалёнными на клиенте, но не на сервере (3)
		if item.IsDeleted {
			if ritem, ok := remoteDataDict[item.ID]; ok {
				if !ritem.IsDeleted {
					idsToDeleteRemote = append(idsToDeleteRemote, item.ID)
				}
			}

			continue
		}

		// запоминаю те, которых нет на сервере (1)
		ritem, ok := remoteDataDict[item.ID]
		if !ok {
			idsToAddRemote = append(idsToAddRemote, item.ID)
			continue
		}

		// запоминаю те, у кого не совпал хеш и дата модификации больше, чем локльно (5)
		if ritem.Hash != item.Hash {
			if ritem.ModificationTimestamp >= item.ModificationTimestamp {
				idsToUpdateLocal = append(idsToUpdateLocal, item.ID)
			}
		}
	}

	for _, ritem := range remoteLst.Items {
		// запоминаю те, которые помечены удалёнными на сервере, но не на клиенте (6)
		if ritem.IsDeleted {
			if litem, ok := localDataDict[ritem.Id]; ok {
				if !litem.IsDeleted {
					idsToDeleteLocal = append(idsToDeleteLocal, litem.ID)
				}
			}
			continue
		}
		// запоминаю те, которые есть на сервере, но не на клиенте (4)
		item, ok := localDataDict[ritem.Id]
		if !ok {
			idsToAddLocal = append(idsToAddLocal, ritem.Id)
			continue
		}
		// запоминаю те, у кого не совпал хеш и дата модификации меньше, чем локально (2)
		if ritem.Hash != item.Hash {
			if ritem.ModificationTimestamp < item.ModificationTimestamp {
				idsToUpdateRemote = append(idsToUpdateRemote, ritem.Id)
			}
		}
	}

	// обрабатываю полученные списки
	// TODO написать функции обновления/удаления данных

	return nil
}

func (uc *Usecase) getDataCipher() (*common.DataCipherHandler, error) {
	sc, err := uc.storage.GetSecretKey()
	if err != nil {
		return nil, err
	}
	cipher, err := common.NewDataCipherHandler(sc)
	if err != nil {
		return nil, err
	}

	return cipher, nil
}

func (uc *Usecase) getModificationTimestamp() int64 {
	return time.Now().UTC().UnixMilli()
}

// getDataHash возвращает хеш данных. Применять к нешифрованным данным.
func (uc *Usecase) getDataHash(data []byte) string {
	hash := common.NewDataHasherSHA256().Hash(data)
	return base64.StdEncoding.EncodeToString(hash)
}

// checkDataHash сверяет хеш данных. Применять к нешифрованным данным.
func (uc *Usecase) checkDataHash(data []byte, datahash string) bool {
	return uc.getDataHash(data) == datahash
}

func getPBDataTypeFromString(str string) pb.DataTypeEnum {
	s := strings.ToUpper(str)
	if s == "CREDENTIALS" {
		return pb.DataTypeEnum_DATA_TYPE_CREDENTIALS
	} else if s == "TEXT" {
		return pb.DataTypeEnum_DATA_TYPE_TEXT
	} else if s == "BINARY" {
		return pb.DataTypeEnum_DATA_TYPE_BINARY
	} else if s == "BANK_CARD" {
		return pb.DataTypeEnum_DATA_TYPE_BANK_CARD
	}

	return pb.DataTypeEnum_DATA_TYPE_TEXT
}
