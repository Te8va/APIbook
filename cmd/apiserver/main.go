package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Te8va/APIbook/internal/app/handler"
	"github.com/Te8va/APIbook/internal/app/repository"
	"github.com/Te8va/APIbook/internal/app/routers"
	"github.com/Te8va/APIbook/internal/app/services"
	logging "github.com/Te8va/APIbook/internal/pkg/logger"
)

func main() {

	bookRep := repository.NewFileBookRepository("books")

	bookSrv := services.NewBookService(bookRep)

	bookHandler := handler.NewBookHandler(bookSrv)

	router := httprouter.New()

	routers.RegisterHandlers(router, bookHandler)

	logging.Logger().Info("Server is running on port 8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logging.Logger().Error("Error while starting the server:", err)
	}
}
