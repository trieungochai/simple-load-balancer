package main

import (
	"net/url"
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
