package cliclient

import (
	"context"
	"log"

	"github.com/KartoonYoko/gophkeeper/internal/controller/cliclient"
	"github.com/KartoonYoko/gophkeeper/internal/storage/clientstorage"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientauth"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientstore"
)

func Run() {
	var err error
	ctx := context.Background()

	conf, err := NewConfig()
	if err != nil {
		log.Fatalf("config err: %s", err)
	}

	tokenstore, err := clientstorage.New(ctx)
	if err != nil {
		log.Fatalf("failed init store: %v", err)
	}
	defer tokenstore.Close()

	sc, err := NewServerConnection(conf, tokenstore)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// usecases
	ucauth := clientauth.New(sc.conn, tokenstore)
	ucstore := clientstore.New(sc.conn, tokenstore)

	controller := cliclient.New(ucauth, ucstore)

	err = controller.Serve(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
