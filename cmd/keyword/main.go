package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/delivery/rest"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print(err)
		log.Fatal("could not load env file, shutting down")
	}

	s := rest.NewServer()
	http.ListenAndServe(":8080", s)
}
