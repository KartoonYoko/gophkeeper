package grpcserver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KartoonYoko/gophkeeper/internal/common/passwordhash"
	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
	serror "github.com/KartoonYoko/gophkeeper/internal/storage/error/auth"
	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/auth/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestController_Register(t *testing.T) {
	ctx := context.Background()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(bootstrapAddressgRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)

	type args struct {
		login   string
		pasword string
	}
	type test struct {
		name            string
		args            args
		prepare         func(mock *mocks.MockStorager)
		statusErrorCode codes.Code
	}
	tests := []test{
		{
			name: "Success",
			args: args{
				login:   "login",
				pasword: "password",
			},
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, nil)
				m.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "Error user already exists",
			args: args{
				login:   "login",
				pasword: "password",
			},
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil, serror.NewLoginAlreadyExistsError("login"))
			},
			statusErrorCode: codes.AlreadyExists,
		},
		{
			name: "Error create user",
			args: args{
				login:   "login",
				pasword: "password",
			},
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of create user"))
			},
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error create refresh token",
			args: args{
				login:   "login",
				pasword: "password",
			},
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil, nil)
				m.EXPECT().
					CreateRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of creating refresh token"))
			},
			statusErrorCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks.NewMockStorager(ctrl)

			if tt.prepare != nil {
				tt.prepare(m)
			}

			usecaseAuth.Storage = m

			request := new(pb.RegisterRequest)
			request.Login = "login"
			request.Password = "password"

			_, err := c.Register(ctx, request)

			if tt.statusErrorCode == 0 {
				require.NoError(t, err)
			} else {
				if e, ok := status.FromError(err); ok {
					require.Equal(t, tt.statusErrorCode.String(), e.Code().String())
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			}

		})
	}
}

func TestController_Logout(t *testing.T) {
	ctx := context.Background()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(bootstrapAddressgRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)
	type test struct {
		name            string
		prepare         func(mock *mocks.MockStorager)
		getJWT          func(userID string) (string, error)
		statusErrorCode codes.Code
	}
	tests := []test{
		{
			name: "Success",
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().RemoveRefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			getJWT: createJWTString,
		},
		{
			name: "Error unauthenticated",
			prepare: func(m *mocks.MockStorager) {
			},
			statusErrorCode: codes.Unauthenticated,
		},
		{
			name: "Error",
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().
					RemoveRefreshToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("some error of removing refresh token"))
			},
			getJWT:          createJWTString,
			statusErrorCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks.NewMockStorager(ctrl)

			if tt.prepare != nil {
				tt.prepare(m)
			}

			var requestCtx context.Context
			if tt.getJWT != nil {
				token, err := tt.getJWT("userID")
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				md := metadata.New(map[string]string{"Authorization": token})
				requestCtx = metadata.NewOutgoingContext(ctx, md)
			} else {
				requestCtx = ctx
			}

			usecaseAuth.Storage = m

			request := new(pb.LogoutRequest)
			request.RefreshToken = "token"

			_, err := c.Logout(requestCtx, request)

			if tt.statusErrorCode == 0 {
				require.NoError(t, err)
			} else {
				if e, ok := status.FromError(err); ok {
					require.Equal(t, tt.statusErrorCode.String(), e.Code().String())
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			}

		})
	}
}

