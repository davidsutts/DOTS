package main

import (
	"log"
	"net/http"
	"strings"
)

type Input struct {
	ID     int64
	Prompt string
}

type Result struct {
	ID     int64
	Value  string
	Prompt int64
	Count  int
}

// writeHandler handles requests made to the api/write route, used for writing data to the datastore.
func writeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check for a valid url request /api/write/<prompt-key>/<data>
	url := strings.Split(strings.TrimPrefix(r.URL.String(), "/"), "/")
	if len(url) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the prompt.
	var input Input
	tx := db.First(&input, url[2])
	if tx.Error != nil {
		log.Printf("could not get prompt: %s", tx.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the count.
	var result Result
	tx = db.First(&result, "prompt = ? AND value = ?", input.ID, url[3])
	if tx.Error != nil {
		// Create a new result value
		result.Count = 1
		result.Prompt = input.ID
		result.Value = url[3]
		db.Save(&result)
		return
	}
	result.Count++
	db.Save(&result)
}
