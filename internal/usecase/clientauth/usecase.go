package clientauth

import (
	"context"
	"log"

	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
	"github.com/KartoonYoko/gophkeeper/internal/storage/clientstorage"
	"google.golang.org/grpc"
)

type Usecase struct {
	client pb.AuthServiceClient

	storage *clientstorage.Storage
}

func New(conn *grpc.ClientConn) *Usecase {
	uc := new(Usecase)

	uc.client = pb.NewAuthServiceClient(conn)
	var err error
	uc.storage, err = clientstorage.New()
	if err != nil {
		// todo либо вернуть ошибку, либо обработать;
		// а лучше прокинуть сюда хранилище через параметр
		log.Fatal(err)
	}

	return uc
}

func (uc *Usecase) Login(ctx context.Context, login string, password string) error {
	return nil
	// TODO
	// попытаться залогиниться
	// если ошибка, то сообщить и выход
	// если успех, то сохранить токен, а также ключ для шифровки/дешифровки
	//

	// request := &pb.LoginRequest{
	// 	Login:    login,
	// 	Password: password,
	// }
	// response, err := uc.client.Login(ctx, request)
	// if err != nil {
	// 	if e, ok := status.FromError(err); ok {
	// 		if e.Code() == codes.NotFound {
	// 			// сообщить ошибкой, что пользователь не найден
	// 		} else {
	// 			// сообщить о неизвестной ошибке
	// 		}
	// 	} else {
	// 		// сообщить о неизвестной ошибке
	// 	}
	// }

	// todo save token

}

func (uc *Usecase) Logout(ctx context.Context) error {
	return nil
}

func (uc *Usecase) Register(ctx context.Context, login string, password string) error {
	request := &pb.RegisterRequest{
		Login:    login,
		Password: password,
	}
	response, err := uc.client.Register(ctx, request)
	if err != nil {
		return err
	}
	
	err = uc.storage.SaveTokens(ctx, response.Token.AccessToken, response.Token.RefreshToken)
	if err != nil {
		return err
	}

	return nil
}
