# Go Load Balancer

This repository contains a simple HTTP load balancer written in Go, which distributes incoming requests to multiple backend servers using a round-robin algorithm. The load balancer performs periodic health checks on backend servers and only forwards requests to healthy servers.

## Features

- **Round-Robin Load Balancing**: Distributes requests among available backend servers in a rotating order.
- **Health Checks**: Periodically checks backend server health, only sending requests to healthy servers.

## Requirements

- Go 1.16 or higher
- Three backend servers (can be run using Python’s built-in HTTP server)

## Setup

1. Clone the repository:

    ```bash
    git clone https://github.com/your-username/go-load-balancer.git
    cd go-load-balancer
    ```

2. Build the project:

    ```bash
    go build -o lb lb.go
    ```

3. Start the backend servers on different ports, for example, using Python’s HTTP server:

    ```bash
    # Terminal 1
    mkdir -p server8080
    echo "<html><body>Hello from server 8080</body></html>" > server8080/index.html
    python3 -m http.server 8080 --directory server8080

    # Terminal 2
    mkdir -p server8081
    echo "<html><body>Hello from server 8081</body></html>" > server8081/index.html
    python3 -m http.server 8081 --directory server8081

    # Terminal 3
    mkdir -p server8082
    echo "<html><body>Hello from server 8082</body></html>" > server8082/index.html
    python3 -m http.server 8082 --directory server8082
    ```

## Running the Load Balancer

Run the load balancer executable on port 80:

```bash
sudo ./lb

