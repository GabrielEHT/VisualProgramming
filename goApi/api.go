package main

import (
	"os"
	"fmt"
	"log"
	"os/exec"
	//"strings"
	//"strconv"
	"net/http"
	"encoding/json"
	//"google.golang.org/grpc"
	"github.com/go-chi/chi/v5"
	"github.com/dgraph-io/dgo/v210"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

func checkStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}

// crear función para los errores
func executeCode(w http.ResponseWriter, r *http.Request) {
	var t map[string]string
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("./script.py")
	if err != nil {
		log.Fatal(err)
	}
	text := []byte(t["data"])
	_, err = file.Write(text)
	if err != nil {
		log.Fatal(err)
	}
	c := exec.Command("python3", "script.py")
	result, err := c.CombinedOutput()
	resultStr := string(result)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	if resultStr == "" {
		fmt.Fprint(w, "Code executed without errors")
	} else {
		fmt.Fprintf(w, "%s", result)
	}
}

type Test struct {
	n string
	a int
}

func saveData(w http.ResponseWriter, r *http.Request) {
	c, err := dgo.DialCloud("https://blue-surf-580096.us-east-1.aws.cloud.dgraph.io/graphql", "ZjAyNGJhZTc4ZmIxMTVkNTM1NmQ3OGQ1YzRkMjAyNDQ=")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	client := dgo.NewDgraphClient(api.NewDgraphClient(c))
	opr := &api.Operation{}
	opr.Schema = `
		name: string @index(exact).
		age : int.
	`

	err = client.Alter(r.Context(), opr)

	if err != nil {
		log.Fatal(err)
	}

	m := &api.Mutation{
		CommitNow: true,
	}

	var j interface{}
	err = json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		log.Fatal(err)
	}
	cosa := Test{"Prueba", 1}
	d, err := json.Marshal(cosa)
	if err != nil {
		log.Fatal(err)
	}
	m.SetJson = d

	t, err := client.NewTxn().Mutate(r.Context(), m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, t)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	router.Get("/", checkStatus)
	router.Post("/scripts", saveData)
	router.Post("/exec", executeCode)
	log.Fatal(http.ListenAndServe(":8080", router))
}