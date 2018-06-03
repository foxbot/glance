package main

import (
	"log"
	"os"
	"strings"

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
	ips := os.Getenv("IP_LIST")
	if key == "" {
		log.Fatalln("Missing IP_LIST")
	}

	server := server.Server{
		Host: host,
		Key:  key,
		Ips:  strings.Split(ips, ","),
	}
	server.Run()
}
