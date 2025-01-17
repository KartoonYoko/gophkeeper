package store

import commondatatype "github.com/KartoonYoko/gophkeeper/internal/common/datatype"

type SaveDataRequestModel struct {
	ID                    string
	UserID                string
	Data                  []byte
	DataType              DataType
	Description           string
	ModificationTimestamp int64
	Hash                  string
}

type SaveDataResponseModel struct {
	DataID string
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

type UpdateDataRequestModel struct {
	ID                    string
	UserID                string
	Hash                  string
	ModificationTimestamp int64
	Data                  []byte
}

type RemoveDataByIDRequestModel struct {
	ID                    string
	UserID                string
	ModificationTimestamp int64
}

type RemoveDataByIDResponseModel struct{}

type UpdateDataResponseModel struct{}

// Тип хранимых данных. Возможные значения:
//   - BANK_CARD
//   - BINARY
//   - CREDENTIALS
//   - TEXT
type DataType string

func (dt *DataType) IsValid() bool {
	switch string(*dt) {
	case commondatatype.DATATYPE_TEXT:
		return true
	case commondatatype.DATATYPE_CREDENTIALS:
		return true
	case commondatatype.DATATYPE_BANK_CARD:
		return true
	case commondatatype.DATATYPE_BINARY:
		return true
	}

	return false
}

func (dt *DataType) String() string {
	return string(*dt)
}
