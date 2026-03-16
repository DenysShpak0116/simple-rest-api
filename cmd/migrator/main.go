package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"simple-rest-api/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	config := config.MustLoad()

	command := "up"
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	db, err := sql.Open("postgres", config.ConnectionString)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	goose.SetBaseFS(nil)
	dir := "migrations"
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	log.Printf("Running goose %s...", command)
	if err := goose.RunContext(context.Background(), command, db, dir); err != nil {
		log.Fatalf("goose %s: %v", command, err)
	}

	log.Println("Migration command completed successfully")
}
