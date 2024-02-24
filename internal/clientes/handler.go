package clientes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler interface {
	CreateTransaction(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db sql.DB
}

func NewHandler(db *sql.DB) Handler {
	return &handler{}
}

type Transaction struct {
	Amount      int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type Customer struct {
	Id      int           `json:"id"`
	Balance int           `json:"balance"`
	Credit  int           `json:"credit"`
	History []Transaction `json:"history"`
}

func (h *handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// decode request body to transaction struct
	var transaction Transaction
	json.NewDecoder(r.Body).Decode(&transaction)

	// connect to db
	db, err := sql.Open("postgres", "user=gorinha dbname=gorinha password=gorinha sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}

	// check if its debit
	if transaction.Type == "d" {
		transaction.Amount = -transaction.Amount
	}

	// get customer balance from database
	var balance int
	var credit int

	err = db.QueryRow("SELECT balance, credit FROM customer WHERE id = $1", 1).Scan(&balance, &credit)

	fmt.Println(balance, transaction.Amount, credit, balance+transaction.Amount+credit)
	if err != nil {
		fmt.Println(err)
	}

	// check if customer has enough balance to make transaction
	if balance+transaction.Amount+credit < 0 {
		fmt.Println("saldo insuficiente")
	}

	// insert transaction into database
	_, err = db.Exec("INSERT INTO transaction (customer_id, amount, type, description) VALUES (1, $1, $2, $3)", transaction.Amount, transaction.Type, transaction.Description)

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte("welcome"))
}
