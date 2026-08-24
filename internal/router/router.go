package router

import (
	"database/sql"
	"net/http"

	"github.com/paudelchetan1112/olx-api/internal/handlers"
)

func Router(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", handlers.List(db))
	mux.HandleFunc("GET /listings/{id}", handlers.GetOneList(db))
	mux.HandleFunc("POST /listings/", handlers.AddList(db))
	mux.HandleFunc("DELETE /listings/{id}", handlers.DeleteOneList(db))

	return mux
}
