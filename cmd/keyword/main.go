package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/delivery/rest"
	"github.com/joho/godotenv"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wd)
	err = godotenv.Load(path.Join(wd, ".env"))
	if err != nil {
		fmt.Print(err)
		log.Fatal("could not load env file, shutting down")
	}

	s := rest.NewServer()
	http.ListenAndServe(":8080", s)
}
