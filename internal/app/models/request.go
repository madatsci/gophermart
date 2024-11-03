package models

type UserReristerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
