package models

type User struct {
	IsLoged  bool
	UserName string
	Profile  string // about imges we can store them in databse as blob dont wory
}

type Data struct {
	Posts []Post
}

type Post struct {
	PostCreator               User
	PostCreatedAt             string
	PostTitle                 string
	PostContent               string
	TotalLikes, TotalDeslikes int
	Categories                []Categorie
}
type Categorie struct {
	CatergoryName string
}
