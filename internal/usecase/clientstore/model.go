package clientstore

import "time"

type BankCardDataModel struct {
	Number     string    `json:"number"`
	CVV        string    `json:"cvv"`
	ValidUntil time.Time `json:"valid_until"`
}

type CredentialDataModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
