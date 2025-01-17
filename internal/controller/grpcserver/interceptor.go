package grpcserver

import (
	"context"
	"strings"

	"github.com/KartoonYoko/gophkeeper/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// InterceptorAuthKey тип ключа перехватчика аутентификации в контексте
type InterceptorAuthKey int

const (
	ctxKeyUserID InterceptorAuthKey = iota // ключ для ID пользователя
)

// interceptorAuth проверяет наличие симметрично подписанного токена
func (c *Controller) interceptorAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// данный перехватчик должен:
	// 	- игнорировать некоторые функции
	// 	- провалидировать полученный токен
	// 	- внести токен в контекст

	if info.FullMethod == "/proto.AuthService/Login" ||
		info.FullMethod == "/proto.AuthService/Register" ||
		info.FullMethod == "/proto.AuthService/RefreshToken" {
		return handler(ctx, req)
	}

	var err error
	var userID string

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Log.Error("can not get metadata")
		return nil, status.Error(codes.Internal, "")
	}

	token, err := getJWTTokenFromMetadata(&md)
	if err != nil {
		logger.Log.Error("get token error", zap.Error(err))
		return nil, err
	}

	userID, err = c.usecaseAuth.ValidateJWTString(token)
	if err != nil {
		logger.Log.Error("can not validate and get user ID: ", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "token is wrong")
	}

	ctx = context.WithValue(ctx, ctxKeyUserID, userID)

	return handler(ctx, req)
}

func getJWTTokenFromMetadata(md *metadata.MD) (string, error) {
	sl := md.Get("Authorization")
	if len(sl) == 0 {
		return "", status.Error(codes.Unauthenticated, "not found token")
	}

	token, _ := strings.CutPrefix(sl[0], "Bearer ")

	if len(token) == 0 {
		return "", status.Error(codes.Unauthenticated, "not found token")
	}

	return token, nil
}
