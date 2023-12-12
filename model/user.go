package model

type UserSign struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserJwt struct {
	Account string `json:"account"`
}
