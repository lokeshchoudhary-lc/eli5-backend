package main

import (
	"eli5/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := server.Create()

	server.Listen(app)

}
