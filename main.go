package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var tmpl *template.Template // Used for executing html templates.
var db *gorm.DB

// main is the main execution flow.
func main() {

	// Get env vars for database setup.
	HOST := os.Getenv("DOTS_HOST")
	USER := os.Getenv("DOTS_USER")
	PWORD := os.Getenv("DOTS_PWORD")
	NAME := os.Getenv("DOTS_NAME")
	PORT := os.Getenv("DOTS_PORT")

	// Connect to the database.
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Australia/Adelaide", HOST, USER, PWORD, NAME, PORT)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	mux := http.NewServeMux()

	// Create file servers to handle file requests for js and html.
	mux.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("s/script/js"))))

	// Define all the valid routes, and their respective handlers.
	mux.HandleFunc("/api/write/", writeHandler)
	mux.HandleFunc("/", indexHandler)

	// Serve all incoming requests.
	http.ListenAndServe("localhost:8080", mux)

}

// indexHandler handles requests made to the index page, and any undefined routes.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	parseAndExecute(w, "index.html", nil)

}

// writeHandler handles requests made to the api/write route, used for writing data to the datastore.
func writeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check for a valid url request /api/write/<data>
	log.Println(r.URL.String())
	url := strings.Split(strings.TrimPrefix(r.URL.String(), "/"), "/")
	if len(url) != 3 {
		log.Println(url[0])
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Print received data.
	log.Println(url[2])
}

// parseAndExecute parses the given file, and executes the template with the given data.
//
// Handles failures and writes header errors.
func parseAndExecute(w http.ResponseWriter, filename string, data ...interface{}) {
	// Parse the given file.
	var err error
	tmpl, err = template.ParseFiles("s/html/" + filename)
	if err != nil {
		log.Printf("couldn't parse /s/html/%s: %v", filename, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Execute the template.
	err = tmpl.ExecuteTemplate(w, filename, data)
	if err != nil {
		log.Printf("couldn't execute %s: %v", filename, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
