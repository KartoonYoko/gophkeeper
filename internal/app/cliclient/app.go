package cliclient

import (
	"context"
	"log"

	"github.com/KartoonYoko/gophkeeper/internal/controller/cliclient"
)

func Run() {
	// todo собрать приложение
	ctx := context.Background()

	controller := cliclient.New()

	err := controller.Serve(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
