package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func main () {

	mux := http.NewServeMux()

	mux.HandleFunc("/api/write", writeHandler)
	mux.HandleFunc("/", indexHandler)

	http.ListenAndServe("localhost:8080", mux)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	
	var err error
	tmpl, err = template.ParseFiles("t/index.html")
	if err != nil {
		log.Printf("couldn't parse t/index.html: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("couldn't execute index.html: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Write([]byte("this is the write route"))
}
