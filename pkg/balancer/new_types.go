package balancer

import "errors"

var (
	ErrNotImplemented   = errors.New("not implemented")
	ErrAlreadyRunning   = errors.New("already running")
	ErrInvalidParam     = errors.New("invalid parameter")
	ErrBrokenConnection = errors.New("broken connection")
)

type LoadBalancer interface {
	// Start launches the load balancer with default settings and waits for incoming requests.
	// Bahavior is similar to ListenAndServe in net/http package - it blocks the current goroutine.
	Start() error
	// StartWithOptions launches the load balancer with specified options.
	StartWithOptions(opts StartOptions) error
	// Stop gracefully stops the load balancer, closing all active connections and breaking the
	Stop() error
}

// Config holds the configuration parameters for the load balancer.
type Config struct {
	Port      uint
	Servers   []string
	Algorithm string
	LogLevel  string
}

// StartOptions holds options to be used when starting the load balancer instead of alrady configured ones.
type StartOptions struct {
}
