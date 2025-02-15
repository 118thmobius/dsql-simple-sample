package domain

type Transaction struct {
	FromId string `json:"from_id"`
	ToId   string `json:"to_id"`
	Amount int    `json:"amount"`
}

func NewTransaction(from Account, to Account, amount int) *Transaction {
	return &Transaction{
		FromId: from.UserId,
		ToId:   to.UserId,
		Amount: amount,
	}
}
