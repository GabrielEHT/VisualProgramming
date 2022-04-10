package main

import (
	"os"
	"fmt"
	"log"
	"os/exec"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/dgraph-io/dgo/v210"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

type User struct {
	Uid      string   `json:"uid"`
	Name     string   `json:"name,omitempty"`
	Password string   `json:"pass,omitempty"`
	Scripts  []Script `json:"scripts"`
	Dtype    []string `json:"dgraph.type"`
}

type Script struct {
	Uid      string   `json:"uid"`
	Name     string   `json:"name,omitempty"`
	Code     string   `json:"code"`
	List     string   `json:"nodeList"`
	Drawflow string   `json:"drawflow"`
	Creator  User     `json:"creator,omitempty"`
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

func executeCode(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	errorCheck(err)

	file, err := os.Create("./script.py")
	errorCheck(err)

	text := []byte(data["data"])
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
	vars["$usr"] = chi.URLParam(r, "user")
	q := `
		query Usr($usr: string) {
			getUsr(func: eq(name, $usr)) {
				uid
			}
		}
	`
	response,err := client.NewReadOnlyTxn().QueryWithVars(r.Context(), q, vars)
	errorCheck(err)

	var jsonResp map[string][]map[string]interface{}
	err = json.Unmarshal(response.Json, &jsonResp)
	errorCheck(err)

	uid := jsonResp["getUsr"][0]["uid"].(string)
	m := &api.Mutation{
		CommitNow: true,
	}
	md := User{
		Uid:uid,
		Dtype:[]string{"User"},
		Scripts:[]Script{{
			Uid:"_:newScript",
			Name:data["name"],
			Code:data["script"],
			List:data["list"],
			Drawflow:data["nodes"],
			Creator:User{
				Uid:uid,
				Dtype:[]string{"User"},
			},
			Dtype:[]string{"Scrìpt"},

		}},
	}

	d, err := json.Marshal(&md)
	errorCheck(err)
	m.SetJson = d
	result, err := client.NewTxn().Mutate(r.Context(), m)
	errorCheck(err)

	fmt.Fprintf(w, "Txn: %v\nUids:%v\nMetrics:%v\n", result.Txn, result.Uids, result.Metrics)
}

func getScriptList(w http.ResponseWriter, r *http.Request) {
	c, err := dgo.DialCloud("https://blue-surf-580096.us-east-1.aws.cloud.dgraph.io/graphql", "ZjAyNGJhZTc4ZmIxMTVkNTM1NmQ3OGQ1YzRkMjAyNDQ=")
	errorCheck(err)
	defer c.Close()

	client := dgo.NewDgraphClient(api.NewDgraphClient(c))

	vars := make(map[string]string)
	vars["$usr"] = chi.URLParam(r, "user")
	q := `
		query Usr($usr: string) {
			getUsr(func: eq(name, $usr)) {
				scripts {
					name
				}
			}
		}
	`
	response,err := client.NewReadOnlyTxn().QueryWithVars(r.Context(), q, vars)
	errorCheck(err)

	var jsonResp map[string][]map[string]interface{}
	err = json.Unmarshal(response.Json, &jsonResp)
	errorCheck(err)

	err = json.NewEncoder(w).Encode(jsonResp["getUsr"][0]["scripts"])
	errorCheck(err)
}

func getScript(w http.ResponseWriter, r *http.Request) {
	c, err := dgo.DialCloud("https://blue-surf-580096.us-east-1.aws.cloud.dgraph.io/graphql", "ZjAyNGJhZTc4ZmIxMTVkNTM1NmQ3OGQ1YzRkMjAyNDQ=")
	errorCheck(err)
	defer c.Close()

	client := dgo.NewDgraphClient(api.NewDgraphClient(c))

	vars := make(map[string]string)
	vars["$st"] = chi.URLParam(r, "script")
	vars["$usr"] = chi.URLParam(r, "user")
	q := `
		query St($usr: string, $st: string) {
			getSt(func: eq(name, $usr)) {
				scripts @filter(eq(name, $st)){
					name
					code
					nodeList
					drawflow
				}
			}
		}
	`
	response,err := client.NewReadOnlyTxn().QueryWithVars(r.Context(), q, vars)
	errorCheck(err)

	var jsonResp map[string][]map[string][]interface{}
	err = json.Unmarshal(response.Json, &jsonResp)
	errorCheck(err)

	err = json.NewEncoder(w).Encode(jsonResp["getSt"][0]["scripts"][0])
	errorCheck(err)
}

func overwriteData(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	errorCheck(err)

	c, err := dgo.DialCloud("https://blue-surf-580096.us-east-1.aws.cloud.dgraph.io/graphql", "ZjAyNGJhZTc4ZmIxMTVkNTM1NmQ3OGQ1YzRkMjAyNDQ=")
	errorCheck(err)
	defer c.Close()

	client := dgo.NewDgraphClient(api.NewDgraphClient(c))

	vars := make(map[string]string)
	vars["$usr"] = chi.URLParam(r, "user")
	vars["$st"] = chi.URLParam(r, "script")
	q := `
		query St($usr: string, $st: string) {
			getUsrSt(func: eq(name, $usr)) {
				scripts @filter(eq(name, $st)) {
					uid
				}
			}
		}
	`
	response,err := client.NewReadOnlyTxn().QueryWithVars(r.Context(), q, vars)
	errorCheck(err)

	var jsonResp map[string][]map[string][]map[string]string
	err = json.Unmarshal(response.Json, &jsonResp)
	errorCheck(err)

	uid := jsonResp["getUsrSt"][0]["scripts"][0]["uid"]
	md := Script{
		Uid:uid,
		Code:data["script"],
		List:data["list"],
		Drawflow:data["nodes"],
		Dtype:[]string{"Script"},
	}
	m := &api.Mutation{
		CommitNow: true,
	}

	d, err := json.Marshal(&md)
	errorCheck(err)
	m.SetJson = d
	result, err := client.NewTxn().Mutate(r.Context(), m)
	errorCheck(err)

	fmt.Fprint(w, result.Metrics)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	router.Get("/", checkStatus)
	router.Post("/exec", executeCode) // cambiar a PUT
	router.Route("/users", func(router chi.Router) {
		router.Get("/{user}", getScriptList)
		router.Post("/{user}", saveData)
		router.Route("/{user}/{script}", func (router chi.Router) {
			router.Get("/", getScript)
			router.Post("/", overwriteData) // cambiar a PUT
		})
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}