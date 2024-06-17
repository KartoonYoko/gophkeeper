package grpcserver

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	ucauth "github.com/KartoonYoko/gophkeeper/internal/usecase/auth"
)

var (
	controller           *Controller
	usecaseAuth          *ucauth.Usecase
	bootstrapAddressgRPC string
)

func TestMain(m *testing.M) {
	bootstrapAddressgRPC = ":8080"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		fmt.Println("start server")
		err := createTestMock()
		if err != nil {
			log.Fatalf("server start error: %v", err)
		}

		err = controller.Serve(ctx)
		if err != nil {
			log.Fatalf("server start error: %v", err)
		}

		fmt.Println("stop server")
	}()

	// ждём пока запуститься сервер
	time.Sleep(2 * time.Second)
	m.Run()
}

// createTestMock собирает контроллер
func createTestMock() error {
	var err error
	usecaseAuth, err = ucauth.New(nil, ucauth.Config{
		RefreshTokenDurationMinute: 60,
		SecretJWTKey:               "some secret jwt key",
		JWTDurationMinute:          10,
		SecretKeySecure:            "some secret key secure",
	})
	if err != nil {
		return fmt.Errorf("failed to create usecase auth: %w", err)
	}

	controller = New(Config{
		BootstrapAddress: bootstrapAddressgRPC,
	}, usecaseAuth, nil)

	return nil
}
