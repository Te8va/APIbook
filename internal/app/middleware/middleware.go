package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/Te8va/APIbook/internal/app/logging"
)

func Log(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logging.Logger.Info("Request", zap.String("URL", r.URL.Path), zap.String("Method", r.Method))

		defer func() {
			if err := recover(); err != nil {
				logging.Logger.Error("An error occurred in the request handler", zap.Any("error", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next(w, r, ps)
	}
}
