package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

type Peel struct {
	Id   int
	Body string
}

func StartHTTPServer() {
	fmt.Println("Server up and running at http://localhost:1337...")
	http.ListenAndServe(":1337", nil)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", "dbname=kananas_development sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		peels := []Peel{}

		rows, err := db.Query("SELECT * FROM peels")
		for rows.Next() {
			peel := Peel{}
			if err := rows.Scan(&peel.Id, &peel.Body); err != nil {
				log.Fatal(err)
			}
			peels = append(peels, peel)
		}

		db.Close()

		tmpl, _ := template.ParseFiles("app/templates/index.html")
		tmpl.Execute(w, peels)
	})

	http.HandleFunc("/peels", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", "dbname=kananas_development sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		body := r.FormValue("body")
		db.Exec("INSERT INTO peels (body) VALUES ($1)", body)

		db.Close()
		http.Redirect(w, r, "/", http.StatusFound)
	})

	StartHTTPServer()
}
