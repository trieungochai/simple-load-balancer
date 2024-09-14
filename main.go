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
type Servcer struct {
	URL       *url.URL
	IsHealthy bool
	Mutex     sync.Mutex
}

func main() {

}
