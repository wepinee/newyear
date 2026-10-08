package main

import (
	"log"
	"net/http"

	"github.com/wepinee/newyear"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/days", newyear.DaysHandler())

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
