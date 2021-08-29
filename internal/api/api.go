package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type API struct {
	r *httprouter.Router
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func New() *API {
	r := httprouter.New()
	r.GET("/v1/health", healthHandler)

	a := API{
		r: r,
	}
	return &a
}

func healthHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "ok!\n")
}
