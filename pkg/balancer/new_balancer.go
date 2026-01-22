package balancer

import (
	"github.com/markamdev/goloba/pkg/logger"
	"github.com/pkg/errors"
)

func NewLoadBalancer(config Config) LoadBalancer {
	return &balancerImpl{
		cfg:    config,
		logger: logger.GetDefaultLogger().WithParam("component", "balancer"),
	}
}

func (b *balancerImpl) Start() error {
	if err := b.validateConfig(); err != nil {
		return errors.Wrap(err, "failed to start load balancer due to invalid configuration")
	}
	b.logger.Debug("starting load balancer", "port", b.cfg.Port, "algorithm", b.cfg.Algorithm, "targets", b.cfg.Servers)
	return b.startListener()
}

func (b *balancerImpl) StartWithOptions(opts StartOptions) error {
	// TODO use opts to override b.cfg parameters
	return ErrNotImplemented
}

func (b *balancerImpl) Stop() error {
	b.logger.Debug("stopping load balancer")
	return b.stopListener()
}

func (b *balancerImpl) validateConfig() error {
	if b.cfg.Port == 0 {
		return errors.Wrap(ErrInvalidParam, "port must be greater than 0")
	}
	if len(b.cfg.Servers) == 0 {
		return errors.Wrap(ErrInvalidParam, "at least one target server must be specified")
	}
	if b.cfg.Algorithm == "" {
		return errors.Wrap(ErrInvalidParam, "load balancing algorithm must be specified")
	}
	if b.cfg.LogLevel == "" {
		return errors.Wrap(ErrInvalidParam, "log level must be specified")
	}
	return nil
}
