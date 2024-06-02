package store

// Тип хранимых данных. Возможные значения:
// 	- BANK_CARD
// 	- BINARY
// 	- CREDENTIALS
// 	- TEXT
type DataType string

func (dt *DataType) IsValid() bool {
	switch string(*dt) {
	case "TEXT":
		return true
	case "CREDENTIALS":
		return true
	case "BANK_CARD":
		return true
	case "BINARY":
		return true
	}

	return false
}

type SaveDataRequestModel struct {
	UserID      string
	BinaryID    string
	Description string
	DataType    DataType
}

type SaveDataResponseModel struct {
}

type GetDataByIDRequestModel struct {
	UserID string
	ID     string
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
