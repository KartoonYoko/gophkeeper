package clientstore

type BankCardDataModel struct {
	Number string `json:"number"`
	CVV    string `json:"cvv"`
}

type CredentialDataModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
