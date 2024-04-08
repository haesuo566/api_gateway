package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Getenv(key string) string {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println(err)
		return ""
	}

	return os.Getenv(key)
}
