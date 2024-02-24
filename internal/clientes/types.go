package clientes

import "time"

type CreateTransactionRequest struct {
	Amount      int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type CreateTransactionResponse struct {
	Balance int `json:"saldo"`
	Credit  int `json:"limite"`
}

type DBTransaction struct {
	Amount      int    `json:"amount"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type DBCustomer struct {
	Id      int `json:"id"`
	Balance int `json:"balance"`
	Credit  int `json:"credit"`
}

type GetBalanceBalanceResponse struct {
	Total  int       `json:"total"`
	Date   time.Time `json:"data_extrato"`
	Credit int       `json:"limite"`
}

type GetBalanceTransactionResponse struct {
	Amount      int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

type GetBalanceResponse struct {
	Balance      GetBalanceBalanceResponse       `json:"saldo"`
	Transactions []GetBalanceTransactionResponse `json:"ultimas_transacoes"`
}
