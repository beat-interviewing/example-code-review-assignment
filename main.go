package main

import (
	"log"
	"net/http"

	"example.com/assignment/code-review/beatly"
)

func main() {
	r, err := beatly.NewBoltStore("beatly.db")
	if err != nil {
		log.Fatal(err)
	}
	s := beatly.NewService(r)
	h := beatly.Handler(s)

	http.ListenAndServe(":8080", h)
}
