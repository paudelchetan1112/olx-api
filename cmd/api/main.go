package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/paudelchetan1112/olx-api/internal/config"
	"github.com/paudelchetan1112/olx-api/internal/db"
	"github.com/paudelchetan1112/olx-api/internal/handlers"
)

func main() {
	cfg := config.MustLoad() //package which config env variable 
	_, err:=db.Connect(cfg.DatabaseUrl)
	if err!=nil{
		log.Fatalf("Main.db.Connect:%v", err)
	}
	fmt.Println("Database Connected")
		fmt.Println("olx-api is running..")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed:%v", err)
	}



}
