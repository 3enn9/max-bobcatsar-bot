package db

import (
	"bobcatsar-max-bot/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

func ConnectionDB(config *config.Config) (*pgx.Conn, error) {
	dataSourceName := fmt.Sprintf(
		"postgres://%s:%s@postgres:5432/%s",
		config.Root,
		config.Password,
		config.Dbname,
	)
	var err error
	var conn *pgx.Conn

	for i := 0; i < 20; i++ {
		conn, err = pgx.Connect(context.Background(), dataSourceName)
		log.Printf("connection: %s", dataSourceName)
		if err != nil {
			log.Printf("❌ Failed to open DB (try %d/20): %v", i+1, err)
		} else if pingErr := conn.Ping(context.Background()); pingErr == nil {
			log.Println("✅ Connected to PostgreSQL")
			return conn, nil
		} else {
			log.Printf("⚠️ Waiting for DB (try %d/20)...", i+1)
			conn.Close(context.Background())
		}

		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to DB: %w", err)
}
