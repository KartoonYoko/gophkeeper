package clientstore

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/KartoonYoko/gophkeeper/internal/common"
	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
	"github.com/KartoonYoko/gophkeeper/internal/storage/clientstorage"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	uccommon "github.com/KartoonYoko/gophkeeper/internal/usecase/common/cliclient"
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
//
// TODO: обернуть обращение к серверу, чтобы выводить в этом случае не ошибку, а предупреждение, что сервер не досутпен, но данные сохранены локально
func (uc *Usecase) CreateTextData(ctx context.Context, text string, description string) error {
	return uc.saveBytes(ctx, []byte(text), description, "TEXT")
}

func (uc *Usecase) CreateBinaryData(ctx context.Context, filepath string, description string) error {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return uc.saveBytes(ctx, b, description, "BINARY")
}

func (uc *Usecase) CreateBankCardData(ctx context.Context, card BankCardDataModel, description string) error {
	b, err := json.Marshal(card)
	if err != nil {
		return err
	}

	return uc.saveBytes(ctx, b, description, "BANK_CARD")
}

func (uc *Usecase) CreateCredentialsData(ctx context.Context, credential CredentialDataModel, description string) error {
	b, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	return uc.saveBytes(ctx, b, description, "CREDENTIALS")
}

func (uc *Usecase) GetDataList(ctx context.Context) ([]clientstorage.GetDataListResponseItemModel, error) {
	userid, err := uc.storage.GetUserID()
	if err != nil {
		return nil, err
	}

	return uc.storage.GetDataList(ctx, userid)
}

func (uc *Usecase) GetDataByID(ctx context.Context, id string) (*clientstorage.GetDataByIDResponseModel, error) {
	res, err := uc.storage.GetDataByID(ctx, id)
	if err != nil {
		return nil, err
	}

	cipher, err := uc.getDataCipher()
	if err != nil {
		return nil, err
	}

	res.Data, err = cipher.Decrypt(res.Data)
	if err != nil {
		return nil, err
	}

	// TODO подумать как обходиться с разными типами данных;
	// скорее всего жто должен делать контроллер, вызывая для разных типов данных разные методы
	return res, nil
}

// Synchronize синхронизирует данные с сервером
func (uc *Usecase) Synchronize(ctx context.Context) error {
	//
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
		return uccommon.NewServerError(err)
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

	idsToDeleteLocal := make([]string, 0)  // удалить локально
	idsToDeleteRemote := make([]string, 0) // удалить на сервере
	idsToAddLocal := make([]string, 0)     // добавить локально
	idsToAddRemote := make([]string, 0)    // добавить на сервере
	idsToUpdateLocal := make([]string, 0)  // обновить локально
	idsToUpdateRemote := make([]string, 0) // обновить на сервере

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

	// обработка полученных списков
	for _, id := range idsToDeleteLocal {
		err = uc.storage.RemoveDataByID(ctx, clientstorage.RemoveDataByIDRequestModel{
			DataID:                id,
			ModificationTimestamp: uc.getModificationTimestamp(),
		})
		if err != nil {
			return err
		}
	}
	for _, id := range idsToDeleteRemote {
		_, err = uc.client.RemoveData(ctx, &pb.RemoveDataRequest{Id: id})
		if err != nil {
			return uccommon.NewServerError(err)
		}
	}
	for _, id := range idsToAddLocal {
		// get data by id
		rd, err := uc.client.GetDataByID(ctx, &pb.GetDataByIDRequest{Id: id})
		if err != nil {
			return err
		}
		r := clientstorage.SaveDataRequestModel{
			ID:                    id,
			Userid:                userID,
			Description:           rd.Description,
			Datatype:              getDataTypeStringFromPB(rd.Type),
			Hash:                  rd.Hash,
			ModificationTimestamp: rd.ModificationTimestamp,
			Data:                  rd.Data,
		}
		err = uc.storage.SaveData(ctx, r)
		if err != nil {
			return err
		}
	}
	for _, id := range idsToAddRemote {
		item, err := uc.storage.GetDataByID(ctx, id)
		if err != nil {
			return err
		}
		r := &pb.SaveDataRequest{
			Id:                    item.ID,
			Description:           item.Description,
			Type:                  getPBDataTypeFromString(item.Datatype),
			Hash:                  item.Hash,
			ModificationTimestamp: item.ModificationTimestamp,
			Data:                  item.Data,
		}
		_, err = uc.client.SaveData(ctx, r)
		if err != nil {
			return uccommon.NewServerError(err)
		}
	}
	for _, id := range idsToUpdateLocal {
		res, err := uc.client.GetDataByID(ctx, &pb.GetDataByIDRequest{Id: id})
		if err != nil {
			return err
		}

		r := clientstorage.UpdateDataRequestModel{
			ID:                    id,
			Hash:                  res.Hash,
			ModificationTimestamp: res.ModificationTimestamp,
			Data:                  res.Data,
		}

		err = uc.storage.UpdateData(ctx, r)
		if err != nil {
			return err
		}
	}
	for _, id := range idsToUpdateRemote {
		res, err := uc.storage.GetDataByID(ctx, id)
		if err != nil {
			return err
		}

		r := &pb.UpdateDataRequest{
			Id:                    res.ID,
			Hash:                  res.Hash,
			ModificationTimestamp: res.ModificationTimestamp,
			Data:                  res.Data,
		}
		_, err = uc.client.UpdateData(ctx, r)
		if err != nil {
			return uccommon.NewServerError(err)
		}
	}
	return nil
}

