package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	HealthEndpoint  = "/v1/health"
	VersionEndponit = "/v1/version"
)

func newRouter() *httprouter.Router {
	r := httprouter.New()
	r.PanicHandler = recovery()
	r.GET(HealthEndpoint, health())
	return r
}

type API struct {
	r *httprouter.Router
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func New() *API {
	r := newRouter()
	a := API{
		r: r,
	}
	return &a
}

type JSONResponse struct {
	Message string
	Code    int
}

func health() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		response(rw, JSONResponse{Message: "OK", Code: http.StatusOK}, http.StatusOK)
	}
}

func recovery() func(rw http.ResponseWriter, r *http.Request, err interface{}) {
	return func(rw http.ResponseWriter, r *http.Request, err interface{}) {
		response(rw, JSONResponse{Message: "InternalServerError", Code: http.StatusInternalServerError}, http.StatusInternalServerError)
	}
}

func response(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "\t")
		err := enc.Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
