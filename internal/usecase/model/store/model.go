package store

type SaveDataRequestModel struct {
	UserID      string
	Data        []byte
	DataType    string
	Description string
}

type GetDataByIDRequestModel struct {
	UserID string
	ID     string
}

type GetDataByIDResponseModel struct {
	Data        []byte
	DataType    string
	Description string
}
