package domain

type Book struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Year      int    `json:"data"`
	IsDeleted bool   `json:"status,omitempty"`
}
