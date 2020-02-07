package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const DefaultPort = "8080"

var headerPrefix string

var tmpl *template.Template

func getServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port != "" {
		return port
	}

	return DefaultPort
}

func HeaderHandler(writer http.ResponseWriter, request *http.Request) {

	log.Println("=====================================")
	log.Println("Handling the request made to " + request.URL.Path + " to client (" + request.RemoteAddr + ")")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
	log.Println(prettyPrint(request.Header))
	fmt.Fprint(writer, table(request.Header))
	log.Println("=====================================")
}

func main() {

	paths := []string{
		"table.html",
	}
	tmpl = template.Must(template.New("table.html").ParseFiles(paths...))
	log.Println("templates parsed, ready to use")

	headerPrefix = os.Getenv("HEADER_PREFIX")
	if headerPrefix != "" {
		log.Println("Headers will be highlighed if starts with: ", headerPrefix)
	}

	log.Println("starting server, listening on port " + getServerPort())

	http.HandleFunc("/", HeaderHandler)
	http.ListenAndServe(":"+getServerPort(), nil)
}

func table(i http.Header) string {

	dataArray := HeaderToArray(i)

	var buffer bytes.Buffer
	err := tmpl.Execute(&buffer, dataArray)
	if err != nil {
		log.Fatalln(err)
	}

	result := buffer.String()

	return result
}

type Header struct {
	Key   string
	Value string
	CC    bool
}

func HeaderToArray(header http.Header) []Header {
	res := make([]Header, 0)

	for name, values := range header {
		for _, value := range values {
			d := Header{Key: name, Value: value}
			if headerPrefix != "" && strings.HasPrefix(name, headerPrefix) {
				d.CC = true
			}
			res = append(res, d)
		}
	}
	return res
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
