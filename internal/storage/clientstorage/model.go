package clientstorage

type tokensFile struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	SecretKey    string `json:"secret_key"`
}
