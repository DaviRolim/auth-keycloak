package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *Config) routes() http.Handler {
	specialContentRole := "vip"
	router := mux.NewRouter().StrictSlash(true)

	router.Use(commonMiddleware)

	router.HandleFunc("/login", app.Login).Methods(http.MethodPost)
	router.Handle("/greet", app.Protect(http.HandlerFunc(app.Greet))).Methods(http.MethodGet)
	router.Handle("/greet-vip", app.ProtectForRole(http.HandlerFunc(app.GreetVip), specialContentRole)).Methods(http.MethodGet)

	http.Handle("/", router)
	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
