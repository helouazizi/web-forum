// internal/database/database.go
// internal/database/db.go
package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"forum/internal/models"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func Create_database() {
	var err error
	Database, err = sql.Open("sqlite3", "./internal/database/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// lets open the schema file to execute the sql commands inside it
	schema, err := os.Open("./internal/database/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer schema.Close()

	// now lets read the schema file using the bufio package
	scanner := bufio.NewScanner(schema)
	var sql_command string
	lineIndex := 0
	for scanner.Scan() {

		lineIndex++
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "--") || strings.HasPrefix(line, "/*") || line == "" {
			continue
		}
		sql_command += line + " "
		// lets execute the sql command
		if strings.HasSuffix(sql_command, "; ") {
			_, err = Database.Exec(sql_command)
			if err != nil {
				log.Fatal(err, " line :", lineIndex)
			}
			// free up the sql command
			sql_command = ""
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("data base creatd succesfully")
}

func Fetch_Database() *[]models.Post {
	// lets connect to our dtatbase
	rows, err := Database.Query("SELECT title , content FROM posts")
	if err != nil {
		fmt.Println("Error executing query:", err)
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()
	// lets iterate over rows and store them our models
	data := &models.Data{}

	for rows.Next() {
		post := &models.Post{}
		rows.Scan(&post.PostTitle, &post.PostContent)
		//log.Println(title, content, "data extracted")
		data.Posts = append(data.Posts, *post)

	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return &data.Posts
}
