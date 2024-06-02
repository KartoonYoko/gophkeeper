package store

type SaveDataRequestModel struct {
	UserID      string
	Data        []byte
	DataType    DataType
	Description string
}

type SaveDataResponseModel struct {
	DataID int
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

// Тип хранимых данных. Возможные значения:
//   - BANK_CARD
//   - BINARY
//   - CREDENTIALS
//   - TEXT
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

func (dt *DataType) String() string {
	return string(*dt)
}
