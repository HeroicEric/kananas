package main

import (
	"fmt"
	"net/http"
)

func StartHTTPServer() {
	fmt.Println("Server up and running at http://localhost:1337...")
	http.ListenAndServe(":1337", nil)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world!")
	})

	StartHTTPServer()
}
