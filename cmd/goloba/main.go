package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/markamdev/goloba/internal/appconfig"
	"github.com/markamdev/goloba/pkg/logger"
	"github.com/spf13/pflag"
)

func main() {
	// load application configuration
	cfg, err := appconfig.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %s", err.Error()))
	}

	if cfg.Help {
		pflag.PrintDefaults()
		return
	}

	// TODO initialize logger here

	if cfg.LogLevel == "debug" {
		fmt.Printf("application config: %+v\n", cfg)
	}

	// launch signal listener without waiting group incrementation
	go startSignalListener()
}

func startSignalListener() {
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt)

	// just wait for signal - no need to save it
	<-sch
	logger.Debug("interrupt signal received - preparing to exit")
	// TODO: add graceful shutdown logic here
}
