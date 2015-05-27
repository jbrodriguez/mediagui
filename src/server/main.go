package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var Version string

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/Volumes/Users/kayak/code/src/jbrodriguez/mediagui/target/build/index.html")
}

func main() {
	log.Printf("starting mediagui v%s ...", Version)

	router := httprouter.New()

	router.HandlerFunc("GET", "/", index)
	router.ServeFiles("/app/*filepath", http.Dir("/Volumes/Users/kayak/code/src/jbrodriguez/mediagui/target/build//app"))

	log.Fatal(http.ListenAndServe(":7623", router))
}
