package handlers

import (
	"fmt"
	"net/http"
	"time"

	"forum/internal/database"
)

func Submit_Post(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	r.ParseForm()
	user_Id := 1
	categories := r.Form["categorie"] ///////////////////////  be carefull
	title := r.FormValue("postTitle")
	content := r.FormValue("postContent")
	time := time.Now().Format(time.DateOnly)
	fmt.Println(categories)

	// lets check for emptyness
	if title == "" || content == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "bad request")
		return
	}
	// lets insert this data to our database
	query := "INSERT INTO posts (user_id,title,content,created_at) VALUES (?,?,?,?)"
	_, err := database.Database.Exec(query, user_Id, title, content, time)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func Craete_Post(w http.ResponseWriter, r *http.Request) {
	pages := Pagess.All_Templates
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	data := database.Fetch_Database(r)
	pages.ExecuteTemplate(w, "createPost.html", data)
	return
}
