package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Log(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t := time.Now()
		log.Default().Println("Request to", r.URL, "method", r.Method)
		next(w, r, ps)
		log.Default().Println("Time passed", time.Since(t))
	}
}
