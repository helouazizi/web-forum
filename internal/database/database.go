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
	"forum/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func NewDatabase() (*Database, error) {
	dbPath := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	return &Database{db}, nil
}

func Create_database() {
	db, err := NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// lets open the schema file to execute the sql commands inside it
	shema_path := os.Getenv("SCHEMA_PATH")
	schema, err := os.Open(shema_path)
	if err != nil {
		logger.LogWithDetails(err)
		log.Fatal(err)
	}
	defer schema.Close()
	// now lets read the schema file using the bufio package
	scanner := bufio.NewScanner(schema)
	var sql_command string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "--") || strings.HasPrefix(line, "/*") || line == "" {
			continue
		}
		sql_command += line + " "
		// lets execute the sql command
		if strings.HasSuffix(sql_command, "; ") {
			_, err = db.Exec(sql_command)
			if err != nil {
				logger.LogWithDetails(err)
				log.Fatal(err)
			}
			// free up the sql command
			sql_command = ""
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func Fetch_Database(r *http.Request, query string, userid int, liked bool) (*models.Data, error) {
	var finalQuery string
	if userid > 0 && !liked { // posts of a single user
		finalQuery = fmt.Sprintf("%s WHERE users.id = %d ORDER  BY posts.created_at DESC;", query, userid)
	} else if userid > 0 && liked { // liked posts
		finalQuery = fmt.Sprintf("%s WHERE  post_reaction.user_id = %d AND  post_reaction.reaction = 1", query, userid)
	} else { // all posts
		finalQuery = fmt.Sprintf("%s ORDER BY posts.created_at DESC", query)
	}

	db, err := NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
	}
	data := &models.Data{}
	stm, err := db.Prepare(finalQuery)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	rows, err := stm.Query()
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	defer rows.Close()

	// lets check if the user have a token
	if t, err := r.Cookie("token"); err == nil {
		if t.Value != "" {
			data.User.IsLoged = true
		}
	}
	// lets extract his username
	userName := r.FormValue("userName")
	Email := r.FormValue("userEmail")
	if Email == "" {
		stm, err := db.Prepare("SELECT userEmail FROM users WHERE userName = ? ")
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		stm.QueryRow(userName).Scan(&Email)
	}
	data.User.UserName = userName
	data.User.UserEmail = Email

	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.PostId, &post.PostTitle, &post.PostContent, &post.TotalLikes, &post.TotalDeslikes, &post.PostCreatedAt, &post.PostCreator, &post.UserID,
		)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		// Fetch categories for the post
		query := "SELECT category FROM categories WHERE post_id = ?"
		stm, err := db.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		rows2, err := stm.Query(post.PostId)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		defer rows2.Close()
		for rows2.Next() {
			categ := &models.Categorie{}
			err := rows2.Scan(&categ.CatergoryName)
			if err != nil {
				logger.LogWithDetails(err)
				return nil, err
			}
			post.Categories = append(post.Categories, *categ)
		}
		data.Posts = append(data.Posts, *post)
	}

	if err := rows.Err(); err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	/// lets fetch cetegories
	query2 := `SELECT category FROM stoke_categories`
	stm, err = db.Prepare(query2)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}

	rows, err = stm.Query()
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		category := &models.Categorie{}
		err := rows.Scan(&category.CatergoryName)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		data.Categories = append(data.Categories, *category)
	}
	if err := rows.Err(); err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}

	return data, nil
}
