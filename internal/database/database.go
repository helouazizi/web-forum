// internal/database/database.go
// internal/database/db.go
package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
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

func Fetch_Database(r *http.Request) *models.Data {
	// lets connect to our dtatbase
	rows, err := Database.Query("SELECT title , content , created_at , total_likes , total_dislikes FROM posts")
	if err != nil {
		fmt.Println("Error executing query:", err)
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()
	// lets iterate over rows and store them our models
	data := &models.Data{}

	// lets check if the user have a token

	if t, err := r.Cookie("token"); err == nil {
		// Only access `t.Value` if no error occurred
		if t.Value != "" {
			data.Userr.IsLoged = true
		}
	}
	// lets extract hus username
	userName := r.FormValue("userName")
	Email := r.FormValue("userEmail")
	if Email == "" {
		Database.QueryRow("SELECT userEmail FROM users WHERE userName = $1", userName).Scan(&Email)
	}
	fmt.Println(userName, Email)
	data.Userr.UserName = userName
	data.Userr.UserEmail = Email

	for rows.Next() {
		post := &models.Post{}
		rows.Scan(&post.PostTitle, &post.PostContent, &post.PostCreatedAt, &post.TotalLikes, &post.TotalDeslikes)
		//log.Println(title, content, "data extracted")
		data.Posts = append(data.Posts, *post)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return data
}
