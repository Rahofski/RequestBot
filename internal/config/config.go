package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
    Backend_URL string
    Gigahat_URL string
    Auth_URL string
)


func Init () error {
    err := godotenv.Load("../../.env")

    if err != nil {
		log.Println("No .env file found, using environment variables")
        return err
	}

    Backend_URL = os.Getenv("BACKEND_URL")

    if Backend_URL == "" {
		log.Fatal("BACKEND_URL is not set in .env file or environment variables")
        return fmt.Errorf("BACKEND_URL is not set in .env file or environment variables")
	}

    Gigahat_URL = os.Getenv("GIGACHAT_URL")
    if Gigahat_URL == "" {
        log.Fatal("GIGACHAT_URL is not set in .env file or environment variables")
        return fmt.Errorf("GIGACHAT_URL is not set in .env file or environment variables")
    }

    Auth_URL = os.Getenv("AUTH_URL")
    if Auth_URL == "" {
        log.Fatal("AUTH_URL is not set in .env file or environment variables")
        return fmt.Errorf("AUTH_URL is not set in .env file or environment variables")
    }

    return nil

}