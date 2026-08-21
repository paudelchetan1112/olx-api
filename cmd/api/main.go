package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/paudelchetan1112/olx-api/internal/config"
)

func main() {
	cfg:=config.MustLoad()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Status":"Okey"}`))
	})
		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Status":"Welcome to the home route"}`))
	})
	srv := http.Server{
		Addr:        ":"+cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed:%v", err)
	}

	fmt.Println("Server is running..")

}
