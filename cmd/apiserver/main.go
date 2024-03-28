package main

import (
	"log"
	"net/http"

	"github.com/Te8va/APIbook/internal/app/server/handler"
	"github.com/Te8va/APIbook/internal/app/server/repository"
	"github.com/Te8va/APIbook/internal/app/server/routers"
	"github.com/Te8va/APIbook/internal/app/server/services"
	logging "github.com/Te8va/APIbook/internal/pkg/logger"
)

const port = ":8080"

func main() {

	err := repository.ApplyMigrations("file://migrations", "postgres://go-book:go-book@localhost:5432/go-book?sslmode=disable")
	if err != nil {
		log.Println(err)
	}

	p, err := repository.NewPgxpool("postgres://go-book:go-book@localhost:5432/go-book?sslmode=disable")
	if err != nil {
		log.Println(err)
		return
	}

	defer p.Close()

	bookRep := repository.NewPostgresBookRepository(p)

	// bookRep := repository.NewFileBookRepository("books")

	bookSrv := services.NewBookService(bookRep)

	bookHandler := handler.NewBookHandler(bookSrv)

	router := routers.NewRouter()
	router.RegisterHandlers(bookHandler)

	logging.Logger().Info("Server is running on port", port)

	err = http.ListenAndServe(port, router)
	if err != nil {
		logging.Logger().Error("Error while starting the server:", err)
	}
}
