package balancer

type balancerImpl struct {
	cfg Config
}

func NewLoadBalancer(config Config) LoadBalancer {
	return &balancerImpl{cfg: config}
}

func (b *balancerImpl) Start() error {
	return nil
}

func (b *balancerImpl) StartWithOptions(opts StartOptions) error {
	return b.Start()
}
