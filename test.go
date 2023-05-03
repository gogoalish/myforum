package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "forum.db")
	count := 0
	query := `INSERT INTO comments values(null, 1, 1, "second", 0)`
	db.Exec(query)
	fmt.Println(count)
}
