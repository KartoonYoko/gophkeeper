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
	storagePostgres "github.com/KartoonYoko/gophkeeper/internal/storage/postgres"
	storageMinio "github.com/KartoonYoko/gophkeeper/internal/storage/miniostorage"
	usecaseAuth "github.com/KartoonYoko/gophkeeper/internal/usecase/auth"
	usecaseStore "github.com/KartoonYoko/gophkeeper/internal/usecase/store"
	"go.uber.org/zap"
)

func Run() {
	ctx := context.Background()

	// логгер
	if err := logger.Initialize("Info"); err != nil {
		log.Fatal(fmt.Errorf("logger init error: %w", err))
	}
	defer logger.Log.Sync()

	// storage
	psConf := storagePostgres.Config{
		ConnectionString: "host=localhost user=postgres password=123 dbname=gophkeeper port=5433 sslmode=disable",
	}
	psSt, err := storagePostgres.New(ctx, psConf)
	if err != nil {
		logger.Log.Error("storage init error", zap.Error(err))
		return
	}
	msConf := storageMinio.Config{
		Endpoint: "",
		AccessKeyID: "",
		SecretAccessKey: "",
	}
	mstorage, err := storageMinio.NewStorage(msConf)
	if err != nil {
		logger.Log.Error("minio storage init error", zap.Error(err))
		return
	}
	// usecases
	passwordSault := "somesault"
	ucAConf := usecaseAuth.Config{
		RefreshTokenDurationMinute: 60,
		SecretJWTKey:               "somesecretjwtkey",
		JWTDurationMinute:          5,
		PasswordSault:              passwordSault,
	}
	ucAuth := usecaseAuth.New(psSt, ucAConf)
	sConf := usecaseStore.Config{
		SecretKeySecure: "",
		DataSecretKey: "",
	}
	usecaseStore.New(sConf, psSt, mstorage)

	// server
	grpcConf := grpcserver.Config{
		BootstrapAddress: ":8080",
	}
	grpcController := grpcserver.New(grpcConf, ucAuth)
	if err := grpcController.Serve(ctx); err != nil {
		logger.Log.Error("grpc serve error: %s", zap.Error(err))
	}
}
