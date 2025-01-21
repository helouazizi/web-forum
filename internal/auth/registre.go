// internal/auth/auth.go
package auth

import (
	"fmt"
	"net/http"

	"forum/internal/database"
	"forum/internal/handlers"

	"golang.org/x/crypto/bcrypt"
)

<<<<<<< HEAD
var (
	pages  handlers.Pages
	server handlers.User
)

func Signup_treatment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
=======
func Register(w http.ResponseWriter, r *http.Request) {
	pages := handlers.Pagess.All_Templates
	if r.Method != http.MethodPost {
>>>>>>> f05865bf8652d85fe6467a0db1b304ff7db4c228
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	if r.URL.Path != "/create_account" || IsCookieSet(r, "session") {
		w.WriteHeader(http.StatusNotFound)
		pages.ExecuteTemplate(w, "error.html", "Page Not Found")
		return

	}

	User := r.FormValue("userName")
	Pass := r.FormValue("userPassword")
	Email := r.FormValue("userEmail")
	fmt.Println(User, Pass, Email)

	if User == "" || Pass == "" || Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		pages.ExecuteTemplate(w, "error.html", "Bad Request")
		return
	}
	Hach_pass, err := bcrypt.GenerateFromPassword([]byte(Pass), 10)
	fmt.Println(Hach_pass)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "Internal Server Error")
		return
	}
<<<<<<< HEAD
	user_already_exist := database.Database.QueryRow("SELECT userName FROM users WHERE userName = $1 ", User).Scan(&User)
	email_already_exist := database.Database.QueryRow("SELECT userEmail FROM users WHERE  userEmail = $1", Email).Scan(&Email)

	if user_already_exist != nil && email_already_exist != nil {
		if user_already_exist == sql.ErrNoRows && email_already_exist == sql.ErrNoRows {
			server.IsAuthenticated = false

			token := uuid.New().String()
			database.Database.Exec("INSERT INTO users  (userName, userPassword, userEmail) VALUES ($1, $2, $3)", User, Hach_pass, Email)
			http.Redirect(w, r, "/", http.StatusFound)
			http.SetCookie(w,
				&http.Cookie{
					Name:   "token",
					Value:  token,
					MaxAge: 3600,
				})
=======
	var userExist bool
	var emailExist bool
	emailErr := database.Database.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE userEmail = $1)", Email).Scan(&emailExist)
	userErr := database.Database.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE userName = $1)", User).Scan(&userExist)
>>>>>>> f05865bf8652d85fe6467a0db1b304ff7db4c228

	if userErr != nil || emailErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.ExecuteTemplate(w, "error.html", "Internal Server Error")
		return
	}
	fmt.Println(emailExist, userExist)
	if userExist || emailExist {
		fmt.Println(emailExist, userExist)
		fmt.Println("exist")
		pages.ExecuteTemplate(w, "error.html", "User already exists")
		return
	} else {
		fmt.Println("not exist")
		_, err := database.Database.Exec("INSERT INTO users (userName,userEmail,userPassword) VALUES ($1, $2, $3)", User, Email, string(Hach_pass))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			pages.ExecuteTemplate(w, "error.html", "Internal Server Error")
			return
		}
	}
	http.Redirect(w, r, "/", 302)
}
