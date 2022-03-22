package main

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func checkStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}

func executeCode(w http.ResponseWriter, r *http.Request) {
	var t map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		log.Println(err)
	}

	result := pythonExecution(t)
	fmt.Fprint(w, result)
}

func pythonExecution(t map[string]interface{}) string {
	var r string;
	var err error;

	for c := range t {
		s := strings.Split(c, ":")
		switch s[0] {
		case "assignation":
			r, err = doAssign(t[c])
		case "operation":
			r, err = doOperation(s[1], t[c])
		}
	}

	if err != nil {
		return "An error ocurred"
	}

	return r
}

func doAssign(val interface{}) (string, error) {
	var r string;
	var err error;
	s := val.([]interface{})
	for _, v := range s {
		switch v.(type) {
		case string:
			r = v.(string)
		case map[string]interface{}:
			m := v.(map[string]interface{})
			for opr, num := range m {
				opr = strings.Split(opr, ":")[1]
				r, err = doOperation(opr, num)
			}
		default:
			log.Printf("Type: %T\n", v)
			r = "default"
		}
	}
	if err != nil {
		return "Error in assign", err
	}

	return r, nil
}

func doOperation(opr string, i interface{}) (string, error) {
	s := i.([]interface{})
	var n1, n2 string
	var r int;
	for _, n := range s {
		log.Printf("Number type and value: %v %T\n", n, n)
		switch n.(type) {
		case string:
			if n1 == "" {
				n1 = n.(string)
			} else {
				n2 = n.(string)
			}
		default:
			doOperation("add", n)
		}
	}

	i1, err1 := strconv.Atoi(n1)
	i2, err2 := strconv.Atoi(n2)

	if err1 != nil {
		return "", err1
	} else if err2 != nil {
		return "", err2
	}

	switch opr {
	case "add":
		r = i1 + i2
	case "sub":
		r = i1 - i2
	case "mul":
		r = i1 * i2
	case "div":
		r = i1 / i2
	}

	log.Println(r)
	return strconv.Itoa(r), nil
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	router.Get("/", checkStatus)
	router.Post("/", executeCode)
	log.Fatal(http.ListenAndServe(":8080", router))
}