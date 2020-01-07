package balancer

import (
	"errors"
	"net"
	"sync"
)

const (
	defServerCount = 5
)

var (
	notifLocker sync.RWMutex
)

// Configuration stores configuration data for load balancer module
type Configuration struct {
	Port    uint
	Servers []string
}

type destCounter struct {
	destination string
	counter     int
}

// Balancer is a load balancing engine object
type Balancer struct {
	// TODO Change this struct int interface + struct with functions
	initialized bool
	active      bool
	port        uint
	servers     []string
	currentID   int
	mapping     map[string]destCounter
	locker      sync.WaitGroup
	listener    net.Listener
}

// New creates new Balancer instance
func New() *Balancer {
	var b Balancer
	b.initialized = false
	b.active = false
	b.currentID = 0
	b.servers = make([]string, defServerCount)
	b.mapping = make(map[string]destCounter)
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

	if cfg.Port == 0 {
		return errors.New("Invalid listening port number")
	}

	b.servers = cfg.Servers
	b.port = cfg.Port
	b.initialized = true
	return nil
}

// Start load balancer functionality
func (b *Balancer) Start() error {
	if !b.initialized {
		return errors.New("Balancer not initialized - failed to start")
	}

	b.active = true
	b.startListener()

	return nil
}

// Stop informs balancer to close all connections
// use #Wait() to be sure that balancer finished working
func (b *Balancer) Stop() error {
	if !b.initialized {
		return errors.New("Balancer not initialized - cannot stop")
	}
	b.active = false
	b.listener.Close()

	return nil
}

// Wait stops execution till balancer completely deinitialized
func (b *Balancer) Wait() {
	b.locker.Wait()
}

// GetDestinationForSource returns destination server (address:port) for given source address
// This function is selecting same destination server for same client address (excluding port number)
func (b *Balancer) GetDestinationForSource(source string) (string, error) {
	if !b.initialized {
		return "", errors.New("Balancer not initializer")
	}

	// lock access to mapping for reading
	notifLocker.RLock()
	defer notifLocker.RUnlock()

	// if destination already found - return it
	if dest, ok := b.mapping[source]; ok {
		return dest.destination, nil
	}

	// temporary use 'round robin' algorithm for destination selection
	dest, err := b.getDestination()
	if err != nil {
		return "", errors.New("Internal Balancer error")
	}

	return dest, nil
}
