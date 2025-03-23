package init

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvironmentalVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Couldn't load .env file")
	}
}
