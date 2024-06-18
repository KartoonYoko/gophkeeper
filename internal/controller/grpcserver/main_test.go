package grpcserver

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	ucauth "github.com/KartoonYoko/gophkeeper/internal/usecase/auth"
	ucstore "github.com/KartoonYoko/gophkeeper/internal/usecase/store"
	"github.com/KartoonYoko/gophkeeper/internal/common/datacipher"
	"github.com/KartoonYoko/gophkeeper/internal/common/passwordhash"
	"github.com/KartoonYoko/gophkeeper/internal/common/secretkeycipher"
)

var (
	controller           *Controller
	usecaseAuth          *ucauth.Usecase
	usecaseStore         *ucstore.Usecase
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
	ph := passwordhash.New()
	sch, err := secretkeycipher.New("some secret key secure")
	if err != nil {
		return fmt.Errorf("failed to create secret key cipher: %w", err)
	}
	usecaseAuth = ucauth.New(nil, ph, sch, ucauth.Config{
		RefreshTokenDurationMinute: 60,
		SecretJWTKey:               "some secret jwt key",
		JWTDurationMinute:          10,
		SecretKeySecure:            "some secret key secure",
	})

	ds, err := datacipher.New("secretkey")
	if err != nil {
		return fmt.Errorf("failed to create data cipher: %w", err)
	}
	usecaseStore = ucstore.New(ucstore.Config{
		SecretKeySecure: "secretkey",
		DataSecretKey:   "secretkey",
	}, nil, nil, ds)

	controller = New(Config{
		BootstrapAddress: bootstrapAddressgRPC,
	}, usecaseAuth, usecaseStore)

	return nil
}

func createJWTString(userID string) (string, error) {
	return usecaseAuth.BuildJWTString(userID)
}

func encrypteData(data []byte) ([]byte, error) {
	h, err := datacipher.New("secretkey")
	if err != nil {
		return nil, err
	}
	
	return h.Encrypt(data), nil
}
