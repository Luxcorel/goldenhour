package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Luxcorel/golden-hour/internal/handlers"
)

func main() {
	port := flag.String("port", "8080", "Usage: -port [PORT]")
	flag.Parse()

	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + *port,
		Handler: router,
	}

	log.Printf("Starting server on :%v", *port)

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("An error occured:", err)
	}
}
