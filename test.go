package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// func main() {
// 	database, _ := sql.Open("sqlite3", "database.db")
// 	defer database.Close()
// 	createTable(database)

// 	message := ""
// 	fmt.Println("start print")
// 	for {
// 		fmt.Scanln(&message)
// 		if message == "show" {
// 			fetchRecords(database, 1)
// 			continue
// 		}
// 		if message == "clear" {
// 			dropRecords(database)
// 			continue
// 		}
// 		addUsers(database, message)
// 	}
// }

// func createTable(db *sql.DB) {
// 	users_table := `CREATE TABLE IF NOT EXISTS users (
//         id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
//         "text" TEXT UNIQUE)`
// 	query, err := db.Prepare(users_table)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	query.Exec()
// 	fmt.Println("Table created successfully!")
// }

// type Post struct {
// 	id      int
// 	content string
// }

// type Comment struct {
// 	id      int
// 	content string
// 	post    *Post
// }

// func addUsers(db *sql.DB, message string) {
// 	records := `INSERT INTO users(text) VALUES (?)`
// 	query, err := db.Prepare(records)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	_, err = query.Exec(message)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func fetchRecords(db *sql.DB, index int) {
// 	record, err := db.Query("SELECT * FROM users")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer record.Close()
// 	for record.Next() {
// 		var id int
// 		var text string
// 		record.Scan(&id, &text)
// 		fmt.Println(id, text)
// 	}
// }

// func dropRecords(db *sql.DB) {
// 	query, err := db.Prepare("DROP TABLE IF EXISTS users;")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	query.Exec()
// 	fmt.Println("database cleared")
// 	createTable(db)
// }

func main() {
	test := []byte("test")
	clean := []byte("clean")
	crypted, _ := bcrypt.GenerateFromPassword(test, 3)
	cryptclean, _ := bcrypt.GenerateFromPassword(clean, 3)
	fmt.Println(string(crypted))
	fmt.Println(bcrypt.CompareHashAndPassword(crypted, clean))
	fmt.Println(bcrypt.CompareHashAndPassword(cryptclean, test))
	fmt.Println(bcrypt.CompareHashAndPassword(cryptclean, clean))
}
