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

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// crear función para los errores
func executeCode(w http.ResponseWriter, r *http.Request) {
	var t map[string]string
	err := json.NewDecoder(r.Body).Decode(&t)
	errorCheck(err)

	file, err := os.Create("./script.py")
	errorCheck(err)

	text := []byte(t["data"])
	_, err = file.Write(text)
	errorCheck(err)

	c := exec.Command("python3", "script.py")
	result, err := c.CombinedOutput()
	resultStr := string(result)
	errorCheck(err)

	if resultStr == "" {
		fmt.Fprint(w, "Code executed without errors")
	} else {
		fmt.Fprintf(w, "%s", result)
	}
}

type Person struct {
	Uid 	 string   `json:"uid,omitempty"`
	Name 	 string	  `json:"name,omitempty"`
	Age 	 int 	  `json:"age,omitempty"`
	Friends  []Person `json:"friends,omitempty"`
	OwnsPets []Animal `json:"ownsPets,omitempty"`
	Dtype 	 []string `json:"dgraph.type,omitempty"`
}

type Animal struct {
	uid   string   `json:"uid,omitempty"`
	name  string   `json:"name,omitempty"`
	dtype []string `json:"dgraph.type,omitempty"`
	owner []Person `json:"owner,omitempty"`
}

func saveData(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	errorCheck(err)
	fmt.Println(data)

	c, err := dgo.DialCloud("https://blue-surf-580096.us-east-1.aws.cloud.dgraph.io/graphql", "ZjAyNGJhZTc4ZmIxMTVkNTM1NmQ3OGQ1YzRkMjAyNDQ=")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	client := dgo.NewDgraphClient(api.NewDgraphClient(c))
	/*opr := &api.Operation{}
	opr.Schema = `
	    name: string @index(term).
	    age: int .
	    friends: [uid] .
	    ownsPets: [uid] .
	    owner: [uid] .

	    type Person {
	    	name
	    	age
	    	friends
	    	ownsPets
	    }

	    type Animal {
	    	name
	    	owner
	    }
	`

	err = client.Alter(r.Context(), opr)

	if err != nil {
		log.Fatal(err)
	}*/

	m := &api.Mutation{
		CommitNow: true,
	}

	cosa := Person{
		Uid: "_:mario",
		Name: "Mario",
		Age: 31,
		Dtype: []string{"Person"},
	}
	d, err := json.Marshal(&cosa)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Json:%s\n", d)
	fmt.Printf("Json:%s", string(d))
	m.SetJson = d

	t, err := client.NewTxn().Mutate(r.Context(), m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Json: %v\nUids:%v\nMetrics:%v\n", t.Json, t.Uids, t.Metrics)
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