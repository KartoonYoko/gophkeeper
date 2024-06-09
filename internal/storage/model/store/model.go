package store

type SaveDataRequestModel struct {
	ID          string
	UserID      string
	BinaryID    string
	Description string
	DataType    string
}

type SaveDataResponseModel struct {
	ID string
}

type GetDataByIDRequestModel struct {
	UserID string
	ID     string
}

type GetDataByIDResponseModel struct {
	ID          string
	UserID      string
	BinaryID    string
	Description string
	DataType    string
}

type RemoveDataByIDRequestModel struct {
	UserID string
	ID     string
}
