package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {

	host := getEnv("DB_HOST", "postgres")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "admin")
	password := getEnv("DB_PASSWORD", "admin")
	dbname := getEnv("DB_NAME", "shortener")
	sslmode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Optional: configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Retry logic
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	for {
		if err := db.PingContext(ctx); err == nil {
			fmt.Println("✅ connected to postgres")
			return db, nil
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("database connection timeout: %w", ctx.Err())
		default:
			fmt.Println("⏳ waiting for database...")
			time.Sleep(2 * time.Second)
		}
	}
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
