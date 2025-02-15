package domain

type Account struct {
	UserId  string `json:"user_id"`
	City    string `json:"city"`
	Balance int    `json:"balance"`
}
