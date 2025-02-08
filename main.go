package main

import (
	"database/sql"
	"flag"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/muesli/termenv"
)

var db *sql.DB

func main() {
	var rm = flag.Int("rm", 0, "ID of the note to remove")
	var add = flag.String("add", "", "Note to add")
	var ls = flag.Bool("ls", false, "List notes")
	flag.Parse()

	var err error
	db, err = sql.Open("sqlite3", "./.welln.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS app (id INTEGER PRIMARY KEY, note TEXT, data INTEGER)")
	if err != nil {
		log.Fatal(err)
	}

	if *add != "" {
		addNote(*add)
	} else if *rm > 0 {
		remove(*rm)
	} else if *ls {
		listNotes()
	} else {
		fmt.Println("No valid flags provided. Use -add, -rm, or -ls.")
	}
}

func addNote(note string) {
	stmt, err := db.Prepare("INSERT INTO app (note, data) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(note, time.Now().UTC().Unix())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Note added successfully.")
}

func remove(id int) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM app WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		fmt.Println("Note with ID", id, "does not exist.")
		return
	}

	_, err = db.Exec("DELETE FROM app WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Note with ID %d removed successfully.\n", id)
}

func listNotes() {
	rows, err := db.Query("SELECT id, note, data FROM app ORDER BY data DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var note string
		var timestamp int64
		err := rows.Scan(&id, &note, &timestamp)
		if err != nil {
			log.Fatal(err)
		}

		t := time.Unix(timestamp, 0).UTC()

		ls := termenv.String(fmt.Sprintf("%d, %s ~ %s\n", id, t.Format("2006-01-02"), note))
		fmt.Print(ls)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
