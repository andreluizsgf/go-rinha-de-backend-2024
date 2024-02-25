package clientes

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	if !checkCustomerExists(customerId) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// decode request body to transaction struct
	var transaction CreateTransactionRequest
	json.NewDecoder(r.Body).Decode(&transaction)

	// insert transaction into database
	_, err := h.db.Exec("INSERT INTO transaction (customer_id, amount, type, description) VALUES ($1, $2, $3, $4)", customerId, transaction.Amount, transaction.Type, transaction.Description)

	if err != nil {
		var errMsg = err.Error()

		if errMsg == "pq: no limit" {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		} else if errMsg == "pq: customer not found" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

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
	customerId, _ := strconv.Atoi(chi.URLParam(r, "customerId"))

	fmt.Println(checkCustomerExists(customerId))

	if !checkCustomerExists(customerId) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// get customer balance from database
	var dbCustomer DBCustomer

	err := h.db.QueryRow("SELECT id, balance, credit FROM customer WHERE id = $1", customerId).Scan(&dbCustomer.Id, &dbCustomer.Balance, &dbCustomer.Credit)

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

func checkCustomerExists(id int) bool {
	return id >= 1 && id <= 5
}
