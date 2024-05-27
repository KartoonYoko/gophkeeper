package grpcserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/KartoonYoko/gophkeeper/internal/controller/common"
	pb "github.com/KartoonYoko/gophkeeper/internal/controller/grpcserver/proto"
	"github.com/KartoonYoko/gophkeeper/internal/logger"
	"google.golang.org/grpc"
)

type Controller struct {
	usecaseAuth AuthUsecase

	conf Config
	pb.AuthServiceServer
}

func New(conf Config, usecaseAuth AuthUsecase) *Controller {
	c := new(Controller)
	c.conf = conf

	c.usecaseAuth = usecaseAuth

	return c
}

type buildJWTStringClaims struct {
	UserID string
}

func (c *Controller) buildJWTString(cl buildJWTStringClaims) (string, error) {
	return common.BuildJWTString(cl.UserID, c.conf.SecretJWTKey, c.conf.JWTDurationMinute)
}

func (c *Controller) getUserIDFromContext(ctx context.Context) (string, error) {
	somevalue := ctx.Value(ctxKeyUserID)
	userID, ok := somevalue.(string)
	if !ok {
		return "", fmt.Errorf("failed to get user ID from context")
	}

	return userID, nil
}

func (c *Controller) Serve(ctx context.Context) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	listen, err := net.Listen("tcp", c.conf.BootstrapAddress)
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	grpcServer := grpc.NewServer()
	// grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
	// 	c.interceptorRequestTime,
	// 	c.interceptorAuth,
	// ))
	pb.RegisterAuthServiceServer(grpcServer, c)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s := <-sigCh
		logger.Log.Info(fmt.Sprintf("got signal %v, attempting graceful shutdown", s))
		cancel()

		grpcServer.GracefulStop()
		wg.Done()
	}()

	logger.Log.Info(fmt.Sprintf("grpc serve on %s", c.conf.BootstrapAddress))
	if err := grpcServer.Serve(listen); err != nil {
		return fmt.Errorf("serve error grpc server: %w", err)
	}
	wg.Wait()
	logger.Log.Info("grpc server stopped gracefully")

	return nil
}