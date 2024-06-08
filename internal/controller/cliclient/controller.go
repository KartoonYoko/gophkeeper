package cliclient

import "context"

type Controller struct{}

func New() *Controller {
	return new(Controller)
}

func (c *Controller) Serve(ctx context.Context) error {
	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}
