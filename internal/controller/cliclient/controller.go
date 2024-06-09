package cliclient

import (
	"context"

	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientauth"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientstore"
)

type Controller struct {
	ucauth  *clientauth.Usecase
	ucstore *clientstore.Usecase
}

var controller Controller

func New(ucauth *clientauth.Usecase, ucstore *clientstore.Usecase) *Controller {
	controller = Controller{
		ucauth:  ucauth,
		ucstore: ucstore,
	}

	return &controller
}

func (c *Controller) Serve(ctx context.Context) error {
	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}
