package routers

import (
	"github.com/julienschmidt/httprouter"

	"github.com/Te8va/APIbook/internal/app/server/handler"
	"github.com/Te8va/APIbook/internal/app/server/middleware"
)

type Router struct {
	*httprouter.Router
}

func NewRouter() *Router {
	return &Router{httprouter.New()}
}

func (r *Router) RegisterHandlers(bookHandler *handler.Book) {
	r.GET("/books/:id", middleware.Log(bookHandler.GetBookByIDHandler))
	r.POST("/books", middleware.Log(bookHandler.AddBookHandler))
	r.DELETE("/books/:id", middleware.Log(bookHandler.DeleteBookHandler))
	r.PUT("/books/:id", middleware.Log(bookHandler.UpdateBookHandler))
}
