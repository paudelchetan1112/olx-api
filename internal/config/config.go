package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseUrl string
}

func MustLoad() Config {
	godotenv.Load("cmd/api/.env")
	port := os.Getenv("PORT")
	fmt.Println("PORT", port)
	if port == "" {
		panic("PORT is required")
	}
	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required")
	}
	DatabaseUrl := os.Getenv("DATABASE_URL")
	if DatabaseUrl == "" {
		panic("DATABASE URL is required")
	}
	return Config{
		Port:        port,
		Env:         env,
		DatabaseUrl: DatabaseUrl,
	}
}
