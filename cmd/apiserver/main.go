package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/Te8va/APIbook/internal/app/handler"
	"github.com/Te8va/APIbook/internal/app/logging"
	"github.com/Te8va/APIbook/internal/app/repository"
	"github.com/Te8va/APIbook/internal/app/routers"
	"github.com/Te8va/APIbook/internal/app/services"
)

func main() {

	logger, err := logging.GetLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("Error while syncing the logger:", zap.Error(err))
		}
	}()

	bookRep := repository.NewFileBookRepository("books")

	bookSrv := services.NewBookService(bookRep)

	bookHandler := handler.NewBookHandler(bookSrv)

	router := httprouter.New()

	routers.RegisterHandlers(router, bookHandler)

	logger.Info("Server is running on port 8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Error("Error while starting the server:", zap.Error(err))
	}
}
