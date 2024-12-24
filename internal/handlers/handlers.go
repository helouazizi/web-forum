// internal/handlers/handlers.go
package handlers

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/utils"
	"html/template"
	"net/http"
)

func init() {
	test := &models.HomePage{}
	// lets get the components path to parse them
	componentsPath, err := utils.GetFolderPath("..", "components")
	if err != nil {
		fmt.Printf(" Error getting components path: %v\n", err)
	}
	// lets parse all components as globale using templ.parseglobal
	// we will use this to render the components in the templates
	ComponentsTemplate, err := template.ParseGlob(componentsPath + "*.html")
	fmt.Println(test,ComponentsTemplate)
}

// lets craete a function to handle the root request
// for our application
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
