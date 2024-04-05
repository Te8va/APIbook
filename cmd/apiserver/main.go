package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"

	bookDomain "github.com/Te8va/APIbook/internal/app/server/domain"
	"github.com/Te8va/APIbook/internal/app/server/handler"
	"github.com/Te8va/APIbook/internal/app/server/repository"
	"github.com/Te8va/APIbook/internal/app/server/routers"
	"github.com/Te8va/APIbook/internal/app/server/services"
	logging "github.com/Te8va/APIbook/internal/pkg/logger"
)

const port = ":8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		logging.Logger().Fatal("Error loading .env file")
	}

	useFile, _ := os.LookupEnv("USE_FILE")
	var bookRep bookDomain.BookRepository
	if useFile != "true" {
		err = repository.ApplyMigrations("file://migrations", "postgres://go-book:go-book@localhost:5432/go-book?sslmode=disable")
		if err != nil {
			logging.Logger().Error("Error to apply migrations:", err)
		}

		p, err := repository.NewPgxpool("postgres://go-book:go-book@localhost:5432/go-book?sslmode=disable")
		if err != nil {
			logging.Logger().Error("Error to create PostgreSQL connection pool:", err)
			return
		}

		defer p.Close()

		bookRep = repository.NewPostgresBookRepository(p)

	} else {
		bookRep = repository.NewFileBookRepository("books")
	}

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
