package dotenv

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Print(fmt.Sprintf("dotenv load error: %v", err))
	}
}

func LoadFromFile(fileName string) {
	if err := godotenv.Load(fileName); err != nil {
		log.Print(fmt.Sprintf("dotenv load error: %v", err))
	}
}
