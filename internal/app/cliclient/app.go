package cliclient

import (
	"context"
	"log"

	"github.com/KartoonYoko/gophkeeper/internal/controller/cliclient"
	"github.com/KartoonYoko/gophkeeper/internal/storage/clientstorage"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientauth"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientstore"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientversion"
)

func Run(vi *VersionInfo) {
	if vi == nil {
		vi = &VersionInfo{}
	}

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
	ucversion := clientversion.New(vi.Version, vi.BuildDate)

	controller := cliclient.New(ucauth, ucstore, ucversion)

	err = controller.Serve(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
