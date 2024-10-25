package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	loadBalancerPort = ":80"
	healthCheckPeriod = 10 * time.Second
	healthCheckPath = "/health" // Path for health check
)

var backendServers = []string{
	"http://localhost:8080",
	"http://localhost:8081",
	"http://localhost:8082",
}
var availableServers []string
var currentServer int
var mu sync.Mutex

func handleConnection(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	if len(availableServers) == 0 {
		http.Error(w, "No available backend servers", http.StatusServiceUnavailable)
		mu.Unlock()
		return
	}

	server := availableServers[currentServer]
	currentServer = (currentServer + 1) % len(availableServers)
	mu.Unlock()

	resp, err := forwardRequest(req, server)
	if err != nil {
		http.Error(w, "Error forwarding request to backend server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyResponse(w, resp)
}

func forwardRequest(req *http.Request, serverURL string) (*http.Response, error) {
	client := &http.Client{}
	newReq, err := http.NewRequest(req.Method, serverURL+req.RequestURI, req.Body)
	if err != nil {
		return nil, err
	}
	newReq.Header = req.Header
	return client.Do(newReq)
}

func copyResponse(w http.ResponseWriter, resp *http.Response) {
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func healthCheckRoutine() {
	for {
		mu.Lock()
		newAvailableServers := []string{}

		for _, server := range backendServers {
			resp, err := http.Get(server + healthCheckPath)
			if err == nil && resp.StatusCode == http.StatusOK {
				newAvailableServers = append(newAvailableServers, server)
				resp.Body.Close()
			} else if resp != nil {
				resp.Body.Close()
			}
		}

		availableServers = newAvailableServers
		mu.Unlock()

		time.Sleep(healthCheckPeriod)
	}
}

func main() {
	availableServers = append(availableServers, backendServers...)
	http.HandleFunc("/", handleConnection)
	fmt.Println("Load balancer started, listening on port", loadBalancerPort)

	go healthCheckRoutine()

	err := http.ListenAndServe(loadBalancerPort, nil)
	if err != nil {
		log.Fatalf("Error starting load balancer: %v", err)
		os.Exit(1)
	}
}



