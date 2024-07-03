package grpcserver

import (
	"context"
	"fmt"

	"github.com/KartoonYoko/gophkeeper/internal/logger"
	pb "github.com/KartoonYoko/gophkeeper/internal/proto"
	ucmodel "github.com/KartoonYoko/gophkeeper/internal/usecase/model/store"
	commondatatype "github.com/KartoonYoko/gophkeeper/internal/common/datatype"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SaveData сохранит данные в хранилище
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
		UserID:                userID,
		Data:                  r.Data,
		DataType:              dt,
		Description:           r.Description,
		ID:                    r.Id,
		Hash:                  r.Hash,
		ModificationTimestamp: r.ModificationTimestamp,
	}
	saveDataResponse, err := c.usecaseStore.SaveData(ctx, ucrequest)
	if err != nil {
		logger.Log.Error("unable save data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	response := new(pb.SaveDataResponse)
	response.DataId = saveDataResponse.DataID

	return response, nil
}

// GetDataByID возвращает данные из хранилища по ID
func (c *Controller) GetDataByID(ctx context.Context, r *pb.GetDataByIDRequest) (*pb.GetDataByIDResponse, error) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		logger.Log.Error("can not get user ID from context", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	request := &ucmodel.GetDataByIDRequestModel{
		UserID: userID,
		ID:     r.Id,
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

// UpdateData обновляет данные в хранилище
func (c *Controller) UpdateData(ctx context.Context, r *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		logger.Log.Error("can not get user ID from context", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	_, err = c.usecaseStore.UpdateData(ctx, &ucmodel.UpdateDataRequestModel{
		ID:                    r.Id,
		Data:                  r.Data,
		UserID:                userID,
		Hash:                  r.Hash,
		ModificationTimestamp: r.ModificationTimestamp,
	})

	if err != nil {
		logger.Log.Error("unable update data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.UpdateDataResponse{}, nil
}

// RemoveData помечает данные в хранилище удалёнными
func (c *Controller) RemoveData(ctx context.Context, r *pb.RemoveDataRequest) (*pb.RemoveDataResponse, error) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		logger.Log.Error("can not get user ID from context", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	_, err = c.usecaseStore.RemoveDataByID(ctx, &ucmodel.RemoveDataByIDRequestModel{
		ID:                    r.Id,
		UserID:                userID,
		ModificationTimestamp: r.ModificationTimestamp,
	})

	if err != nil {
		logger.Log.Error("remove data error", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.RemoveDataResponse{}, nil
}

// GetMetaDataList возвращает список метаданных
func (c *Controller) GetMetaDataList(ctx context.Context, r *pb.GetMetaDataListRequest) (*pb.GetMetaDataListResponse, error) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		logger.Log.Error("can not get user ID from context", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	rs, err := c.usecaseStore.GetUserDataList(ctx, userID)
	if err != nil {
		logger.Log.Error("unable get user data list", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	response := new(pb.GetMetaDataListResponse)
	for _, item := range rs.Items {
		dt, err := getProtoDataTypeFromUsecaseDataType(item.DataType)
		if err != nil {
			logger.Log.Error("wrong data type", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "internal error")
		}

		response.Items = append(response.Items, &pb.GetMetaDataListItemResponse{
			Id:                    item.ID,
			Type:                  dt,
			Hash:                  item.Hash,
			Description:           item.Description,
			ModificationTimestamp: item.ModificationTimestamp,
			IsDeleted:             item.IsDeleted,
		})
	}

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
	case commondatatype.DATATYPE_TEXT:
		return pb.DataTypeEnum_DATA_TYPE_TEXT, nil
	case commondatatype.DATATYPE_BINARY:
		return pb.DataTypeEnum_DATA_TYPE_BINARY, nil
	case commondatatype.DATATYPE_BANK_CARD:
		return pb.DataTypeEnum_DATA_TYPE_BANK_CARD, nil
	case commondatatype.DATATYPE_CREDENTIALS:
		return pb.DataTypeEnum_DATA_TYPE_CREDENTIALS, nil
	}

	return pb.DataTypeEnum(0), fmt.Errorf("invalid data type %s", dataType)
}
