// internal/handlers/handlers.go
package handlers

import (
	"forum/internal/utils"
	"html/template"
	"net/http"
)

type Pages struct {
	All_Templates *template.Template
}
type Form struct {
	Title  string
	Button string
}

var pages Pages

func init() {
	var err error
	path, err := utils.GetFolderPath("..", "templates")
	if err != nil {
		panic(err)
	}
	pages.All_Templates, err = template.ParseGlob(path + "/*.html")
	if err != nil {
		panic(err)
	}

}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Page not found")
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}

	pages.All_Templates.ExecuteTemplate(w, "home.html", nil)

}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	data := Form{
		Title:  "Login",
		Button: "Login",
	}

	pages.All_Templates.ExecuteTemplate(w, "login.html", data)
}
func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	data := Form{
		Title:  "Create Account",
		Button: "Create Account",
	}

	pages.All_Templates.ExecuteTemplate(w, "login.html", data)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}

	pages.All_Templates.ExecuteTemplate(w, "createpost.html", nil)
}

func Serve_Static(w http.ResponseWriter, r *http.Request) {
	path, _ := utils.GetFolderPath("..", "static")
	fs := http.FileServer(http.Dir(path))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
}
