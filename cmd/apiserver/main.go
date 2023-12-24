package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Te8va/APIbook/internal/app/handler"
	"github.com/Te8va/APIbook/internal/app/middleware"
	"github.com/Te8va/APIbook/internal/app/repository"
	"github.com/Te8va/APIbook/internal/app/service"
)

// @title Books API
// @version 1.0
// @description Сервис хранения информации о книгах

// @host localhost:8080
// @BasePath /

// @Tag.name Books
// @Tag.description Группа запросов для управления списком книг

// @Schemes http

func main() {

	// var bookRep domain.BookRepository = repository.NewFileBookRepository("books.json")

	bookRep := repository.NewFileBookRepository("M:\\GoLang\\APIbook\\cmd\\apiserver\\books") // TODO: use path like that "./books"

	bookSrv := service.NewBookService(bookRep)

	bookHandler := handler.NewBookHandler(bookSrv)

	router := httprouter.New()

	// TODO: make a func to register handlers (if you want)
	router.GET("/books/:id", middleware.Log(bookHandler.GetBookByIDHandler))
	router.POST("/books", middleware.Log(bookHandler.AddBookHandler))
	router.DELETE("/books/:id", middleware.Log(bookHandler.DeleteBookHandler))
	router.PUT("/books/:id", middleware.Log(bookHandler.UpdateBookHandler))

	router.GET("/swagger/:any", bookHandler.SwaggerHandler)

	fmt.Println("Сервер запущен на порту 8080") // TODO: use logger instead of println (zap/logrus)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
