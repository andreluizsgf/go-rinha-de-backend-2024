package clientes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	CreateTransaction(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) Handler {
	return &handler{db: db}
}

func (h *handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	customerId, _ := strconv.Atoi(chi.URLParam(r, "customerId"))

	// decode request body to transaction struct
	var transaction CreateTransactionRequest
	json.NewDecoder(r.Body).Decode(&transaction)

	// check if its debit
	if transaction.Type == "d" {
		transaction.Amount = -transaction.Amount
	}

	// get customer balance from database
	var dbCustomer DBCustomer

	err := h.db.QueryRow("SELECT id, balance, credit FROM customer WHERE id = $1", customerId).Scan(&dbCustomer.Id, &dbCustomer.Balance, &dbCustomer.Credit)

	// check if customer exists
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// check if customer has enough balance to make transaction
	if dbCustomer.Balance+transaction.Amount+dbCustomer.Credit < 0 {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	// insert transaction into database
	_, err = h.db.Exec("INSERT INTO transaction (customer_id, amount, type, description) VALUES (1, $1, $2, $3)", transaction.Amount, transaction.Type, transaction.Description)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// get customer with balance updated balance
	var updatedCustomer DBCustomer
	err = h.db.QueryRow("SELECT balance, credit FROM customer WHERE id = $1", customerId).Scan(&updatedCustomer.Balance, &updatedCustomer.Credit)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := CreateTransactionResponse{Balance: updatedCustomer.Balance, Credit: updatedCustomer.Credit}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	customerId := chi.URLParam(r, "customerId")

	// get customer balance from database
	var dbCustomer DBCustomer

	err := h.db.QueryRow("SELECT id, balance, credit FROM customer WHERE id = $1", customerId).Scan(&dbCustomer.Id, &dbCustomer.Balance, &dbCustomer.Credit)

	// check if customer exists
	if err != nil && err == sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// get last transactions
	rows, err := h.db.Query("SELECT amount, type, description, created_at FROM transaction WHERE customer_id = $1 ORDER BY created_at DESC LIMIT 10", customerId)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := GetBalanceResponse{Balance: GetBalanceBalanceResponse{Total: dbCustomer.Balance, Date: time.Now().UTC(), Credit: dbCustomer.Credit}}

	for rows.Next() {
		var transaction GetBalanceTransactionResponse

		err = rows.Scan(&transaction.Amount, &transaction.Type, &transaction.Description,
			&transaction.CreatedAt)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response.Transactions = append(response.Transactions, transaction)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
