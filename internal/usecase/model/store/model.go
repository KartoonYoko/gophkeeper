package store

type SaveDataRequestModel struct {
	UserID      string
	Data        []byte
	DataType    string
	Description string
}
