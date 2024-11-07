package handlers

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
	"os"
)

var dbPool *pgxpool.Pool

func Start() error {
	connStr := os.Getenv("DATABASE_URL") // Ensure DATABASE_URL is set in the environment
	var err error
	dbPool, err = pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	return nil
}

func Stop() {
	dbPool.Close()
}
