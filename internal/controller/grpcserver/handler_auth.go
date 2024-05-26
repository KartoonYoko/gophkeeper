package grpcserver

import (
	"context"

	"github.com/KartoonYoko/gophkeeper/internal/controller/common"
	pb "github.com/KartoonYoko/gophkeeper/internal/controller/grpcserver/proto"
	"github.com/KartoonYoko/gophkeeper/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login авторизует пользователя и отдаёт ему ключ для шифрования
func (c *Controller) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "uninmplemented")
}

// Logout
func (c *Controller) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "uninmplemented")
}

// RefreshToken обновляет токен доступа
func (c *Controller) RefreshToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "uninmplemented")
}

// Register добавляет новго пользователя в систему и отдаёт ему ключ для шифрования
func (c *Controller) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	c.usecaseAuth.Register(ctx, request.Login, request.Password)
	accessToken, err := common.BuildJWTString("1", c.conf.SecretJWTKey)
	if err != nil {
		logger.Log.Error("failed to build jwt string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &pb.RegisterResponse{
		Token: &pb.Token{
			AccessToken:  accessToken,
			RefreshToken: "",
		},
	}, nil
}
