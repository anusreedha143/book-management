package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq" // Go package for sql postgres driver
	"readinglist.demo.io/internal/data"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	dsn  string // connection string to the database
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API Server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("READINGLIST_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	if cfg.dsn == "" {
		cfg.dsn = os.Getenv("READINGLIST_DB_DSN")
		if cfg.dsn == "" {
			cfg.dsn = "postgres://readinglistdbuser:vikky@postgres_db:5432/readinglist?sslmode=disable"
		}
		logger.Printf("Connection string %v", cfg.dsn)
	}

	if cfg.dsn == "" {
		logger.Fatal("No database DSN provided")
	}

	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("database connection pool establised")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)

}
