package grpcserver

import (
	"context"
	"errors"

	"github.com/KartoonYoko/gophkeeper/internal/logger"
	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
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
			logger.Log.Info("login or password not found", zap.String("login", request.Login))
			return nil, status.Errorf(codes.Unauthenticated, "login and(or) password not found")
		}

		var loginNotFoundError *ucauth.LoginNotFoundError
		if errors.As(err, &loginNotFoundError) {
			logger.Log.Info("login not found", zap.String("login", request.Login))
			return nil, status.Errorf(codes.Unauthenticated, "login and(or) password not found")
		}

		logger.Log.Error("failed to login user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	accessToken, err := c.usecaseAuth.BuildJWTString(result.UserID)
	if err != nil {
		logger.Log.Error("failed to build jwt string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.LoginResponse{
		Token: &pb.Token{
			AccessToken:  accessToken,
			RefreshToken: result.RefreshToken,
		},
		SecretKey: result.SecretKey,
		UserId:    result.UserID,
	}, nil
}

// Logout
func (c *Controller) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		logger.Log.Error("can not get user ID from context", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	err = c.usecaseAuth.Logout(ctx, userID, request.RefreshToken)
	if err != nil {
		logger.Log.Error("can not logout", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.LogoutResponse{}, nil
}

// RefreshToken обновляет токен доступа
//
// TODO если refresh token недействительный, то возвращать ошибку "неавторизован"
func (c *Controller) RefreshToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	res, err := c.usecaseAuth.RefreshToken(ctx, request.Token.RefreshToken)
	if err != nil {
		var exsterror *ucauth.RefreshTokenNotFoundError
		if errors.As(err, &exsterror) {
			return nil, status.Errorf(codes.Unauthenticated, "refresh token not found")
		}

		logger.Log.Error("can not refresh token", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	accessToken, err := c.usecaseAuth.BuildJWTString(res.UserID)
	if err != nil {
		logger.Log.Error("failed to build jwt string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.RefreshTokenResponse{
		Token: &pb.Token{
			AccessToken:  accessToken,
			RefreshToken: res.RefreshToken,
		},
	}, nil
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

	accessToken, err := c.usecaseAuth.BuildJWTString(result.UserID)
	if err != nil {
		logger.Log.Error("failed to build jwt string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.RegisterResponse{
		Token: &pb.Token{
			AccessToken:  accessToken,
			RefreshToken: result.RefreshToken,
		},
		SecretKey: result.SecretKey,
		UserId:    result.UserID,
	}, nil
}
