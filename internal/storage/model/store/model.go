package store

type SaveDataRequestModel struct {
	ID                    string
	UserID                string
	BinaryID              string
	Description           string
	DataType              string
	Hash                  string
	ModificationTimestamp int64
}

type SaveDataResponseModel struct {
	ID                    string
	UserID                string
	Hash                  string
	ModificationTimestamp int64
}

type UpdateDataResponseModel struct{}

type UpdateDataRequestModel struct {
	ID                    string
	UserID                string
	Hash                  string
	ModificationTimestamp int64
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
	UserID                string
	ID                    string
	ModificationTimestamp int64
}

type GetUserDataListRequestModel struct {
	UserID string
}

type GetUserDataListResponseModel struct {
	Items []*GetUserDataListResponseItemModel
}

type GetUserDataListResponseItemModel struct {
	ID                    string
	Description           string
	DataType              string
	Hash                  string
	ModificationTimestamp int64
	IsDeleted             bool
}
