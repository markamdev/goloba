package balancer

import (
	"errors"
	"log"
)

const (
	defServerCount = 5
)

// Configuration stores configuration data for load balancer module
type Configuration struct {
	Servers []string
	Logger  *log.Logger
}

// Balancer is a load balancing engine object
type Balancer struct {
	initialized bool
	servers     []string
	currentID   int
}

// New creates new Balancer instance
func New() *Balancer {
	var b Balancer
	b.initialized = false
	b.currentID = 0
	b.servers = make([]string, defServerCount)
	return &b
}

// Init initializes balancer with necessary data (has to be called before any other call)
func (b *Balancer) Init(cfg Configuration) error {
	if b.initialized {
		// already initialized - it's not an error
		return nil
	}

	if len(cfg.Servers) == 0 {
		return errors.New("At least one server needed")
	}

	b.servers = cfg.Servers
	b.initialized = true
	return nil
}

// GetServer returns server (address:port) selected by balancing algorithm
func (b *Balancer) GetServer() (string, error) {
	if !b.initialized {
		return "", errors.New("Balancer not initializer")
	}

	b.currentID++
	return b.servers[b.currentID%len(b.servers)], nil
}
