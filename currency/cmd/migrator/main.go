package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	db "github.com/lekss361/curserv2/currency/internal/db"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mode := flag.String("mode", "up", "Migration mode: up or down")
	flag.Parse()

	database, err := db.New()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer database.Close()

	done := make(chan error, 1)
	go func() {
		switch *mode {
		case "down":
			done <- db.MigrateDown(database)
		default:
			done <- db.Migrate(database)
		}
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Fatalf("Migration %s failed: %v", *mode, err)
		}
		log.Printf("Migration %s completed successfully", *mode)

	case <-ctx.Done():
		log.Println("Received shutdown signal, aborting migration")
		os.Exit(1)
	}
}
