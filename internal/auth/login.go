package auth

import (
	"database/sql"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Log_in(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method Not Allowed")
		return
	}
	if IsCookieSet(r, "token") {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Page Not Found")
		return
	}
	UserName := r.FormValue("userName")
	Password := r.FormValue("userPassword")

	if UserName == "" || Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Bad Request")
		return
	}

	// lets check the valid data from user
	var pasword string
	var username string
	err := database.Database.QueryRow("SELECT userName , userPassword  FROM users WHERE  username= $1 ", UserName).Scan(&username, &pasword)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			pages.All_Templates.ExecuteTemplate(w, "sign_in.html", "Invalid Password or UserName") // should execute login page here for no rows err
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "internal server error")
		return

	}
	if err := bcrypt.CompareHashAndPassword([]byte(pasword), []byte(Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		pages.All_Templates.ExecuteTemplate(w, "sign_in.html", "Invalid Password or UserName")
		return
	}
	// lets update the soken
	Token := uuid.New().String()
	query := "UPDATE users SET token = ? WHERE userName = ?"
	res, err := database.Database.Exec(query, Token, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	rowsAffected, _ := res.RowsAffected()
	// it this statement the user not exist
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		pages.All_Templates.ExecuteTemplate(w, "sign_in.html", "Invalid Password or UserName")
		return
	}
	// after all that set the token
	cookie := &http.Cookie{
		Name:   "token",
		Value:  Token,
		MaxAge: 3600,
		Path:   "/",
	}
	http.SetCookie(w, cookie)
	//r.AddCookie(cookie)
	//log.Println(UserName, "logged in")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
