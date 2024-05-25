package grpcserver

import (
	"context"
	"strings"

	"github.com/KartoonYoko/gophkeeper/internal/controller/common"
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
	keyUserID InterceptorAuthKey = iota // ключ для ID пользователя
)

// interceptorAuth проверяет наличие симметрично подписанного токена
func (c *Controller) interceptorAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var err error
	var userID string

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Log.Error("can not get metadata")
		return nil, status.Error(codes.Internal, "")
	}

	sl := md.Get("Authorization")
	if len(sl) == 0 {
		userID, err = c.setAuthorizationMetadata(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		token, _ := strings.CutPrefix(sl[0], "Bearer ")
		userID, err = common.ValidateAndGetUserID(token, c.conf.SecretJWTKey)
		if err != nil {
			logger.Log.Error("can not validate and get user ID: ", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "token is wrong")
		}
	}

	ctx = context.WithValue(ctx, keyUserID, userID)

	return handler(ctx, req)
}

func (c *Controller) setAuthorizationMetadata(ctx context.Context) (string, error) {
	return "", nil
	// var err error
	// var userID string

	// userID, err = c.ucAuth.GetNewUserID(ctx)
	// if err != nil {
	// 	logger.Log.Error("can not get new user ID: ", zap.Error(err))
	// 	return "", status.Error(codes.Internal, "")
	// }
	// jwt, err := common.BuildJWTString(userID)
	// if err != nil {
	// 	logger.Log.Error("can not build JWT string: ", zap.Error(err))
	// 	return "", status.Error(codes.Internal, "")
	// }
	// bearerStr := fmt.Sprintf("Bearer %s", jwt)
	// grpc.SetHeader(ctx, metadata.New(map[string]string{"Authorization": bearerStr}))

	// return userID, nil
}
