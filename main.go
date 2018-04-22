package main

import (
	"log"
	"os"

	"github.com/dabbotorg/glance/server"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	server := server.Server{
		Host: os.Getenv("HOST_ADDR"),
	}
	server.Run()
}
