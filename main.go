package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"
)

// Tracks which server to send the next request to and uses a mutex to ensure the logic for selecting servers is thread-safe
type LoadBalancer struct {
	Current int
	Mutex   sync.Mutex
}

// Represents a backend server with a URL and a health status. The mutex ensures that the health status can be updated or checked safely across multiple requests.
type Server struct {
	URL       *url.URL
	IsHealthy bool
	Mutex     sync.Mutex
}

// When the load balancer receives a request, it forwards the request to the next available server using a reverse proxy.
// In Golang, the httputil package provides a built-in way to handle reverse proxying, and we will use it in our code through the ReverseProxy function:
func (s *Server) ReverseProxy() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(s.URL)
}

type Config struct {
	Port                string   `json:"port"`
	HealthCheckInterval string   `json:"healthCheckInterval"`
	Servers             []string `json:"servers"`
}

func loadConfig(file string) (Config, error) {
	var config Config

	// Read the contents of the config file
	data, err := os.ReadFile(file)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal JSON data into the Config struct
	err = json.Unmarshal(data, &config)
	if err != nil {
		// Return an empty config and the error if unmarshaling fails
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Return the successfully populated config
	return config, nil
}

// health check function that runs in given interval to check health of servers.
// This healthCheck function performs periodic health checks on a backend server
// using the HTTP HEAD request to see if the server is reachable and responding with a status code of 200 OK.
func healthChecks(s *Server, healthCheckInterval time.Duration) {
	// Ticker for periodic health checks
	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()

	// Create an HTTP client with a custom timeout
	client := &http.Client{
		Timeout: 5 * time.Second, // Adjust timeout as needed
	}

	// Runs the health check periodically at intervals of healthCheckInterval
	for range ticker.C {
		// Send an HTTP HEAD request to check if the server is up
		res, err := client.Head(s.URL.String())

		// Lock the server's mutex to update health status
		s.Mutex.Lock()

		if err != nil {
			fmt.Printf("Error checking %s: %v\n", s.URL, err)
			s.IsHealthy = false
		} else if res.StatusCode != http.StatusOK {
			fmt.Printf("%s is down (status code: %d)\n", s.URL, res.StatusCode)
			s.IsHealthy = false
		} else {
			s.IsHealthy = true
		}

		// Ensure the response body is closed
		if res != nil {
			res.Body.Close()
		}

		// Unlock the mutex after updating the status
		s.Mutex.Unlock()
	}
}

// round robin algorithm implementation to distribute load across servers
func (lb *LoadBalancer) getNextServer(servers []*Server) *Server {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	for i := 0; i < len(servers); i++ {
		idx := lb.Current % len(servers)
		nextServer := servers[idx]
		lb.Current++

		nextServer.Mutex.Lock()
		isHealthy := nextServer.IsHealthy
		nextServer.Mutex.Unlock()

		if isHealthy {
			return nextServer
		}
	}

	return nil
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	healthCheckInterval, err := time.ParseDuration(config.HealthCheckInterval)
	if err != nil {
		log.Fatalf("Invalid health check interval: %s", err.Error())
	}

	var servers []*Server
	for _, serverUrl := range config.Servers {
		u, _ := url.Parse(serverUrl)
		server := &Server{URL: u, IsHealthy: true}
		servers = append(servers, server)
		go healthChecks(server, healthCheckInterval)
	}

	lb := LoadBalancer{Current: 0}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server := lb.getNextServer(servers)
		if server == nil {
			http.Error(w, "No healthy server available", http.StatusServiceUnavailable)
			return
		}

		// adding this header just for checking from which server the request is being handled.
		// this is not recommended from security perspective as we don't want to let the client know which server is handling the request.
		w.Header().Add("X-Forwarded-Server", server.URL.String())
		server.ReverseProxy().ServeHTTP(w, r)
	})

	log.Println("Starting load balancer on port", config.Port)
	err = http.ListenAndServe(config.Port, nil)
	if err != nil {
		log.Fatalf("Error starting load balancer: %s\n", err.Error())
	}
}
