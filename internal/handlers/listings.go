package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	uuid "github.com/google/uuid"
	"github.com/paudelchetan1112/olx-api/internal/models"
)
type ListingHandler struct{
	db *sql.DB
}
func NewListingHandler(db *sql.DB)*ListingHandler{
	return &ListingHandler{
		db: db,
	}
}
func addFilter(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	var params = map[string]string{
		"title":       "title",
		"description": "description",
		"price":       "price",
		"city":        "city",
		"status":      "status",
		"created_at":  "created_at",
	}
	v := 1
	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + "=$" + strconv.Itoa(v)
			v++
			args=append(args, value)
		}
	}
	return query, args
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}
func isValidSortField(field string) bool {
	validField := map[string]bool{
		"title":       true,
		"description": true,
		"price":       true,
		"city":        true,
		"status":      true,
		"created_at":  true,
	}
	return validField[field]
}
func addSorting(r *http.Request, query string) string {
	sortparams := r.URL.Query()["sortby"]
	fmt.Println(sortparams)
	if len(sortparams) > 0 {
		query += " ORDER BY"
		for i, param := range sortparams {
			fmt.Println(param)
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortField(field) && !isValidSortOrder(order) {
				continue
			}

			if i > 0 {
				query += ","

			}
			query += " " + field + " " + order

		}
	}
	return query
}
func (lh *ListingHandler)Get(w http.ResponseWriter, r *http.Request) {
		var lists []models.List
		query := "SELECT id, title, description, price, city, status, created_at FROM listings WHERE 1=1"
		var args []interface{}
		query, args = addFilter(r, query, args)
		query = addSorting(r, query)

		rows, err := lh.db.Query(query, args...)
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
			Status string        `json:"status"`
			Count  int           `json:"count"`
			Data   []models.List `json:"data"`
		}{
			Status: "success",
			Count:  len(lists),
			Data:   lists,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}


func (lh *ListingHandler)GetOne(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		fmt.Println(id)
		if err != nil {

			log.Printf("uuid_conversion %v", err)
			http.Error(w, "id conversion error", http.StatusBadRequest)
			return
		}
		var list models.List
		err = lh.db.QueryRow("SELECT id, title, description, price, city, status,created_at  FROM listings WHERE id=$1", id).Scan(&list.ID, &list.Title, &list.Description, &list.Price, &list.City, &list.Status, &list.CreatedAt)
		if err == sql.ErrNoRows {
			log.Printf("List not found:%v", err)
			http.Error(w, "List not found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("query:%v", err)
			http.Error(w, "List not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)

	}

func(lh *ListingHandler)AddNew(w http.ResponseWriter, r *http.Request) {
		var newLists []models.List
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&newLists)
		if err != nil {
			http.Error(w, "Invalid Request Body", http.StatusBadRequest)
			return
		}

		for _, list := range newLists {
			val := reflect.ValueOf(list)

			for i := 0; i < val.NumField(); i++ {
				fieldValue := val.Field(i)

				if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
					http.Error(w, "All field are required", http.StatusBadRequest)
					return
				}

			}

		}
		stmt, err := lh.db.Prepare(`INSERT INTO listings(title, description, price, city, status) VALUES($1, $2, $3, $4, $5) RETURNING id, created_at`)
		if err != nil {

			log.Printf("Error preparing query:%v", err)
			http.Error(w, "Error while preparing insert Query", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		addedList := make([]models.List, 0, len(newLists))
		for _, newlist := range newLists {
			err := stmt.QueryRow(newlist.Title, newlist.Description, newlist.Price, newlist.City, newlist.Status).Scan(&newlist.ID, &newlist.CreatedAt)
			if err != nil {
				log.Printf("Error while inserting into database:%v", err)
				http.Error(w, "Error inserting data into the database", http.StatusInternalServerError)

				return

			}

			addedList = append(addedList, newlist)

		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := struct {
			Status string        `json:"status"`
			Count  int           `json:"count"`
			Data   []models.List `json:"data"`
		}{
			Status: "success",
			Count:  len(addedList),
			Data:   addedList,
		}
		json.NewEncoder(w).Encode(response)

	}


func(lh *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		fmt.Println(id)
		if err != nil {

			log.Printf("something went wrong%v", err)
			http.Error(w, "something went wrong", http.StatusBadRequest)
			return
		}

		result, err := lh.db.Exec("DELETE FROM listings WHERE id=$1", id)
		if  err != nil {
			log.Printf("Delete:%v", err)
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
		response := struct {
			Status string    `json:"status"`
			// Id     uuid.UUID `json:"id"`
		}{
			Status: "List Successfully Deleted",
			// Id:     id,
		}
		json.NewEncoder(w).Encode(response)

	}


func(lh *ListingHandler) DeleteMultipleList(w http.ResponseWriter, r *http.Request) {
		var ids []uuid.UUID
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			log.Printf("Error: decoding json data, %v ", err)

			http.Error(w, "Error while decoding json data", http.StatusInternalServerError)
			return
		}
		tx, err := lh.db.Begin()
		if err != nil {
			log.Printf("Error while begin transaction, %v", err)
		}
		stmt, err := tx.Prepare("DELETE FROM listings WHERE id=$1")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error while preparing SQL query", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		deletedIds := []uuid.UUID{}
		for _, id := range ids {
			result, err := stmt.Exec(id)
			if err != nil {
				tx.Rollback()
				log.Printf("Error No Row:,%v ", err)
				http.Error(w, "Error while deleting list", http.StatusNotFound)
				return
			}
			if err != nil {
				http.Error(w, "Error while retrieving deleted data", http.StatusInternalServerError)
				return
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return
			}

			if rowsAffected > 0 {
				deletedIds = append(deletedIds, id)
			}
			if rowsAffected < 1 {
				tx.Rollback()
				http.Error(w, fmt.Sprintf("Id %d doesn't exist", id), http.StatusNotFound)

				return

			}
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			http.Error(w, "Error while commmiting transcation", http.StatusInternalServerError)

		}

	}


