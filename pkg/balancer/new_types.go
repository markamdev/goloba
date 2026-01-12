package balancer

type LoadBalancer interface {
	Start() error
	StartWithOptions(opts StartOptions) error
}

type Config struct {
	Port      uint
	Servers   []string
	Algorithm string
	LogLevel  string
}

type StartOptions struct {
}
