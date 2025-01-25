package auth

import (
	"fmt"
	"forum/internal/database"
	"forum/internal/handlers"
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
	Update_Token(w, r, "token", "")
	http.SetCookie(w, &http.Cookie{
		Name:     "token",         // name of the cookie
		Value:    "",              // clear the cookie value
		Expires:  time.Unix(0, 0), // set expiration time to a time in the past
		Path:     "/",             // scope of the cookie
		HttpOnly: true,            // prevent JavaScript access
		Secure:   true,            // ensure cookie is only sent over HTTPS
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func Update_Token(w http.ResponseWriter, r *http.Request, tokenName string, token string) {
	pages := handlers.Pagess
	cokie, err := r.Cookie(tokenName)
	if err != nil {
		// If no cookie exists, handle appropriately
		w.WriteHeader(http.StatusUnauthorized)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Unauthorized: No token provided")
		return
	}

	fmt.Println("Token from cookie:", cokie.Value)

	// Update the token to NULL or an empty string in the database
	query := "UPDATE users SET token = $1 WHERE token = $2"
	_, err = database.Database.Exec(query, token, cokie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Internal server error")
		return
	}

	fmt.Println("Token removed from database")
}