func TestController_Login(t *testing.T) {
	ctx := context.Background()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(bootstrapAddressgRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)

	type args struct {
		login    string
		password string
	}
	type test struct {
		name            string
		args            args
		prepare         func(mock *mocks.MockStorager) error
		statusErrorCode codes.Code
	}
	tests := []test{
		{
			name: "Success",
			args: args{
				login:    "login",
				password: "password",
			},
			prepare: func(m *mocks.MockStorager) error {
				h := passwordhash.New()
				hash, err := h.Hash("password")
				if err != nil {
					return err
				}

				m.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(&smodel.GetUserByLoginResponseModel{
						UserID:    "userID",
						Password:  hash,
						SecretKey: "secretKey",
					}, nil)

				m.EXPECT().
					CreateRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				return nil
			},
		},
		{
			name: "Error login not found",
			args: args{
				login:    "login",
				password: "password",
			},
			prepare: func(m *mocks.MockStorager) error {
				m.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(nil, serror.NewLoginNotFoundError("login"))

				return nil
			},
			statusErrorCode: codes.Unauthenticated,
		},
		{
			name: "Error get user by login",
			args: args{
				login:    "login",
				password: "password",
			},
			prepare: func(m *mocks.MockStorager) error {
				m.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of getting user by login"))

				return nil
			},
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error login or password not found",
			args: args{
				login:    "login",
				password: "password",
			},
			prepare: func(m *mocks.MockStorager) error {
				m.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(&smodel.GetUserByLoginResponseModel{
						UserID:    "userID",
						Password:  "not valid hash",
						SecretKey: "secretKey",
					}, nil)

				return nil
			},
			statusErrorCode: codes.Unauthenticated,
		},
		{
			name: "Error create refresh token",
			args: args{
				login:    "login",
				password: "password",
			},
			prepare: func(m *mocks.MockStorager) error {
				h := passwordhash.New()
				hash, err := h.Hash("password")
				if err != nil {
					return err
				}

				m.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(&smodel.GetUserByLoginResponseModel{
						UserID:    "userID",
						Password:  hash,
						SecretKey: "secretKey",
					}, nil)

				m.EXPECT().
					CreateRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of creating refresh token"))

				return nil
			},
			statusErrorCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks.NewMockStorager(ctrl)

			if tt.prepare != nil {
				tt.prepare(m)
			}

			usecaseAuth.Storage = m

			request := new(pb.LoginRequest)
			request.Login = tt.args.login
			request.Password = tt.args.password

			_, err := c.Login(ctx, request)

			if tt.statusErrorCode == 0 {
				require.NoError(t, err)
			} else {
				if e, ok := status.FromError(err); ok {
					require.Equal(t, tt.statusErrorCode.String(), e.Code().String())
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			}

		})
	}
}

func TestController_RefreshToken(t *testing.T) {
	ctx := context.Background()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(bootstrapAddressgRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)

	type args struct {
		accessToken  string
		refreshToken string
	}
	type test struct {
		name            string
		args            args
		prepare         func(mock *mocks.MockStorager)
		statusErrorCode codes.Code
	}
	tests := []test{
		{
			name: "Success",
			args: args{
				accessToken:  "accessToken",
				refreshToken: "refreshToken",
			},
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().
					GetRefreshToken(gomock.Any(), gomock.Any()).
					Return(&smodel.GetRefreshTokenResponseModel{
						TokenID:   "tokenID",
						UserID:    "userID",
						ExpiredAt: time.Now().Add(time.Hour * 24),
					}, nil)

				m.EXPECT().
					UpdateRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&smodel.UpdateRefreshTokenResponseModel{
						Token:     "tokenID",
						ExpiredAt: time.Now().Add(time.Hour * 24),
					}, nil)
			},
		},
		{
			name: "Error refresh token not found",
			args: args{
				accessToken:  "accessToken",
				refreshToken: "refreshToken",
			},
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().
					GetRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil, serror.NewNotFoundError(fmt.Errorf("refresh token not found")))
			},
			statusErrorCode: codes.Unauthenticated,
		},
		{
			name: "Error can not get refresh token",
			args: args{
				accessToken:  "accessToken",
				refreshToken: "refreshToken",
			},
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().
					GetRefreshToken(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of getting refresh token"))
			},
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error refresh token expired",
			args: args{
				accessToken:  "accessToken",
				refreshToken: "refreshToken",
			},
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().
					GetRefreshToken(gomock.Any(), gomock.Any()).
					Return(&smodel.GetRefreshTokenResponseModel{
						TokenID:   "tokenID",
						UserID:    "userID",
						ExpiredAt: time.Now().Add(-1 * time.Hour * 24),
					}, nil)
			},
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error update refresh token",
			args: args{
				accessToken:  "accessToken",
				refreshToken: "refreshToken",
			},
			prepare: func(m *mocks.MockStorager) {
				m.EXPECT().
					GetRefreshToken(gomock.Any(), gomock.Any()).
					Return(&smodel.GetRefreshTokenResponseModel{
						TokenID:   "tokenID",
						UserID:    "userID",
						ExpiredAt: time.Now().Add(time.Hour * 24),
					}, nil)

				m.EXPECT().
					UpdateRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of updating refresh token"))
			},
			statusErrorCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mocks.NewMockStorager(ctrl)

			if tt.prepare != nil {
				tt.prepare(m)
			}

			usecaseAuth.Storage = m

			request := new(pb.RefreshTokenRequest)
			request.Token = new(pb.Token)
			request.Token.AccessToken = tt.args.accessToken
			request.Token.RefreshToken = tt.args.refreshToken

			_, err := c.RefreshToken(ctx, request)

			if tt.statusErrorCode == 0 {
				require.NoError(t, err)
			} else {
				if e, ok := status.FromError(err); ok {
					require.Equal(t, tt.statusErrorCode.String(), e.Code().String())
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			}

		})
	}
}
