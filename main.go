package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var tmpl *template.Template // Used for executing html templates.

var (
	db    *gorm.DB
	HOST  = os.Getenv("DOTS_HOST")
	USER  = os.Getenv("DOTS_USER")
	PWORD = os.Getenv("DOTS_PWORD")
	NAME  = os.Getenv("DOTS_NAME")
	PORT  = os.Getenv("DOTS_PORT")
)

// main is the main execution flow.
func main() {

	// Connect to database.
	var err error
	db, err = dbConnect()
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()

	// Create file servers to handle file requests for js and html.
	mux.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("s/script/js"))))

	// Define all the valid routes, and their respective handlers.
	mux.HandleFunc("/api/write/", writeHandler)
	mux.HandleFunc("/api/read/", readHandler)
	mux.HandleFunc("/", indexHandler)

	// Serve all incoming requests.
	http.ListenAndServe("localhost:8080", mux)

}

// indexHandler handles requests made to the index page, and any undefined routes.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	parseAndExecute(w, "index.html", nil)

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

// dbConnect initialises a connection with the database and returns a reference to a gorm.DB.
func dbConnect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Australia/Sydney", HOST, USER, PWORD, NAME, PORT)
	for i := 0; i < 5; i++ {
		gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Printf("connected to database: %s", NAME)
			return gormDB, err
		}
		log.Printf("failed to connect to %s: attempt %d", NAME, i+1)
	}
	return new(gorm.DB), fmt.Errorf("exceeded max retries and couldn't connect to database")

}
