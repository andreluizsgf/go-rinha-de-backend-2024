package main

import (
	"database/sql"
	"net/http"

	"example/rinha-de-backend-2024/internal/clientes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
)

func main() {
	// setup
	db, err := sql.Open("postgres", "user=gorinha dbname=gorinha password=gorinha sslmode=disable")

	if err != nil {
		panic(err)
	}

	// routing
	r := chi.NewRouter()
	h := clientes.NewHandler(db)

	r.Use(middleware.Logger)

	clientes.AddRoutes(r, h)

	// server
	http.ListenAndServe(":3000", r)
}
