package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Print the source IP
	fmt.Printf("Received request from %s\n", r.RemoteAddr)

	// Print the HTTP method, URI, and HTTP version
	fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)

	// Print the Host header
	fmt.Printf("Host: %s\n", r.Host)

	// Print the User-Agent header
	fmt.Printf("User-Agent: %s\n", r.UserAgent())

	// Print the Accept header (if available)
	fmt.Printf("Accept: %s\n", r.Header.Get("Accept"))

	// Send a simple response back
	fmt.Fprintln(w, "Request received!")
}

func main() {
	// Setup HTTP server with the handler function
	http.HandleFunc("/", handler)

	// Start the server on port 8080
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


