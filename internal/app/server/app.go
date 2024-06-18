/*
Package server пакет для запуска серверного приложения
*/
package server

import (
	"context"
	"fmt"
	"log"

	"github.com/KartoonYoko/gophkeeper/internal/common/datacipher"
	"github.com/KartoonYoko/gophkeeper/internal/common/passwordhash"
	"github.com/KartoonYoko/gophkeeper/internal/common/secretkeycipher"
	grpcserver "github.com/KartoonYoko/gophkeeper/internal/controller/grpcserver"
	"github.com/KartoonYoko/gophkeeper/internal/logger"
	storageMinio "github.com/KartoonYoko/gophkeeper/internal/storage/miniostorage"
	storagePostgres "github.com/KartoonYoko/gophkeeper/internal/storage/postgres"
	usecaseAuth "github.com/KartoonYoko/gophkeeper/internal/usecase/auth"
	usecaseStore "github.com/KartoonYoko/gophkeeper/internal/usecase/store"
	"go.uber.org/zap"
)

// Run запуск серверного приложения
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
	sch, err := secretkeycipher.New(config.UserSecretKeySecure)
	if err != nil {
		logger.Log.Error("secret key cipher init error", zap.Error(err))
		return
	}
	ph := passwordhash.New()
	ucAuth := usecaseAuth.New(psSt, ph, sch, ucAConf)
	sConf := usecaseStore.Config{
		SecretKeySecure: config.UserSecretKeySecure,
		DataSecretKey:   config.DataSecretKeySecure,
	}
	d, err := datacipher.New(config.DataSecretKeySecure)
	if err != nil {
		logger.Log.Error("data cipher init error", zap.Error(err))
		return
	}
	ucStore := usecaseStore.New(sConf, psSt, mstorage, d)

	// server
	grpcConf := grpcserver.Config{
		BootstrapAddress: config.ServerAddress,
	}
	grpcController := grpcserver.New(grpcConf, ucAuth, ucStore)
	if err := grpcController.Serve(ctx); err != nil {
		logger.Log.Error("grpc serve error: %s", zap.Error(err))
	}
}
