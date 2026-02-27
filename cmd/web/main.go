package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"readinglist.demo.io/internal/models"
)

type application struct {
	readinglist *models.ReadinglistModel
}

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "HTTP network address")
	endpoint := flag.String("endpoint", os.Getenv("READINGLIST_API_ENDPOINT"), "Endpoint for the readinglist web service")

	app := &application{
		readinglist: &models.ReadinglistModel{Endpoint: *endpoint},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Printf("Starting the server on %s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
