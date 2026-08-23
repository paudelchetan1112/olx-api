package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/paudelchetan1112/olx-api/internal/models"
)

func List(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var lists []models.List
		rows, err := db.Query("SELECT id, title, description, price, city, status, created_at FROM listings ORDER BY created_at DESC")
		if err != nil {
			log.Printf("Query:%v", err)
			http.Error(w, "Error while preparing sql statement", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var list models.List
			err = rows.Scan(&list.ID, &list.Title, &list.Description, &list.Price, &list.City, &list.Status, &list.CreatedAt)
			if err != nil {
				fmt.Println("Error while scanning query")
				http.Error(w, "Error while scanning query", http.StatusInternalServerError)

			}
			lists = append(lists, list)

		}
		response := struct {
			Count int           `json:"count"`
			Data  []models.List `json:"data"`
		}{
			Count: len(lists),
			Data:  lists,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}
}

