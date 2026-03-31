package db

import (
	"bobcatsar-max-bot/internal/config"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func ConnectionDB(config *config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.Root,
		config.Password,
		config.Dbname,
	)
	var err error
	var db *sql.DB

	for i := 0; i < 20; i++ {
		db, err = sql.Open("pgx", dataSourceName)
		log.Printf("connection: %s", dataSourceName)

		if err != nil {
			log.Printf("❌ Failed to open DB (try %d/20): %v", i+1, err)
		} else if pingErr := db.Ping(); pingErr == nil {
			log.Println("✅ Connected to PostgreSQL")
			return db, nil
		} else {
			log.Printf("⚠️ Waiting for DB (try %d/20)...")
			db.Close()
		}

		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to DB: %w", err)
}
