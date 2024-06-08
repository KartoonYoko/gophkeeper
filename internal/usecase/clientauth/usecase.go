package clientauth

import (
	"context"
	"fmt"

	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
	"github.com/KartoonYoko/gophkeeper/internal/storage/clientstorage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Usecase struct {
	client pb.AuthServiceClient

	storage *clientstorage.Storage
}

func New(conn *grpc.ClientConn, store *clientstorage.Storage) *Usecase {
	uc := new(Usecase)

	uc.client = pb.NewAuthServiceClient(conn)
	uc.storage = store

	return uc
}

func (uc *Usecase) Login(ctx context.Context, login string, password string) error {
	// попытаться залогиниться
	// если ошибка, то сообщить и выход
	// если успех, то сохранить токен
	request := &pb.LoginRequest{
		Login:    login,
		Password: password,
	}
	response, err := uc.client.Login(ctx, request)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				// сообщить ошибкой, что пользователь не найден
				return fmt.Errorf("login or password not found")
			} else {
				// сообщить о неизвестной ошибке
				return fmt.Errorf("unhandled error: %w", err)
			}
		} else {
			// сообщить о неизвестной ошибке
			return fmt.Errorf("unhandled error: %w", err)
		}
	}

	// save token
	err = uc.storage.SaveTokens(ctx, response.Token.AccessToken, response.Token.RefreshToken, response.SecretKey)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) Logout(ctx context.Context) error {
	_, rt, err := uc.storage.GetTokens()
	if err != nil {
		return err
	}

	request := &pb.LogoutRequest{
		RefreshToken: rt,
	}
	_, err = uc.client.Logout(ctx, request)
	if err != nil {
		return err
	}

	err = uc.storage.RemoveTokens()
	if err != nil {
		return err
	}

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

	err = uc.storage.SaveTokens(ctx, response.Token.AccessToken, response.Token.RefreshToken, response.SecretKey)
	if err != nil {
		return err
	}

	return nil
}
