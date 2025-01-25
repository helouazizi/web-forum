package auth

import (
	"fmt"
	"forum/internal/database"
	"forum/internal/handlers"
	"log"
	"net/http"
	"time"
)

func Log_out(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	// lets check in first that is already have a session
	if !IsCookieSet(r, "token") {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "page not founttt")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",         // name of the cookie
		Value:    "",              // clear the cookie value
		Expires:  time.Unix(0, 0), // set expiration time to a time in the past
		Path:     "/",             // scope of the cookie
		HttpOnly: true,            // prevent JavaScript access
		Secure:   true,            // ensure cookie is only sent over HTTPS
	})
	log.Print("A User logged out")
	// lets remove his token from database
	token, err := r.Cookie("token")
	if err == nil {
		fmt.Println(token.Value)
		_, err := database.Database.Exec("UPDATE users SET token = $1 WHERE token = $1 ", "", token.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			pages.All_Templates.ExecuteTemplate(w, "error.html", "internal server error")
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
