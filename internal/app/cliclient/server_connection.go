package cliclient

import (
	"context"
	"fmt"
	"log"

	"github.com/KartoonYoko/gophkeeper/internal/storage/clientstorage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
)

// ServerConnection фасад для работы с соединением и его перехватчиками
type ServerConnection struct {
	tokenstore *clientstorage.Storage

	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

// NewServerConnection конструктор
func NewServerConnection(conf *Config, tokenstore *clientstorage.Storage) (*ServerConnection, error) {
	sc := new(ServerConnection)

	conn, err := grpc.NewClient(
		conf.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(sc.authInterpector, sc.refreshTokenInterpector))

	if err != nil {
		log.Fatal(err)
	}

	sc.client = pb.NewAuthServiceClient(conn)
	sc.tokenstore = tokenstore
	sc.conn = conn

	return sc, nil
}

// Close закрывает соединение
func (sc *ServerConnection) Close() error {
	return sc.conn.Close()
}

// authInterpector перехватчик для подстановки токена доступа в метаданные
func (sc *ServerConnection) authInterpector(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	if method == "/proto.AuthService/Login" ||
		method == "/proto.AuthService/Register" ||
		method == "/proto.AuthService/RefreshToken" {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
	at, _, err := sc.tokenstore.GetTokens()
	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{"Authorization": at})
	ctx = metadata.NewOutgoingContext(ctx, md)
	err = invoker(ctx, method, req, reply, cc, opts...)

	return err
}

// refreshTokenInterpector перехватчик для автоматического обновления токена доступа
func (sc *ServerConnection) refreshTokenInterpector(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	if method == "/proto.AuthService/Login" ||
		method == "/proto.AuthService/Register" ||
		method == "/proto.AuthService/RefreshToken" {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		e, ok := status.FromError(err)
		if !ok {
			return err
		}

		if e.Code() != codes.Unauthenticated {
			return err
		}

		at, rt, err := sc.tokenstore.GetTokens()
		if err != nil {
			return status.Error(codes.Unauthenticated, fmt.Sprintf("client failed refresh tokens: %s", err))
		}

		res, err := sc.client.RefreshToken(ctx, &pb.RefreshTokenRequest{
			Token: &pb.Token{
				AccessToken:  at,
				RefreshToken: rt,
			},
		})
		if err != nil {
			return status.Error(codes.Unauthenticated, fmt.Sprintf("client failed refresh tokens: %s", err))
		}

		err = sc.tokenstore.UpdateTokens(res.Token.AccessToken, res.Token.RefreshToken)
		if err != nil {
			return status.Error(codes.Unauthenticated, fmt.Sprintf("client failed refresh tokens: %s", err))
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}

	return err
}
