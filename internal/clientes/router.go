package clientes

import "github.com/go-chi/chi/v5"

func AddRoutes(r chi.Router, h Handler) {
	r.Route("/clientes/{customerId}", func(r chi.Router) {
		r.Post("/transacoes", h.CreateTransaction)
		r.Get("/extrato", h.GetBalance)
	})
}
