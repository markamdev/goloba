package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/markamdev/goloba/internal/appconfig"
	"github.com/markamdev/goloba/pkg/balancer"
	"github.com/markamdev/goloba/pkg/logger"
	"github.com/spf13/pflag"
)

func main() {
	// load application configuration
	cfg, err := appconfig.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %s\n", err.Error())
		os.Exit(1)
	}

	if cfg.Help {
		pflag.PrintDefaults()
		return
	}

	logger.SetDefaultLogger(logger.NewBasicLogger())
	logger.SetLevel(logger.ParseLogLevel(cfg.LogLevel))

	logger.Debug("creating new load balancer", "config", cfg)
	// initialize/create load balancer
	bl := balancer.NewLoadBalancer(appconfig.AppConfigToBalancerConfig(cfg))

	// launch signal listener without waiting group incrementation
	go startSignalListener(func() {
		err := bl.Stop()
		if err != nil {
			logger.Error("failed to stop load balancer", "error", err.Error())
		}
	})

	err = bl.Start()
	if err != nil {
		logger.Fatal("failed to start load balancer", "error", err.Error())
	}
}

func startSignalListener(stopRoutine func()) {
	logger.Debug("starting signal listener for interrupt signals")
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt)

	// just wait for signal - no need to save it
	<-sch
	logger.Debug("interrupt signal received - preparing to exit")

	stopRoutine()
}
