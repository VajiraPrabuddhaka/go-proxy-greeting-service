package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get the base URL from the environment variable
	baseURL := os.Getenv("GREETING_SERVICE_URL")
	if baseURL == "" {
		log.Fatal("GREETING_SERVICE_URL environment variable is not set")
	}

	// Handler for the service
	http.HandleFunc("/proxy-greet", func(w http.ResponseWriter, r *http.Request) {
		// Get the name from the query parameters
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Name parameter is missing", http.StatusBadRequest)
			return
		}

		// Construct the full URL for the greeting service
		fullURL := fmt.Sprintf("%s/greet?name=%s", baseURL, name)

		// Make the GET request to the greeting service
		resp, err := http.Get(fullURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to make request: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read response body: %v", err), http.StatusInternalServerError)
			return
		}

		// Write the response from the greeting service to the client
		w.Write(body)
	})

	// Start the HTTP server
	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
