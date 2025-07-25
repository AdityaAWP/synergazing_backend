package main

import (
	"log"

	"github.com/joho/godotenv"
	"synergazing.com/synergazing/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load ENV")
	}

	config.ConnectEnvDBConfig()
}
