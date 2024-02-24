package clientes

import "github.com/go-chi/chi/v5"

func AddRoutes(r chi.Router, h Handler) {
	r.Post("/", h.CreateTransaction)
}
