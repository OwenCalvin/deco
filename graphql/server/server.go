package server

import (
	"bytes"
	"dego/graphql/definition"
	"dego/graphql/executor"
	"dego/graphql/language/parser"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var schema definition.Schema

type Context struct {
}

type Request struct {
	OperationName string      `json:"operationName"`
	Variables     interface{} `json:"variables"`
	Query         string      `json:"query"`
	Mutation      string      `json:"mutation"`
}

func (r *Request) String() string {
	return r.Query
}

func parseRequest(request string) interface{} {
	req := &Request{}
	json.Unmarshal([]byte(request), req)
	fmt.Println(req)
	doc, _ := parser.Parse(parser.ParseParams{
		Source: req.Query,
		Options: parser.ParseOptions{
			NoLocation: false,
			NoSource:   false,
		},
	})

	return executor.Execute(&executor.ExecutionParams{
		Schema: schema,
		AST:    doc,
	})
}

func parseBody(r *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return buf.String()
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	res := parseRequest(parseBody(r))
	resJson, _ := json.Marshal(res)

	fmt.Fprint(w, string(resJson))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	switch r.Method {
	case "POST":
		handlePost(w, r)
	case "GET":
		handleGet(w, r)
	default:
		fmt.Fprintf(w, "ok")
	}
}

func Serve(sch definition.Schema, path string, port string) {
	schema = sch
	http.HandleFunc(path, handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
