package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sync"
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

}
