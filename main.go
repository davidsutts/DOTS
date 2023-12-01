package main

import (
	"log"
	"net/http"
)

func main () {

	mux := http.NewServeMux()

	mux.HandleFunc("/api/write", writeHandler)
	mux.HandleFunc("/", indexHandler)

	http.ListenAndServe("localhost:8080", mux)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Write([]byte("Hello Web"))
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Write([]byte("this is the write route"))
}
