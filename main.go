package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func StartHTTPServer() {
	fmt.Println("Server up and running at http://localhost:1337...")
	http.ListenAndServe(":1337", nil)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("app/templates/index.html")

		tmpl.Execute(w, nil)
	})

	StartHTTPServer()
}
