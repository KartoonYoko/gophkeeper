package grpcserver

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
	sfmodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/filestore"
	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/store"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/store/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestController_SaveData(t *testing.T) {
	ctx := context.Background()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(bootstrapAddressgRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewStoreServiceClient(conn)

	type args struct {
		dataType pb.DataTypeEnum
	}
	type test struct {
		name            string
		args            args
		prepare         func(mock *mocks.MockStorager, mf *mocks.MockFileStorager)
		getJWT          func(userID string) (string, error)
		statusErrorCode codes.Code
	}
	tests := []test{
		{
			name: "Success text",
			args: args{
				dataType: pb.DataTypeEnum_DATA_TYPE_TEXT,
			},
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) {
				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&sfmodel.SaveDataResponseModel{}, nil)
				m.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&smodel.SaveDataResponseModel{}, nil)
			},
			getJWT: createJWTString,
		},
		{
			name: "Success binary",
			args: args{
				dataType: pb.DataTypeEnum_DATA_TYPE_BINARY,
			},
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) {
				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&sfmodel.SaveDataResponseModel{}, nil)
				m.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&smodel.SaveDataResponseModel{}, nil)
			},
			getJWT: createJWTString,
		},
		{
			name: "Success bank card",
			args: args{
				dataType: pb.DataTypeEnum_DATA_TYPE_BANK_CARD,
			},
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) {
				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&sfmodel.SaveDataResponseModel{}, nil)
				m.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&smodel.SaveDataResponseModel{}, nil)
			},
			getJWT: createJWTString,
		},
		{
			name: "Success credentials",
			args: args{
				dataType: pb.DataTypeEnum_DATA_TYPE_CREDENTIALS,
			},
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) {
				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&sfmodel.SaveDataResponseModel{}, nil)
				m.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&smodel.SaveDataResponseModel{}, nil)
			},
			getJWT: createJWTString,
		},
		{
			name: "Error file storeage save",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) {
				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of saving file"))
			},
			getJWT:          createJWTString,
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error storage save",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) {
				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&sfmodel.SaveDataResponseModel{}, nil)
				m.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of saving data"))
				mf.EXPECT().
					RemoveDataByID(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			getJWT:          createJWTString,
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error storage save and delete file",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) {
				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&sfmodel.SaveDataResponseModel{}, nil)
				m.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of saving data"))
				mf.EXPECT().
					RemoveDataByID(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("can not delete saved file"))
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
			mf := mocks.NewMockFileStorager(ctrl)

			if tt.prepare != nil {
				tt.prepare(m, mf)
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

			usecaseStore.FileStorage = mf
			usecaseStore.Storage = m

			request := new(pb.SaveDataRequest)
			request.Type = tt.args.dataType

			_, err := c.SaveData(requestCtx, request)

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

func TestController_GetDataByID(t *testing.T) {
	ctx := context.Background()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(bootstrapAddressgRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewStoreServiceClient(conn)

	type test struct {
		name            string
		prepare         func(mock *mocks.MockStorager, mf *mocks.MockFileStorager) error
		getJWT          func(userID string) (string, error)
		statusErrorCode codes.Code
	}
	tests := []test{
		{
			name: "Success",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				d, err := encrypteData([]byte("some data"))
				if err != nil {
					return err
				}
				m.EXPECT().
					GetDataByID(gomock.Any(), gomock.Any()).
					Return(&smodel.GetDataByIDResponseModel{
						DataType: "TEXT",
					}, nil)
				mf.EXPECT().
					GetDataByID(gomock.Any(), gomock.Any()).
					Return(&sfmodel.GetDataByIDResponseModel{
						Data: d,
					}, nil)

				return nil
			},
			getJWT: createJWTString,
		},
		{
			name: "Error wrong data type",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				d, err := encrypteData([]byte("some data"))
				if err != nil {
					return err
				}
				m.EXPECT().
					GetDataByID(gomock.Any(), gomock.Any()).
					Return(&smodel.GetDataByIDResponseModel{
						DataType: "WRONG DATA TYPE",
					}, nil)
				mf.EXPECT().
					GetDataByID(gomock.Any(), gomock.Any()).
					Return(&sfmodel.GetDataByIDResponseModel{
						Data: d,
					}, nil)

				return nil
			},
			getJWT:          createJWTString,
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error get storage data",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					GetDataByID(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of getting data"))
				return nil
			},
			getJWT:          createJWTString,
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error get file storage data",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					GetDataByID(gomock.Any(), gomock.Any()).
					Return(&smodel.GetDataByIDResponseModel{}, nil)
				mf.EXPECT().
					GetDataByID(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of getting file"))
				return nil
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
			mf := mocks.NewMockFileStorager(ctrl)

			if tt.prepare != nil {
				err = tt.prepare(m, mf)
				require.NoError(t, err)
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

			usecaseStore.FileStorage = mf
			usecaseStore.Storage = m

			request := new(pb.GetDataByIDRequest)
			request.Id = "dataID"

			_, err := c.GetDataByID(requestCtx, request)

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

func TestController_UpdateData(t *testing.T) {
	ctx := context.Background()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(bootstrapAddressgRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewStoreServiceClient(conn)

	type test struct {
		name            string
		prepare         func(mock *mocks.MockStorager, mf *mocks.MockFileStorager) error
		getJWT          func(userID string) (string, error)
		statusErrorCode codes.Code
	}

	tests := []test{
		{
			name: "Success",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {

				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&sfmodel.SaveDataResponseModel{
						ID:     "dataID",
						UserID: "UserID",
					}, nil)

				m.EXPECT().
					UpdateData(gomock.Any(), gomock.Any()).
					Return(&smodel.UpdateDataResponseModel{}, nil)

				return nil
			},
			getJWT: createJWTString,
		},
		{
			name: "Error save file",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {

				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of save data"))

				return nil
			},
			getJWT:          createJWTString,
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error update data",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {

				mf.EXPECT().
					SaveData(gomock.Any(), gomock.Any()).
					Return(&sfmodel.SaveDataResponseModel{
						ID:     "dataID",
						UserID: "UserID",
					}, nil)

				m.EXPECT().
					UpdateData(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error of update data"))

				return nil
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
			mf := mocks.NewMockFileStorager(ctrl)

			if tt.prepare != nil {
				err = tt.prepare(m, mf)
				require.NoError(t, err)
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

			usecaseStore.FileStorage = mf
			usecaseStore.Storage = m

			request := new(pb.UpdateDataRequest)
			request.Id = "dataID"
			request.Data = []byte("some data")
			request.Hash = "some hash"

			_, err := c.UpdateData(requestCtx, request)

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

func TestController_RemoveData(t *testing.T) {
	ctx := context.Background()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(bootstrapAddressgRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewStoreServiceClient(conn)

	type test struct {
		name            string
		prepare         func(mock *mocks.MockStorager, mf *mocks.MockFileStorager) error
		getJWT          func(userID string) (string, error)
		statusErrorCode codes.Code
	}

	tests := []test{
		{
			name: "Success",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					RemoveDataByID(gomock.Any(), gomock.Any()).
					Return(nil)

				mf.EXPECT().
					RemoveDataByID(gomock.Any(), gomock.Any()).
					Return(nil)

				return nil
			},
			getJWT: createJWTString,
		},
		{
			name: "Error remove data",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					RemoveDataByID(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("some error of remove data"))

				return nil
			},
			getJWT:          createJWTString,
			statusErrorCode: codes.Internal,
		},
		{
			name: "Error file remove",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					RemoveDataByID(gomock.Any(), gomock.Any()).
					Return(nil)

				mf.EXPECT().
					RemoveDataByID(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("some error of file remove"))

				return nil
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
			mf := mocks.NewMockFileStorager(ctrl)

			if tt.prepare != nil {
				err = tt.prepare(m, mf)
				require.NoError(t, err)
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

			usecaseStore.FileStorage = mf
			usecaseStore.Storage = m

			request := new(pb.RemoveDataRequest)
			request.Id = "dataID"

			_, err := c.RemoveData(requestCtx, request)

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

func TestController_GetMetaDataList(t *testing.T) {
	ctx := context.Background()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(bootstrapAddressgRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewStoreServiceClient(conn)

	type test struct {
		name            string
		prepare         func(mock *mocks.MockStorager, mf *mocks.MockFileStorager) error
		getJWT          func(userID string) (string, error)
		statusErrorCode codes.Code
	}

	tests := []test{
		{
			name: "Success text",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					GetUserDataList(gomock.Any(), gomock.Any()).
					Return(&smodel.GetUserDataListResponseModel{
						Items: []*smodel.GetUserDataListResponseItemModel{
							{
								DataType: "TEXT",
							},
						},
					}, nil)

				return nil
			},
			getJWT: createJWTString,
		},
		{
			name: "Success binary",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					GetUserDataList(gomock.Any(), gomock.Any()).
					Return(&smodel.GetUserDataListResponseModel{
						Items: []*smodel.GetUserDataListResponseItemModel{
							{
								DataType: "BINARY",
							},
						},
					}, nil)

				return nil
			},
			getJWT: createJWTString,
		},
		{
			name: "Success credentials",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					GetUserDataList(gomock.Any(), gomock.Any()).
					Return(&smodel.GetUserDataListResponseModel{
						Items: []*smodel.GetUserDataListResponseItemModel{
							{
								DataType: "CREDENTIALS",
							},
						},
					}, nil)

				return nil
			},
			getJWT: createJWTString,
		},
		{
			name: "Success bank card",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					GetUserDataList(gomock.Any(), gomock.Any()).
					Return(&smodel.GetUserDataListResponseModel{
						Items: []*smodel.GetUserDataListResponseItemModel{
							{
								DataType: "BANK_CARD",
							},
						},
					}, nil)

				return nil
			},
			getJWT: createJWTString,
		},
		{
			name: "Error wrong data type",
			prepare: func(m *mocks.MockStorager, mf *mocks.MockFileStorager) error {
				m.EXPECT().
					GetUserDataList(gomock.Any(), gomock.Any()).
					Return(&smodel.GetUserDataListResponseModel{
						Items: []*smodel.GetUserDataListResponseItemModel{
							{
								DataType: "WRONG DATA TYPE",
							},
						},
					}, nil)

				return nil
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
			mf := mocks.NewMockFileStorager(ctrl)

			if tt.prepare != nil {
				err = tt.prepare(m, mf)
				require.NoError(t, err)
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

			usecaseStore.FileStorage = mf
			usecaseStore.Storage = m

			request := new(pb.GetMetaDataListRequest)
			_, err := c.GetMetaDataList(requestCtx, request)

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
