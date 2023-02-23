package models

type Book struct {
	ID     string `json:"id" bson:"_id"`
	Title  string `json:"title" bson:"title"`
	Author string `json:"author" bson:"author"`
	Year   string `json:"year" bson:"year"`
	Url    string `json:"url" bson:"url"`
	Filename    string `json:"filename" bson:"filename"`	
}
