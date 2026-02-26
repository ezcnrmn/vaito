package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("no DB_DSN was provided")
	}

	selfUrl := fmt.Sprintf(":%s", os.Getenv("STORAGE_PORT"))
	_ = fmt.Sprintf("%s:%s", os.Getenv("GATEWAY_HOST"), os.Getenv("GATEWAY_PORT"))

	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	server := &http.Server{
		Addr: selfUrl,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from storage"))
		}),
	}

	log.Println("storage started on ", selfUrl)
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
