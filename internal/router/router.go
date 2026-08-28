package router

import (
	"database/sql"
	"net/http"

	"github.com/paudelchetan1112/olx-api/internal/handlers"
)

func Router(db *sql.DB) *http.ServeMux {
	lh := handlers.NewListingHandler(db)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", lh.Get)
	mux.HandleFunc("GET /listings/{id}", lh.GetOne)
	mux.HandleFunc("POST /listings/", lh.AddNew)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)
	mux.HandleFunc("PUT /listings/{id}", lh.Put)

	return mux
}
