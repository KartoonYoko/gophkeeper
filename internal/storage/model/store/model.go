package store

type SaveDataRequestModel struct {
	UserID      string
	BinaryID    string
	Description string
	DataType    string
}

type SaveDataResponseModel struct {
	ID int
}

type GetDataByIDRequestModel struct {
	UserID string
	ID     int
}

type GetDataByIDResponseModel struct {
	ID          int
	UserID      string
	BinaryID    string
	Description string
	DataType    string
}

type RemoveDataByIDRequestModel struct {
	UserID string
	ID     string
}
