package initialializers

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	fmt.Println("Loading env variables")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
