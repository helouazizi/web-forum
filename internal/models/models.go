package models

type User struct {
	IsLoged   bool
	UserName  string
	UserEmail string
	Profile   string // about imges we can store them in databse as blob dont wory
}

type Data struct {
	Userr User
	Posts []Post
}

type Post struct {
	PostCreator                       string
	PostCreatedAt                     string
	PostTitle                         string
	PostContent                       string
	TotalLikes, TotalDeslikes, PostId int
	Categories                        []Categorie
}
type Categorie struct {
	CatergoryName string
}
