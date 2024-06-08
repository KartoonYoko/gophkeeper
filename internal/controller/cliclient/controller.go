package cliclient

import (
	"context"

	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientauth"
)

type Controller struct{
	ucauth *clientauth.Usecase
}

var controller Controller

func New(ucauth *clientauth.Usecase) *Controller {
	controller = Controller{
		ucauth: ucauth,
	}

	return &controller
}

func (c *Controller) Serve(ctx context.Context) error {
	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}
