package main

import (
	"encoding/json"
	"net/http"
	"log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"fmt"
)

type Test struct {
	Title string
	Body string
	Moment int
}

type Test2 struct {
	Name string
	Instruction []string
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	t := Test{"A test", "Hello world", 1253}
	json.NewEncoder(w).Encode(t)
}

func executeCode(w http.ResponseWriter, r *http.Request) {
	var t map[string]interface{}
	//var t Test2
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Println(err)
	}
	for s := range t {
		log.Println(s)
	}
	json.NewEncoder(w).Encode(&t)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	router.Get("/", myHandler)
	router.Post("/", executeCode)
	log.Fatal(http.ListenAndServe(":8080", router))
}