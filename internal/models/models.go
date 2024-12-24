// internal/models/models.go
// models/models.go
package models

import "html/template"

type Components struct {
	Header *template.Template
	Aside  *template.Template
	Postt   *template.Template
	Title  string
	Form   *template.Template
	Footer *template.Template
}
type HomePage struct {
	Header *template.Template
	Aside  *template.Template
	Title  string
	Posts  []*Post
	Footer *template.Template
}

type Post struct {
	Content *template.Template
}

type SignPage struct {
	Header *template.Template
	Form   *template.Template
	Footer *template.Template
}

type CreatePage struct {
	Header     *template.Template
	CraeteForm *template.Template
	Footer     *template.Template
}
