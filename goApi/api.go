package main

import (
	"encoding/json"
	"net/http"
	"log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Test struct {
	Title string
	Body string
	Moment int
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	m := Test{"A test", "Hello world", 1253}
	json.NewEncoder(w).Encode(m)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	router.Get("/", myHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}