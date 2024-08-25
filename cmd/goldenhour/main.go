package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Luxcorel/goldenhour/internal/handlers"
)

func main() {
	port := flag.String("port", "8080", "Usage: -port [PORT]")
	flag.Parse()

	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + *port,
		Handler: router,
	}

	calc := handlers.NewGetCalc()
	router.HandleFunc("GET /{$}", handlers.GetHome)
	router.HandleFunc("GET /calc", calc.GetCalc)

	log.Printf("Starting server on :%v", *port)

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("An error occured:", err)
	}
}
