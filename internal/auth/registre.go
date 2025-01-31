// internal/auth/auth.go
package auth

import (
	"log"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess.All_Templates
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	if IsCookieSet(r, "token") {
		w.WriteHeader(http.StatusNotFound)
		pages.ExecuteTemplate(w, "error.html", "not found")
		return
	}
	r.ParseForm()
	userName := r.FormValue("userName")
	userPassword := r.FormValue("userPassword")
	Email := r.FormValue("userEmail")
	Token := uuid.New().String()
	// lets check for emptyness
	if userName == "" || userPassword == "" || Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "Bad Request")
		return
	}
	Hach_pass, err := bcrypt.GenerateFromPassword([]byte(userPassword), 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "internal server error")
		return
	}
	var userExist bool
	var emailExist bool
	emailErr := database.Database.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE userEmail = $1)", Email).Scan(&emailExist)
	userErr := database.Database.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE userName = $1)", userName).Scan(&userExist)

	if userErr != nil || emailErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "Internal Server Error")
		return
	}

	if userExist || emailExist {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "sign_up.html", "Invalid Password or UserName")
		return
	} else {
		_, err := database.Database.Exec("INSERT INTO users (userName,userEmail,userPassword,token) VALUES ($1, $2, $3 , $4 )", userName, Email, string(Hach_pass), Token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate(w, "error.html", "Internal Server Error")
			return
		}
		log.Printf("%s account has been created", userName)
	}

	cookie := &http.Cookie{
		Name:   "token",
		Value:  Token,
		MaxAge: 3600,
		Path:   "/",
	}

	http.SetCookie(w, cookie)
	//r.AddCookie(cookie)
	http.Redirect(w, r, "/", http.StatusFound)
	//return
}
