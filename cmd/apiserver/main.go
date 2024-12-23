package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Te8va/APIbook/internal/app/server/config"
	bookDomain "github.com/Te8va/APIbook/internal/app/server/domain"
	"github.com/Te8va/APIbook/internal/app/server/handler"
	"github.com/Te8va/APIbook/internal/app/server/repository"
	"github.com/Te8va/APIbook/internal/app/server/routers"
	"github.com/Te8va/APIbook/internal/app/server/services"
	"github.com/Te8va/APIbook/internal/pkg/logger"
)

func main() {
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Logger().Fatal("Error loading configuration:", err)
	}

	var bookRep bookDomain.BookRepository
	var p *pgxpool.Pool

	if !cfg.UseFile {
		logger.Logger().Info("Using PostgreSQL-based repository.")

		migrator, err := repository.NewMigrator(cfg.MigrationsPath, cfg.PostgresConn)
		if err != nil {
			logger.Logger().Fatal("Error to create migrator:", err)
		}

		logger.Logger().Info("Running migrations...")
		if err := repository.ApplyMigrations(migrator); err != nil {
			logger.Logger().Fatal("Error applying migrations:", err)
		}

		p, err = repository.NewPgxpool(cfg.PostgresConn)
		logger.Logger().Info("Postgres connection string: ", cfg.PostgresConn)
		if err != nil {
			logger.Logger().Fatal("Error creating PostgreSQL connection pool:", err)
		}

		if p != nil {
			defer p.Close()
		}

		bookRep = repository.NewPostgresBookRepository(p)
	} else {
		var err error
		logger.Logger().Info("Using file-based repository.")
		bookRep, err = repository.NewFileBookRepository("books")
		if err != nil {
			logger.Logger().Error("Failed to initialize file repository:", err)
		}
	}

	bookSrv := services.NewBookService(bookRep)
	bookHandler := handler.NewBookHandler(bookSrv)

	router := routers.NewRouter()
	router.RegisterHandlers(bookHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort),
		Handler: router,
	}

	go func() {
		logger.Logger().Info("Server is running on port", cfg.ServicePort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger().Fatal("ListenAndServe error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Logger().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger().Fatal("Server forced to shutdown:", err)
	}

	if p != nil {
		logger.Logger().Info("Closing database connections...")
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer closeCancel()

		done := make(chan struct{})
		go func() {
			p.Close()
			close(done)
		}()

		select {
		case <-done:
			logger.Logger().Info("Database connections closed successfully.")
		case <-closeCtx.Done():
			logger.Logger().Warn("Timeout exceeded while waiting for database connections to close.")
		}
	}

	logger.Logger().Info("Server exiting")
}
