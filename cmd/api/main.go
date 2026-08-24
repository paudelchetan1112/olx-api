package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/paudelchetan1112/olx-api/internal/config"
	"github.com/paudelchetan1112/olx-api/internal/db"
	"github.com/paudelchetan1112/olx-api/internal/router"
)

func main() {
	cfg := config.MustLoad() //package which config env variable 
	db, err:=db.Connect(cfg.DatabaseUrl)
	if err!=nil{
		log.Fatalf("Main.db.Connect:%v", err)
	}
	fmt.Println("Database Connected")
		fmt.Println("olx-api is running..")
	
	
	router:=router.Router(db)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed:%v", err)
	}



}
