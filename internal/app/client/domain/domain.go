package domain

const BaseURL = "http://localhost:8080/books"

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"data"`
}
