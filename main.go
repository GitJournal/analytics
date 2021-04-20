package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mailru/go-clickhouse"
)

type tuple struct {
	name  string
	value string
}

func main() {
	connect, err := sql.Open("clickhouse", "http://127.0.0.1:8123/default")
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS example3 (
			country_code FixedString(2),
			t        Tuple(String, String)
		) engine=Memory
	`)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Going to insert")

	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO example3 (country_code, t) VALUES (?, ?)")
	)
	defer stmt.Close()

	if _, err := stmt.Exec(
		"RU",
		tuple{"a", "b"},
		// clickhouse.Array([]tuple{{"a", "b"}, {"c", "d"}}),
		clickhouse.Array([]tuple{}),
	); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

}
