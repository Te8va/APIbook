package routers

import (
	"github.com/julienschmidt/httprouter"

	"github.com/Te8va/APIbook/internal/app/handler"
	"github.com/Te8va/APIbook/internal/app/middleware"
)

func RegisterHandlers(router *httprouter.Router, bookHandler *handler.Book) {
	router.GET("/books/:id", middleware.Log(bookHandler.GetBookByIDHandler))
	router.POST("/books", middleware.Log(bookHandler.AddBookHandler))
	router.DELETE("/books/:id", middleware.Log(bookHandler.DeleteBookHandler))
	router.PUT("/books/:id", middleware.Log(bookHandler.UpdateBookHandler))
}
