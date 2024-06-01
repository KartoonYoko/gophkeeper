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
	BinaryID    string
	Description string
	DataType DataType
}

type SaveDataResponseModel struct {
}
