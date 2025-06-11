package initialializers

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	fmt.Println("Loading env variables")
	err := godotenv.Load(".env")
		fmt.Println("env loaded")


	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
