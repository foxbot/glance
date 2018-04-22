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
	host := os.Getenv("HOST_ADDR")
	if host == "" {
		log.Fatalln("Missing HOST_ADDR")
	}
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		log.Fatalln("Missing SECRET_KEY")
	}

	server := server.Server{
		Host: host,
		Key:  key,
	}
	server.Run()
}
