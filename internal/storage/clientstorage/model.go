package clientstorage

type SaveDataRequestModel struct {
	Filename    string
	Userid      string
	Description string
	Datatype    string
	Data        []byte
}

type credentialsFile struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	SecretKey    string `json:"secret_key"`
	UserID       string `json:"user_id"`
}
