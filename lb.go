package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "sync"
    "time"
)

// Server struct for each backend server
type Server struct {
    URL    *url.URL
    Alive  bool
    mutex  sync.RWMutex
}

// LoadBalancer struct to manage servers and requests
type LoadBalancer struct {
    Servers          []*Server
    RoundRobinCount  int
    mutex            sync.Mutex
    HealthCheckURL   string
    HealthCheckInterval time.Duration
}

// NewServer creates a new server instance
func NewServer(serverURL string) *Server {
    parsedURL, _ := url.Parse(serverURL)
    return &Server{
        URL:   parsedURL,
        Alive: true,
    }
}

// SetAlive updates the alive status of the server
func (s *Server) SetAlive(alive bool) {
    s.mutex.Lock()
    s.Alive = alive
    s.mutex.Unlock()
}

// IsAlive checks if the server is alive
func (s *Server) IsAlive() bool {
    s.mutex.RLock()
    alive := s.Alive
    s.mutex.RUnlock()
    return alive
}

// NewLoadBalancer initializes the load balancer
func NewLoadBalancer(servers []string, healthCheckInterval time.Duration) *LoadBalancer {
    var serverList []*Server
    for _, serverURL := range servers {
        server := NewServer(serverURL)
        serverList = append(serverList, server)
    }
    lb := &LoadBalancer{
        Servers:          serverList,
        HealthCheckURL:   "/",
        HealthCheckInterval: healthCheckInterval,
    }
    return lb
}

// GetNextAvailableServer returns the next server in round-robin
func (lb *LoadBalancer) GetNextAvailableServer() *Server {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()

    serverCount := len(lb.Servers)
    for i := 0; i < serverCount; i++ {
        lb.RoundRobinCount = (lb.RoundRobinCount + 1) % serverCount
        nextServer := lb.Servers[lb.RoundRobinCount]

        if nextServer.IsAlive() {
            return nextServer
        }
    }
    return nil
}

// ServeHTTP handles incoming requests and forwards them to the backend server
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    server := lb.GetNextAvailableServer()
    if server == nil {
        http.Error(w, "No available servers", http.StatusServiceUnavailable)
        return
    }

    proxyURL := server.URL.ResolveReference(r.URL)
    req, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
    if err != nil {
        http.Error(w, "Failed to create request", http.StatusInternalServerError)
        return
    }

    req.Header = r.Header

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        server.SetAlive(false)
        http.Error(w, "Server error", http.StatusServiceUnavailable)
        return
    }
    defer resp.Body.Close()

    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)

    fmt.Printf("Forwarded request to %s; received response: %d\n", server.URL, resp.StatusCode)
}

// healthCheck pings the server to see if it's alive
func (lb *LoadBalancer) healthCheck(server *Server) {
    resp, err := http.Get(server.URL.String() + lb.HealthCheckURL)
    if err != nil || resp.StatusCode != http.StatusOK {
        server.SetAlive(false)
        fmt.Printf("Server %s is down\n", server.URL)
    } else {
        server.SetAlive(true)
        fmt.Printf("Server %s is healthy\n", server.URL)
    }
}

// StartHealthChecks periodically checks each server's health
func (lb *LoadBalancer) StartHealthChecks() {
    ticker := time.NewTicker(lb.HealthCheckInterval)
    go func() {
        for range ticker.C {
            for _, server := range lb.Servers {
                lb.healthCheck(server)
            }
        }
    }()
}

func main() {
    servers := []string{
        "http://localhost:8080",
        "http://localhost:8081",
        "http://localhost:8082",
    }
    lb := NewLoadBalancer(servers, 10*time.Second)
    lb.StartHealthChecks()

    fmt.Println("Starting load balancer on port 80")
    log.Fatal(http.ListenAndServe(":80", lb))
}




