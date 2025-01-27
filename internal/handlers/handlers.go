package handlers

import (
	"log"
	"net/http"
	"text/template"

	"forum/internal/database"
)

type Pages struct {
	All_Templates *template.Template
}

var Pagess Pages

func ParseTemplates() {
	var err error
	Pagess.All_Templates, err = template.ParseGlob("./web/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	Pagess.All_Templates, err = Pagess.All_Templates.ParseGlob("./web/components/*.html")
	if err != nil {
		log.Fatal(err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Page not found")
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed hassan")
		return
	}
	data := database.Fetch_Database(r)
	Pagess.All_Templates.ExecuteTemplate(w, "home.html", data)
}

func Sign_Up(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed hassan")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "sign_up.html", nil)
	return
}

func Sign_In(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Pagess.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed hassan")
		return
	}
	Pagess.All_Templates.ExecuteTemplate(w, "sign_in.html", nil)
	return
}

func Serve_Static(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./web/static"))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	return
}
