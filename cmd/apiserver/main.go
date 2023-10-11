package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/Te8va/APIbook.git/iternal/handler"
	"github.com/Te8va/APIbook.git/iternal/domain"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func main() {

	err := handler.ReadBooksFromFile()
	if err != nil {
		fmt.Println("Ошибка при чтении книг из файла:", err)
		return
	}

	bookHandler := handler.NewBookHandler(domain.Book{})

	router := httprouter.New()

	router.GET("/books", bookHandler.GetBooksHandler)
	router.GET("/books/:id", bookHandler.GetBookByIDHandler)
	router.POST("/books", bookHandler.AddBookHandler)
	router.DELETE("/books/:id", bookHandler.DeleteBookHandler)
	router.PUT("/books/:id", bookHandler.UpdateBookHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		router.ServeHTTP(w, r)
	})

	fmt.Println("Сервер запущен на порту 8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}

}

// type Book struct {
// 	ID     int    `json:"id"`
// 	Title  string `json:"title"`
// 	Author string `json:"author"`
// 	Data   int    `json:"data"`
// }

// var books []Book

// const filePath = "books.json"

// func readBooksFromFile() error {
// 	file, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return fmt.Errorf("Ошибка чтения файла: %v", err)
// 	}

// 	err = json.Unmarshal(file, &books)
// 	if err != nil {
// 		return fmt.Errorf("Ошибка декодирования JSON: %v", err)
// 	}
// 	return nil
// }

// func saveBooksToFile() error{
// 	data, err := json.Marshal(books)
// 	if err != nil {
// 		return fmt.Errorf("Ошибка кодирования JSON: %v", err)
// 	}

// 	err = os.WriteFile(filePath, data, fs.ModePerm)
// 	if err != nil {
// 		return fmt.Errorf("Ошибка записи в файл:%v", err)
// 	}
// 	return nil
// }

// func getBooksHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(books)
// }

// func getBookByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	id, err := validateID(ps.ByName("id"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	for _, book := range books {
// 		if book.ID == id {
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(book)
// 			return
// 		}
// 	}

// 	http.Error(w, "Book not found", http.StatusNotFound)
// }

// func addBookHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	var newBook Book
// 	err := json.NewDecoder(r.Body).Decode(&newBook)
// 	if err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	newBook.ID = len(books) + 1 //TODO посложнее id

// 	books = append(books, newBook)

// 	err = saveBooksToFile()
// 	if err != nil {
// 		http.Error(w, "Ошибка при сохранении в файл", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(newBook)
// }

// func deleteBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	id, err := validateID(ps.ByName("id"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	index := -1
// 	for i, book := range books {
// 		if book.ID == id {
// 			index = i
// 			break
// 		}
// 	}

// 	if index == -1 {
// 		http.Error(w, "Book not found", http.StatusNotFound)
// 		return
// 	}

// 	books = append(books[:index], books[index+1:]...)

// 	err = saveBooksToFile()
// 	if err != nil {
// 		http.Error(w, "Ошибка при сохранении в файл", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// func updateBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

// 	id, err := validateID(ps.ByName("id"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var updatedBook Book
// 	err = json.NewDecoder(r.Body).Decode(&updatedBook)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	index := -1
// 	for i, book := range books {
// 		if book.ID == id {
// 			index = i
// 			break
// 		}
// 	}

// 	if index == -1 {
// 		http.Error(w, "Book not found", http.StatusNotFound)
// 		return
// 	}

// 	books[index] = updatedBook

// 	err = saveBooksToFile()
// 	if err != nil {
// 		http.Error(w, "Ошибка при сохранении в файл", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// func validateID(param string) (int, error) {
// 	id, err := strconv.Atoi(param)
// 	if err != nil {
// 		return 0, fmt.Errorf("Invalid ID: %v", err)
// 	}
// 	return id, nil
// }
