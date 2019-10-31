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

type destCounter struct {
	destination string
	counter     int
}

// Balancer is a load balancing engine object
type Balancer struct {
	initialized bool
	servers     []string
	currentID   int
	mapping     map[string]destCounter
}

// New creates new Balancer instance
func New() *Balancer {
	var b Balancer
	b.initialized = false
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

	b.servers = cfg.Servers
	b.initialized = true
	return nil
}

// GetDestination returns destination server (address:port) selected by round robin algorithm
// This function can be used when connection source (client) address is ignored when redirecting traffic
func (b *Balancer) GetDestination() (string, error) {
	if !b.initialized {
		return "", errors.New("Balancer not initializer")
	}

	// in future there will be another way of destination selection
	b.currentID++
	return b.servers[b.currentID%len(b.servers)], nil
}

// GetDestinationForSource returns destination server (address:port) for given source address
// This function is selecting same destination server for same client address (excluding port number)
func (b *Balancer) GetDestinationForSource(source string) (string, error) {
	if !b.initialized {
		return "", errors.New("Balancer not initializer")
	}

	// if destination already found - return it
	if dest, ok := b.mapping[source]; ok {
		return dest.destination, nil
	}

	// temporary use 'round robin' algorithm for destination selection
	dest, err := b.GetDestination()
	if err != nil {
		return "", errors.New("Internal Balancer error")
	}

	return dest, nil
}

// NotifyOpened informs balanacing algorithm, that connection to particular destination has been opened
func (b *Balancer) NotifyOpened(source, destination string) error {
	if !b.initialized {
		return errors.New("Balancer not initialized")
	}

	// if destination already saved - increment counter
	if dest, ok := b.mapping[source]; ok {
		if dest.destination != destination {
			// already saved destination is other than given one
			return errors.New("Notifying different destination than already saved")
		}
		dest.counter++
		b.mapping[source] = dest
		return nil
	}

	// save new mapping information
	b.mapping[source] = destCounter{counter: 1, destination: destination}
	return nil
}

// NotifyClosed informs balancing algorithm, that connection to particular destination has been closed
func (b *Balancer) NotifyClosed(source, destination string) error {
	if !b.initialized {
		return errors.New("Balancer not initialized")
	}

	dest, ok := b.mapping[source]
	// if no mapping found - some error occured
	if !ok {
		return errors.New("Closing not mapped connection")
	}
	if dest.destination != destination {
		// already saved destination is other than given one
		return errors.New("Closing connection with destination different than already saved")
	}

	dest.counter--
	if dest.counter == 0 {
		delete(b.mapping, source)
	} else {
		b.mapping[source] = dest
	}
	return nil
}
