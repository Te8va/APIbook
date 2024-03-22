package main

import (
	"net/http"

	"github.com/Te8va/APIbook/internal/app/server/handler"
	"github.com/Te8va/APIbook/internal/app/server/repository"
	"github.com/Te8va/APIbook/internal/app/server/routers"
	"github.com/Te8va/APIbook/internal/app/server/services"
	logging "github.com/Te8va/APIbook/internal/pkg/logger"
)

const port = ":8080"

func main() {

	bookRep := repository.NewFileBookRepository("books")

	bookSrv := services.NewBookService(bookRep)

	bookHandler := handler.NewBookHandler(bookSrv)

	router := routers.NewRouter()
	router.RegisterHandlers(bookHandler)

	logging.Logger().Info("Server is running on port", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		logging.Logger().Error("Error while starting the server:", err)
	}
}
