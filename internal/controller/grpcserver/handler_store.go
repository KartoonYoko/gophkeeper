package grpcserver

import (
	"context"
	"fmt"

	"github.com/KartoonYoko/gophkeeper/internal/logger"
	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
	ucmodel "github.com/KartoonYoko/gophkeeper/internal/usecase/model/store"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) SaveData(ctx context.Context, r *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		logger.Log.Error("can not get user ID from context", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	dt, err := getUsecaseDataTypeFromProtoDataType(r.Type)
	if err != nil {
		logger.Log.Error("unable cast data type", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	ucrequest := &ucmodel.SaveDataRequestModel{
		UserID:      userID,
		Data:        r.Data,
		DataType:    dt,
		Description: r.Description,
	}
	saveDataResponse, err := c.usecaseStore.SaveData(ctx, ucrequest)
	if err != nil {
		logger.Log.Error("unable save data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	response := new(pb.SaveDataResponse)
	response.DataId = int32(saveDataResponse.DataID)

	return response, nil
}

func (c *Controller) GetDataByID(ctx context.Context, r *pb.GetDataByIDRequest) (*pb.GetDataByIDResponse, error) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		logger.Log.Error("can not get user ID from context", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	
	request := &ucmodel.GetDataByIDRequestModel{
		UserID: userID,
		ID: int(r.Id),
	}
	getDataResponse, err := c.usecaseStore.GetDataByID(ctx, request)
	if err != nil {
		logger.Log.Error("unable get data by id", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	dt, err := getProtoDataTypeFromUsecaseDataType(getDataResponse.DataType)
	if err != nil {
		logger.Log.Error("wrong data type", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	response := new(pb.GetDataByIDResponse)
	response.Data = getDataResponse.Data
	response.Description = getDataResponse.Description
	response.Type = dt

	return response, nil
}	

func getUsecaseDataTypeFromProtoDataType(dataType pb.DataTypeEnum) (ucmodel.DataType, error) {
	switch dataType {
	case pb.DataTypeEnum_DATA_TYPE_TEXT:
		return ucmodel.DataType("TEXT"), nil
	case pb.DataTypeEnum_DATA_TYPE_BINARY:
		return ucmodel.DataType("BINARY"), nil
	case pb.DataTypeEnum_DATA_TYPE_BANK_CARD:
		return ucmodel.DataType("BANK_CARD"), nil
	case pb.DataTypeEnum_DATA_TYPE_CREDENTIALS:
		return ucmodel.DataType("CREDENTIALS"), nil
	}

	return ucmodel.DataType(""), fmt.Errorf("invalid data type %s", dataType)
}

func getProtoDataTypeFromUsecaseDataType(dataType string) (pb.DataTypeEnum, error) {
	switch dataType {
	case "TEXT":
		return pb.DataTypeEnum_DATA_TYPE_TEXT, nil
	case "BINARY":
		return pb.DataTypeEnum_DATA_TYPE_BINARY, nil
	case "BANK_CARD":
		return pb.DataTypeEnum_DATA_TYPE_BANK_CARD, nil
	case "CREDENTIALS":
		return pb.DataTypeEnum_DATA_TYPE_CREDENTIALS, nil
	}

	return pb.DataTypeEnum(0), fmt.Errorf("invalid data type %s", dataType)
}