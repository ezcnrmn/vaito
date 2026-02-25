package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

// TODO:
const default_dsn string = "postgres://vaito:pa55word@postgres_db:5432/vaito?sslmode=disable"

func main() {
	dsn := flag.String("db-dsn", default_dsn, "PostgteSQL DSN")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	server := &http.Server{
		Addr: ":4001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from storage"))
		}),
	}

	fmt.Println("storage started")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// TODO: вынести в конфиг
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
