// internal/database/database.go
// internal/database/db.go
package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func init() {
	var err error
	Database, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
<<<<<<< HEAD
	createTable := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"username" TEXT,
		"password" TEXT
	);`
	_, err = Database.Exec(createTable)
=======

	// lets open the schema file to execute the sql commands inside it
	schema, err := os.Open("./internal/database/schema.sql")
>>>>>>> f05865bf8652d85fe6467a0db1b304ff7db4c228
	if err != nil {
		log.Fatal(err)
	}
}
