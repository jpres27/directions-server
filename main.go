package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Directions struct {
	Destination string
	How         string
}

func main() {
	fmt.Println("Server starting...")
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		// TODO: figure out how to connect this with a real database and use that
		directions := map[string][]Directions{
			"Directions": {
				{Destination: "Stacks", How: "Idk"},
				{Destination: "Study spaces", How: "Idk"},
				{Destination: "Showers", How: "Idk"},
			},
		}

		tmpl.Execute(w, directions)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		destination := r.PostFormValue("destination")
		how := r.PostFormValue("how")
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "directions-list-element", Directions{Destination: destination, How: how})
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/submit/", h2)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
