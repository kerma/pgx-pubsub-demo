package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/oklog/ulid/v2"

	"github.com/jackc/pgx/v4"
)

func main() {
	var (
		kind string
		dir  string
		url  string
	)
	flag.StringVar(&kind, "kind", "payment", "Event kind")
	flag.StringVar(&dir, "dir", "db/migrations", "migrations directory")
	flag.StringVar(&url, "database", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "Connection URI")
	flag.Parse()

	{ // apply migrations
		m, err := migrate.New("file://"+dir, url)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Unable to apply down migrations: %s", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Unable to apply migrations: %s", err)
		}
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	for {
		data := map[string]string{
			"id":   ulid.Make().String(),
			"kind": kind,
			"body": "example message body",
		}
		body, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("marshal failed: %s", err)
		}

		var id int
		err = conn.QueryRow(ctx, "insert into events(data) values($1) returning id", body).Scan(&id)
		if err != nil {
			log.Fatalf("query failed: %s", err)
		}
		log.Printf("event inserted: %d - %s\n", id, kind)
		time.Sleep(time.Second)
	}
}
