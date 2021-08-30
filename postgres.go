package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
)

func postgresConnect(ctx context.Context) (*pgx.Conn, error) {
	passwordBytes, err := os.ReadFile("secrets/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read postgres secret: %v\n", err)
		os.Exit(1)
	}

	password := string(passwordBytes)
	password = strings.TrimSuffix(password, "\n")
	password = url.QueryEscape(password)

	url := fmt.Sprintf("postgresql://postgres:%s@db.tefpmcttotopcptdivsj.supabase.co:5432/postgres", password)
	// url := fmt.Sprintf("postgresql://postgres:%s@127.0.0.1:5432/postgres", "vish_")

	cfg, err := pgx.ParseConfig(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse postgres config: %v\n", err)
		os.Exit(1)
	}

	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = conn.Ping(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed to database: %v\n", err)
		os.Exit(1)
	}

	return conn, err
}
