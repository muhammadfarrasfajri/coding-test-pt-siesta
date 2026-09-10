package model

type Account struct {
	ID      string
	Balance float64
}

type TransferRequest struct {
	FromAccount string  `json:"from_account" binding:"required"`
	ToAccount   string  `json:"to_account" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}
