package graphql

import (
	"bytes"
	"dego/reflector"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

func parseRequest(request string) {
	req := &Request{}
	json.Unmarshal([]byte(request), req)
	fmt.Println(req)
	if req.Query != "" {

	} else if req.Mutation != "" {

	}
}

func parseBody(r *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return buf.String()
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	parseRequest(parseBody(r))
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

func New(path string, port string, types ...interface{}) {
	allReflections := reflector.ReflectTypes(types...)
	fmt.Println(allReflections)

	http.HandleFunc(path, handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
