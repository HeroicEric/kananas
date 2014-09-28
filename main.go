package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type Peel struct {
	Id   int
	Body string
}

func startHTTPServer() {
	fmt.Println("Server up and running at http://localhost:1337...")
	http.ListenAndServe(":1337", nil)
}

func withMetrics(l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		next.ServeHTTP(w, r)
		l.Printf("%s %s took %s", r.Method, r.URL, time.Since(began))
	})
}

func indexHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		peels := []Peel{}

		rows, _ := db.Query("SELECT * FROM peels")
		for rows.Next() {
			peel := Peel{}
			if err := rows.Scan(&peel.Id, &peel.Body); err != nil {
				log.Fatal(err)
			}
			peels = append(peels, peel)
		}

		tmpl, _ := template.ParseFiles("app/templates/index.html")
		tmpl.Execute(w, peels)
	})
}

func createPeelHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("body")

		db.Exec("INSERT INTO peels (body) VALUES ($1)", body)

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func main() {
	db, err := sql.Open("postgres", "dbname=kananas_development sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "", 0)

	http.Handle("/", withMetrics(logger, indexHandler(db)))
	http.Handle("/peels", createPeelHandler(db))

	startHTTPServer()
}
