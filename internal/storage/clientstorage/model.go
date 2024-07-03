package clientstorage

type SaveDataRequestModel struct {
	ID          string
	Userid      string
	Description string
	// Возможные значения:
	// 	- TEXT
	// 	- BINARY
	// 	- CREDENTIALS
	// 	- BANK_CARD
	Datatype              string
	Hash                  string
	ModificationTimestamp int64
	Data                  []byte
}

type UpdateDataRequestModel struct {
	ID                    string
	Hash                  string
	ModificationTimestamp int64
	Data                  []byte
}

type GetDataListResponseItemModel struct {
	ID          string
	UserID      string
	Description string
	Datatype    string
}

type GetDataListToSynchronizeItemModel struct {
	ID                    string
	UserID                string
	Description           string
	Datatype              string
	Hash                  string
	IsDeleted             bool
	ModificationTimestamp int64
}

type GetDataByIDResponseModel struct {
	ID                    string
	Userid                string
	Description           string
	Datatype              string
	Hash                  string
	ModificationTimestamp int64
	Data                  []byte
}

type GetTextDataByIDResponseModel struct {
	Description string
	Text        string
}

type GetBinaryDataByIDResponseModel struct {
	Description string
	Text        string
}

type GetBankCardDataByIDResponseModel struct {
	Description string
	Number      string
	CVV         string
}

type GetCredentialsDataByIDResponseModel struct {
	Description string
	Login       string
	Password    string
}

type RemoveDataByIDRequestModel struct {
	DataID                string
	ModificationTimestamp int64
}

type credentialsFile struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	SecretKey    string `json:"secret_key"`
	UserID       string `json:"user_id"`
}
