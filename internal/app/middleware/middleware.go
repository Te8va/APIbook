package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Log(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel() // TODO: use logger, do not create new context
		next(w, r.WithContext(ctx), ps)
	}
}
