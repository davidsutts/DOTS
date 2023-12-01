package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var tmpl *template.Template // Used for executing html templates.

// main is the main execution flow.
func main() {

	mux := http.NewServeMux()

	// Define all the valid routes, and their respective handlers.
	mux.HandleFunc("/api/write", writeHandler)
	mux.HandleFunc("/", indexHandler)

	// Serve all incoming requests.
	http.ListenAndServe("localhost:8080", mux)

}

// indexHandler handles requests made to the index page, and any undefined routes.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	parseAndExecute(w, "t/index.html", nil)

}

// writeHandler handles requests made to the api/write route, used for writing data to the datastore.
func writeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Write([]byte("this is the write route"))
}

// parseAndExecute parses the given file, and executes the template with the given data.
//
// Handles failures and writes header errors.
func parseAndExecute(w http.ResponseWriter, filepath string, data ...interface{}) {
	// Parse the given file.
	var err error
	tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		log.Printf("couldn't parse %s: %v", filepath, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the name of the file.
	s := strings.Split(filepath, "/")
	filename := s[len(s)-1]

	// Execute the template.
	err = tmpl.ExecuteTemplate(w, filename, data)
	if err != nil {
		log.Printf("couldn't execute %s: %v", filename, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
