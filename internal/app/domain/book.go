package domain

type Book struct {
	ID     string `json:"id" example:"8502ab55-6750-4c53-8126-acc1ba19f801"`
	Title  string `json:"title" example:"The book"`
	Author string `json:"author" example:"The author"`
	Year   int    `json:"data" example:"2023"`
	Status string `json:"status,omitempty" example:"deleted"`
}
