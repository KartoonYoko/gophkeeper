package cliclient

import (
	"context"

	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientauth"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientstore"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/clientversion"
)

// Controller контроллер клиентского приложения
type Controller struct {
	ucauth    *clientauth.Usecase
	ucstore   *clientstore.Usecase
	ucversion *clientversion.Usecase
}

// controller экземпляр контроллера
var controller Controller

// New ининциализирует экземпляр контроллера в пакете
func New(
	ucauth *clientauth.Usecase,
	ucstore *clientstore.Usecase,
	ucversion *clientversion.Usecase) *Controller {
	controller = Controller{
		ucauth:    ucauth,
		ucstore:   ucstore,
		ucversion: ucversion,
	}

	return &controller
}

// Serve запустить обработку cli команд
func (c *Controller) Serve(ctx context.Context) error {
	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}
