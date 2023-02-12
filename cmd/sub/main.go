package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v4"
)

func main() {
	var (
		channel string
		lock    int
		url     string
	)
	flag.IntVar(&lock, "lock", 1, "advisory lock to aquire")
	flag.StringVar(&channel, "channel", "payment", "Channel to subscribe to")
	flag.StringVar(&url, "database", "postgres://postgres:postgres@localhost:5432/postgres?ssl=mode=disable", "Connection URI")
	flag.Parse()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	log.Printf("waiting for lock (%d)...\n", lock)
	_, err = conn.Exec(ctx, "select pg_advisory_lock($1)", lock)
	if err != nil {
		log.Fatalf("Unable to aquire lock: %v\n", err)
	}

	_, err = conn.Exec(ctx, "listen "+channel)
	if err != nil {
		log.Fatalf("Unable to listen: %v\n", err)
	}

	log.Printf("Subscribed to %s", channel)
	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("notification from %s: %s\n", notification.Channel, notification.Payload)
	}
}
