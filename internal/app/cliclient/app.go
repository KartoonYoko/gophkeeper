package cliclient

import (
	"context"
	"log"

	"github.com/KartoonYoko/gophkeeper/internal/controller/cliclient"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run() {
	// todo собрать приложение
	ctx := context.Background()

	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// usecases
	ucauth := clientauth.New(conn)

	controller := cliclient.New(ucauth)

	err = controller.Serve(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
