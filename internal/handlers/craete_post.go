package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/internal/database"
)

func Craete_Post(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	data := database.Fetch_Database(r)
	pages.ExecuteTemplate(w, "createPost.html", data)
	//return
}
func Submit_Post(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	r.ParseForm()
	cookies, err := r.Cookie("token")
	if err != nil {
		fmt.Println(err)
		return
	}
	categories := r.Form["categorie"] ///////////////////////  be carefull
	title := r.FormValue("postTitle")
	content := r.FormValue("postContent")
	time := time.Now().Format(time.RFC3339)
	fmt.Println(categories)

	// lets check for emptyness
	if title == "" || content == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}

	var userID int
	query1 := "SELECT id FROM users WHERE token = ?"
	err = database.Database.QueryRow(query1, cookies.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found for the given token")
		} else {
			log.Println("Database query error:", err)
		}
		return
	}

	// Insert post into database
	query := "INSERT INTO posts (user_id, title, content, created_at) VALUES (?, ?, ?, ?)"
	result, err := database.Database.Exec(query, userID, title, content, time)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	postID, err := result.LastInsertId()
	if err != nil {
		log.Println("Failed to retrieve last inserted post ID:", err)
		return
	}

	// Insert categories into post_categories table
	for _, category := range categories {
		var categoryID int
		categoryQuery := "SELECT id FROM categories WHERE category_name = ?"
		err := database.Database.QueryRow(categoryQuery, category).Scan(&categoryID)
		if err == sql.ErrNoRows {
			insertCategoryQuery := "INSERT INTO categories (category_name) VALUES (?)"
			res, err := database.Database.Exec(insertCategoryQuery, category)
			if err != nil {
				log.Println("Error inserting new category:", err)
				continue
			}
			categoryID64, _ := res.LastInsertId()
			categoryID = int(categoryID64)
		} else if err != nil {
			log.Println("Error retrieving category ID:", err)
			continue
		}

		insertPostCategoryQuery := "INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)"
		_, err = database.Database.Exec(insertPostCategoryQuery, postID, categoryID)
		if err != nil {
			log.Println("Error inserting into post_categories:", err)
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
