package cliclient

import (
	"context"
	"log"

	"github.com/KartoonYoko/gophkeeper/internal/controller/cliclient"
	"github.com/KartoonYoko/gophkeeper/internal/storage/clientstorage"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var tokenstore *clientstorage.Storage

func Run() {
	var err error
	// todo собрать приложение
	ctx := context.Background()
	tokenstore, err = clientstorage.New()
	if err != nil {
		log.Fatalf("failed init store: %v", err)
	}

	conn, err := grpc.NewClient(
		":8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterpector))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// usecases
	ucauth := clientauth.New(conn, tokenstore)

	controller := cliclient.New(ucauth)

	err = controller.Serve(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func authInterpector(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	if method == "" {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
	at, _, err := tokenstore.GetTokens()
	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{"Authorization": at})
	ctx = metadata.NewOutgoingContext(ctx, md)
	err = invoker(ctx, method, req, reply, cc, opts...)

	return err
}
