/*
Package server пакет для запуска серверного приложения
*/
package server

import (
	"context"
	"fmt"
	"log"

	grpcserver "github.com/KartoonYoko/gophkeeper/internal/controller/grpcserver"
	"github.com/KartoonYoko/gophkeeper/internal/logger"
	storageMinio "github.com/KartoonYoko/gophkeeper/internal/storage/miniostorage"
	storagePostgres "github.com/KartoonYoko/gophkeeper/internal/storage/postgres"
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

	config, err := NewConfig()
	if err != nil {
		logger.Log.Error("config init error", zap.Error(err))
		return
	}

	{
		logger.Log.Info("database dsn", zap.Any("dsn", fmt.Sprintf("%v", config)))
	}

	// storage
	psConf := storagePostgres.Config{
		ConnectionString: config.DatabaseDsn,
	}
	psSt, err := storagePostgres.New(ctx, psConf)
	if err != nil {
		logger.Log.Error("storage init error", zap.Error(err))
		return
	}
	msConf := storageMinio.Config{
		Endpoint:        config.MinioAddress,
		AccessKeyID:     config.MinioAccessKeyID,
		SecretAccessKey: config.MinioSecretAccessKey,
	}
	mstorage, err := storageMinio.NewStorage(msConf)
	if err != nil {
		logger.Log.Error("minio storage init error", zap.Error(err))
		return
	}
	// usecases
	ucAConf := usecaseAuth.Config{
		SecretJWTKey:               config.SecretJWTKey,
		JWTDurationMinute:          config.JWTTokenLifetimeMinutes,
		RefreshTokenDurationMinute: config.RefreshTokenLifeimeMinutes,
	}
	ucAuth, err := usecaseAuth.New(psSt, ucAConf)
	if err != nil {
		logger.Log.Error("usecase auth init error", zap.Error(err))
		return
	}
	sConf := usecaseStore.Config{
		SecretKeySecure: config.UserSecretKeySecure,
		DataSecretKey:   config.DataSecretKeySecure,
	}
	ucStore, err := usecaseStore.New(sConf, psSt, mstorage)
	if err != nil {
		logger.Log.Error("usecase store init error", zap.Error(err))
		return
	}

	// server
	grpcConf := grpcserver.Config{
		BootstrapAddress: config.ServerAddress,
	}
	grpcController := grpcserver.New(grpcConf, ucAuth, ucStore)
	if err := grpcController.Serve(ctx); err != nil {
		logger.Log.Error("grpc serve error: %s", zap.Error(err))
	}
}
