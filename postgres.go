package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jackc/pgx/v4"
)

func main() {
	password := url.QueryEscape("foo")
	url := fmt.Sprintf("postgresql://postgres:%s@db.tefpmcttotopcptdivsj.supabase.co:5432/postgres", password)

	cfg, err := pgx.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)
}
