package server

import (
	"bytes"
	"deco/graphql/definition"
	"deco/graphql/language/parser"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var Schema definition.Schema

type Context struct {
}

type Request struct {
	OperationName string      `json:"operationName"`
	Variables     interface{} `json:"variables"`
	Query         string      `json:"query"`
	Mutation      string      `json:"mutation"`
}

type Response struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func (r *Request) String() string {
	return r.Query
}

func parseRequest(request string) (res interface{}, err error) {
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

	return Schema.Execute(doc)
}

func parseBody(r *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return buf.String()
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	res, err := parseRequest(parseBody(r))

	finalRes := Response{
		Data:  res,
		Error: err,
	}
	resJSON, _ := json.Marshal(finalRes)

	fmt.Fprint(w, string(resJSON))
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
	Schema = sch
	http.HandleFunc(path, handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
