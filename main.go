package main

import (
	"github.com/pantuza/go-sample-app/api"
	"log"
)

func main() {

	log.Println("Initializing go-sample-application")
	api.RunServer()
}
