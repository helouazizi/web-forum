// internal/database/database.go
// internal/database/db.go
package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"forum/internal/models"

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
	// Connect to the database
	query := `
		SELECT 
			posts.id, posts.title, posts.content, posts.total_likes, posts.total_dislikes, posts.total_comments, posts.created_at,
			users.userName
		FROM 
			posts
		INNER JOIN 
			users
		ON 
			posts.user_id = users.id
		ORDER BY 
			posts.created_at DESC
	`
	rows, err := Database.Query(query)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil
	}
	defer rows.Close()

	// Initialize data model
	data := &models.Data{}

	// Check if the user has a token
	if t, err := r.Cookie("token"); err == nil && t.Value != "" {
		data.User.IsLoged = true
	}

	// Extract user information
	userName := r.FormValue("userName")
	Email := r.FormValue("userEmail")
	if Email == "" {
		Database.QueryRow("SELECT userEmail FROM users WHERE userName = ?", userName).Scan(&Email)
	}

	data.User.UserName = userName
	data.User.UserEmail = Email

	// Fetch posts
	for rows.Next() {
		post := models.Post{}
		err := rows.Scan(
			&post.PostId, &post.PostTitle, &post.PostContent, &post.TotalLikes, &post.TotalDeslikes, &post.TotalComments, &post.PostCreatedAt, &post.PostCreator,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		// Fetch categories for each post
		categoryQuery := "SELECT category_name FROM categories INNER JOIN post_categories ON categories.id = post_categories.category_id WHERE post_categories.post_id = ?"
		categoryRows, err := Database.Query(categoryQuery, post.PostId)
		if err == nil {
			for categoryRows.Next() {
				var category models.Categorie
				categoryRows.Scan(&category.CatergoryName)
				post.Categories = append(post.Categories, category)
			}
			categoryRows.Close()
		}

		data.Posts = append(data.Posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
	}

	return data
}
