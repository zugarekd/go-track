package handlers

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

var dbPool *pgxpool.Pool

type Config struct {
	connStr string
}

func Init() error {
	config = Config{
		connStr: os.Getenv("DATABASE_URL"),
	}
	return nil
}

var config Config

func Start() error {
	var err error
	dbPool, err = pgxpool.Connect(context.Background(), config.connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	return nil
}

func Stop() {
	dbPool.Close()
}
