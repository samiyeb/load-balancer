package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Hello world, the web server

	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request from %s\n", r.RemoteAddr)
		fmt.Printf("GET / HTTP/1.1 %s\n", r.RemoteAddr)
		fmt.Printf("Host: localhost %s\n", r.RemoteAddr)
		fmt.Printf("User-Agent: curl/7.85.0 %s\n", r.RemoteAddr)
		fmt.Printf("Accept: */* %s\n", r.RemoteAddr)
	}

	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

