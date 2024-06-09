package clientstore

import (
	"context"
	"strings"

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

func (uc *Usecase) CreateTextData(ctx context.Context, text string) error {
	cipher, err := uc.getDataCipher()
	if err != nil {
		return err
	}

	id := uuid.New().String()
	encrypted := cipher.Encrypt([]byte(text))

	userID, err := uc.storage.GetUserID()
	if err != nil {
		return err
	}
	r := clientstorage.SaveDataRequestModel{
		Filename:    id,
		Userid:      userID,
		Description: "",
		Datatype:    "TEXT",
		Data:        encrypted,
	}
	err = uc.storage.SaveData(ctx, r)
	if err != nil {
		return err
	}

	rr := &pb.SaveDataRequest{
		Description: r.Description,
		Type:        getPBDataTypeFromString(r.Datatype),
	}
	uc.client.SaveData(ctx)
	// TODO
	// - шифрую данные и работаю далее только с шифрованными данными
	// - сохраняю локльно на диске (под каким именем сохранять? использовать guid, придётся поменять схему БД на сервере)
	// - сохраняю информацию в БД
	// - пробую сохранить удалённо
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
