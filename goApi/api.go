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

type User struct {
	Uid      string   `json:"uid"`
	Name     string   `json:"name"`
	Password string   `json:"pass"`
	Scripts  []Script `json:"scripts"`
	Dtype    []string `json:"dgraph.type"`
}

type Script struct {
	Uid      string   `json:"uid"`
	Name     string   `json:"name"`
	Code     string   `json:"code"`
	List     string   `json:"nodeList"`
	Drawflow string   `json:"drawflow"`
	Dtype    []string `json:"dgraph.type"`
}

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

func saveData(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	errorCheck(err)

	c, err := dgo.DialCloud("https://blue-surf-580096.us-east-1.aws.cloud.dgraph.io/graphql", "ZjAyNGJhZTc4ZmIxMTVkNTM1NmQ3OGQ1YzRkMjAyNDQ=")
	errorCheck(err)
	defer c.Close()
	client := dgo.NewDgraphClient(api.NewDgraphClient(c))

	vars := make(map[string]string)
	vars["$usr"] = data["user"]
	vars["$pwd"] = data["password"]
	q := `
		query Usr($usr: string, $pwd: string) {
			getUsr(func: eq(name, $usr)) {
				uid
				checkpwd(pass, $pwd)
			}
		}
	`
	response,err := client.NewReadOnlyTxn().QueryWithVars(r.Context(), q, vars)
	errorCheck(err)
	var jsonResp map[string][]map[string]interface{}
	err = json.Unmarshal(response.Json, &jsonResp)
	errorCheck(err)
	if jsonResp["getUsr"][0]["checkpwd(pass)"] == false {
		fmt.Fprint(w, "Wrong password")
	} else {
		uid := jsonResp["getUsr"][0]["uid"].(string)
		m := &api.Mutation{
			CommitNow: true,
		}
		md := User{
			Uid:uid,
			Name:data["user"],
			Password:data["password"],
			Dtype:[]string{"User"},
			Scripts:[]Script{{
				Uid:"_:newScript",
				Name:data["name"],
				Code:data["script"],
				List:data["list"],
				Drawflow:data["nodes"],
				Dtype:[]string{"Scrìpt"},

			}},

		}
		d, err := json.Marshal(&md)
		errorCheck(err)
		m.SetJson = d
		result, err := client.NewTxn().Mutate(r.Context(), m)
		errorCheck(err)
		fmt.Fprintf(w, "Json: %v\nUids:%v\nMetrics:%v\n", result.Txn, result.Uids, result.Metrics)
	}
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