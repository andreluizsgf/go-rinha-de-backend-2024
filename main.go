package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"example/rinha-de-backend-2024/internal/clientes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
)

func main() {
	// setup
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", os.Getenv("PG_HOST"), os.Getenv("PG_USER"), os.Getenv("PG_DB"), os.Getenv("PG_PASSWORD"))

	db, err := sql.Open("postgres", dsn)

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
