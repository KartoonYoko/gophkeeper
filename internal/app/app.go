/*
Package app пакет для запуска приложения
*/
package app

import (
	"context"
	"fmt"
	"log"

	grpcserver "github.com/KartoonYoko/gophkeeper/internal/controller/grpcserver"
	"github.com/KartoonYoko/gophkeeper/internal/logger"
	"go.uber.org/zap"
)

func Run() {
	ctx := context.Background()

	// логгер
	if err := logger.Initialize("Info"); err != nil {
		log.Fatal(fmt.Errorf("logger init error: %w", err))
	}
	defer logger.Log.Sync()

	grpcConf := grpcserver.Config{
		BootstrapAddress: ":8080",
		SecretJWTKey:     "somesecretjwtkey",
	}
	grpcController := grpcserver.New(grpcConf)
	if err := grpcController.Serve(ctx); err != nil {
		logger.Log.Error("grpc serve error: %s", zap.Error(err))
	}
}
