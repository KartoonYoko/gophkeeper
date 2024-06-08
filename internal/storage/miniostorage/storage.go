package miniostorage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"

	filestoremodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/filestore"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	conf Config

	client *minio.Client
}

func NewStorage(conf Config) (*Storage, error) {
	var err error
	s := new(Storage)

	s.conf = conf

	s.client, err = minio.New(s.conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.conf.AccessKeyID, s.conf.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) SaveData(ctx context.Context, request *filestoremodel.SaveDataRequestModel) (*filestoremodel.SaveDataResponseModel, error) {
	contentType := "application/octet-stream"

	bucketName := request.UserID
	isExists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !isExists {
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	objName, err := s.generateDataKey()
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(request.Data)
	_, err = s.client.PutObject(
		ctx,
		bucketName,
		objName,
		reader,
		int64(len(request.Data)),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}

	response := new(filestoremodel.SaveDataResponseModel)
	response.ID = objName

	return response, nil
}

func (s *Storage) GetDataByID(ctx context.Context, request *filestoremodel.GetDataByIDRequestModel) (*filestoremodel.GetDataByIDResponseModel, error) {
	obj, err := s.client.GetObject(ctx, request.UserID, request.ID, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	response := new(filestoremodel.GetDataByIDResponseModel)
	response.Data, err = io.ReadAll(obj)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Storage) RemoveDataByID(ctx context.Context, request *filestoremodel.RemoveDataByIDRequestModel) error {
	return s.client.RemoveObject(ctx, request.UserID, request.ID, minio.RemoveObjectOptions{})
}

func (s *Storage) generateDataKey() (string, error) {
	n := 16
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ``, err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
