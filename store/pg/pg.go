package pg

import (
	"context"
	"fmt"
	"hwAgentApp/config"
	"os"

	"github.com/jackc/pgx/v4"
)

type DB struct {
	*pgx.Conn
}

func Dial() (*DB, error) {
	cfg := config.Get()
	if cfg.PgAddr == "" || cfg.PgPort == "" ||
		cfg.PgDB == "" || cfg.PgUser == "" || cfg.PgPassword == "" {
		return nil, nil
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"

	pgOpts := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.PgUser, cfg.PgPassword, cfg.PgAddr, cfg.PgPort, cfg.PgDB,
	)
	pgDB, err := pgx.Connect(context.TODO(), pgOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	var n int
	if err = pgDB.QueryRow(context.TODO(), "SELECT 1").Scan(&n); err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return &DB{pgDB}, nil
}
