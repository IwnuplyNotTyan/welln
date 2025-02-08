package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS app (id INTEGER PRIMARY KEY, note TEXT, data INTEGER)")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO app (note, data) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(os.Args[1], time.Now().UTC().Unix()) // Store time as Unix timestamp (int64)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT data, note FROM app")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp int64
		var note string
		if err := rows.Scan(&timestamp, &note); err != nil {
			log.Fatal(err)
		}
		t := time.Unix(timestamp, 0).UTC()
		fmt.Printf("%s ~ %s\n", t.Format("2006-01-02"), note)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
