package models

type UserReristerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BalanceWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}
