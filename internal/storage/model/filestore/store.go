package filestore

type SaveDataResponseModel struct {
	ID     string
	UserID string
}

type SaveDataRequestModel struct {
	ID     string
	Data   []byte
	UserID string
}

type GetDataByIDRequestModel struct {
	UserID string
	ID     string
}

type GetDataByIDResponseModel struct {
	Data []byte
}

type RemoveDataByIDRequestModel struct {
	UserID string
	ID     string
}