func (uc *Usecase) UpdateTextData(ctx context.Context, dataid string, text string) error {
	// формируем данные для обновления
	cipher, err := uc.getDataCipher()
	if err != nil {
		return err
	}
	encrypted := cipher.Encrypt([]byte(text))
	hash := uc.getDataHash([]byte(text))
	modts := uc.getModificationTimestamp()

	err = uc.storage.UpdateData(ctx, clientstorage.UpdateDataRequestModel{
		ID:                    dataid,
		Hash:                  hash,
		ModificationTimestamp: modts,
		Data:                  encrypted,
	})
	if err != nil {
		return err
	}

	_, err = uc.client.UpdateData(ctx, &pb.UpdateDataRequest{
		Id:                    dataid,
		Hash:                  hash,
		ModificationTimestamp: modts,
		Data:                  encrypted,
	})
	if err != nil {
		return uccommon.NewServerError(err)
	}

	return nil
}

func (uc *Usecase) RemoveDataByID(ctx context.Context, dataid string) error {
	modts := uc.getModificationTimestamp()
	err := uc.storage.RemoveDataByID(ctx, clientstorage.RemoveDataByIDRequestModel{
		DataID:                dataid,
		ModificationTimestamp: modts,
	})
	if err != nil {
		return err
	}
	_, err = uc.client.RemoveData(ctx, &pb.RemoveDataRequest{
		Id:                    dataid,
		ModificationTimestamp: modts,
	})
	if err != nil {
		return uccommon.NewServerError(err)
	}

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

func (uc *Usecase) saveBytes(ctx context.Context, b []byte, description string, datatype string) error {
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
	encrypted := cipher.Encrypt(b)
	hash := uc.getDataHash(b)
	modts := uc.getModificationTimestamp()
	userID, err := uc.storage.GetUserID()
	if err != nil {
		return err
	}

	// пробуем сохранить локально
	r := clientstorage.SaveDataRequestModel{
		ID:                    id,
		Userid:                userID,
		Description:           description,
		Datatype:              datatype,
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
		return uccommon.NewServerError(err)
	}

	return nil
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

func getDataTypeStringFromPB(dt pb.DataTypeEnum) string {
	switch dt {
	case pb.DataTypeEnum_DATA_TYPE_CREDENTIALS:
		return "CREDENTIALS"
	case pb.DataTypeEnum_DATA_TYPE_TEXT:
		return "TEXT"
	case pb.DataTypeEnum_DATA_TYPE_BINARY:
		return "BINARY"
	case pb.DataTypeEnum_DATA_TYPE_BANK_CARD:
		return "BANK_CARD"
	}

	return "TEXT"
}
