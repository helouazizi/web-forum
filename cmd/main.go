// cmd/main.go
package main

import (
	"fmt"
	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/handlers"
	"log"
	"net/http"
)

func init() {
	database.Create_database()
	handlers.ParseTemplates()
}

func main() {
	// Get the current working directory
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/create_account", auth.Register)
	http.HandleFunc("/log_in", auth.Log_in)
	http.HandleFunc("/log_out", auth.Log_out)
	http.HandleFunc("/create_Post", handlers.Craete_post)
	http.HandleFunc("/static/", handlers.Serve_Static)
	fmt.Println("server is running on port 8080 ... http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
