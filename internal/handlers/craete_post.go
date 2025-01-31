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
	time := time.Now().Format(time.DateTime)
	fmt.Println(categories)

	// lets check for emptyness
	if title == "" || content == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}
	var userID int
	query1 := "SELECT id FROM users WHERE token = ? "
	err = database.Database.QueryRow(query1, cookies.Value).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found for the given token")
		} else {
			log.Println("Database query error:", err)
		}
		return
	}
	// lets insert this data to our database
	query := "INSERT INTO posts (user_id,title,content,created_at) VALUES (?,?,?,?)"
	_, err = database.Database.Exec(query, userID, title, content, time)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	//r.AddCookie(cookies)
	http.Redirect(w, r, "/", http.StatusFound)
}
