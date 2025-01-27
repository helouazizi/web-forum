// cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/handlers"
	"forum/pkg/config"
	"forum/pkg/logger"
)

func init() {
	database.Create_database()
	handlers.ParseTemplates()
}

func main() {
	// Get the current working directory
	logger, err := logger.Create_Logger()
	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}
	defer logger.Close()

	// lets load the configuration
	configuration := config.LoadConfig()
	fmt.Printf("Server starting on port: %d >>> http://localhost:8080\n", configuration.Port)

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/sign_up", handlers.Sign_Up)
	http.HandleFunc("/register", auth.Register)
	http.HandleFunc("/sign_in", handlers.Sign_In)
	http.HandleFunc("/log_in", auth.Log_in)
	http.HandleFunc("/log_out", auth.Log_out)
	http.HandleFunc("/create_post", handlers.Craete_Post)
	http.HandleFunc("/submit_post", handlers.Submit_Post)
	http.HandleFunc("/static/", handlers.Serve_Static)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configuration.Port), nil))
}
