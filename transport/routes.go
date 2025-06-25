package transport

import (
	"net/http"

	"github.com/AdityaHarindar/Richpanel-Assignment/store"

	"github.com/gorilla/mux"
)

// NewRouter sets up the mux router with handlers wired to the store
func NewRouter(s store.Store, c store.Cache) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/posts", GetAllHandler(s, c)).Methods(http.MethodGet)
	r.HandleFunc("/posts/", PostHandler(s, c)).Methods(http.MethodPost)
	r.HandleFunc("/posts/{key}", GetHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts/{key}", PutHandler(s)).Methods(http.MethodPut)
	r.HandleFunc("/posts/{key}", DeleteHandler(s)).Methods(http.MethodDelete)
	r.HandleFunc("/health", HealthCheckHandler)
	return r
}
