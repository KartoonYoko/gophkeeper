package grpcserver

import (
	"context"
	"errors"

	pb "github.com/KartoonYoko/gophkeeper/internal/controller/grpcserver/proto"
	"github.com/KartoonYoko/gophkeeper/internal/logger"
	ucauth "github.com/KartoonYoko/gophkeeper/internal/usecase/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login авторизует пользователя и отдаёт ему ключ для шифрования
func (c *Controller) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, err := c.usecaseAuth.Login(ctx, request.Login, request.Password)
	if err != nil {
		var exsterror *ucauth.LoginOrPasswordNotFoundError
		if errors.As(err, &exsterror) {
			return nil, status.Errorf(codes.Unauthenticated, "login and(or) password not found")
		}

		logger.Log.Error("failed to login user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	accessToken, err := c.buildJWTString(buildJWTStringClaims{
		UserID: result.UserID,
	})
	if err != nil {
		logger.Log.Error("failed to build jwt string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.LoginResponse{
		Token: &pb.Token{
			AccessToken:  accessToken,
			RefreshToken: result.RefreshToken,
		},
	}, nil
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
	result, err := c.usecaseAuth.Register(ctx, request.Login, request.Password)
	if err != nil {
		var exsterror *ucauth.LoginAlreadyExistsError
		if errors.As(err, &exsterror) {
			return nil, status.Errorf(codes.AlreadyExists, "login already exists")
		}

		logger.Log.Error("failed to register user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	accessToken, err := c.buildJWTString(buildJWTStringClaims{
		UserID: result.UserID,
	})
	if err != nil {
		logger.Log.Error("failed to build jwt string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.RegisterResponse{
		Token: &pb.Token{
			AccessToken:  accessToken,
			RefreshToken: result.RefreshToken,
		},
	}, nil
}
