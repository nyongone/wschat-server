package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	Host string
	Port string
}

var Env = new(Environment)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	Env.Host = os.Getenv("APP_HOST")
	Env.Port = os.Getenv("APP_PORT")
}
