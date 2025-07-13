package main

import (
	"os"
	"os/signal"
	"strings"

	"github.com/namsral/flag"

	"github.com/markamdev/goloba/internal/version"
	"github.com/markamdev/goloba/pkg/balancer"
	"github.com/markamdev/goloba/pkg/logger"
)

var (
	// logLevel    = flag.String("log-level", "info", "Log level: debug, info, warn, error, fatal")
	// loggingFile = flag.String("log-file", "", "Output file for logs")
	help    = flag.Bool("h", false, "Print help screen")
	port    = flag.Int("port", 8060, "GoLoBa listening port")
	servers = flag.String("targets", "", "List of comma separated target servers")
	// configFile  = flag.String("config", "", "Configuration file")
)

var blnc *balancer.Balancer

func main() {
	logger.SetLevel(logger.GlbDebug)
	logger.Info("goloba - simple TCP load balancer", "version", version.Version)

	flag.Parse()

	// if requested print help and exit
	if *help {
		// TODO add help message here
		flag.PrintDefaults()
		os.Exit(0)
	}

	// prepare balancer
	blnc = balancer.New()
	err := blnc.Init(balancer.Configuration{
		Port:    uint(*port),
		Servers: strings.Split(strings.Trim(*servers, "\""), ","),
	})
	if err != nil {
		logger.Fatal("failed to init balancer: ", "error", err.Error())
	}

	// start balancer
	err = blnc.Start()
	if err != nil {
		logger.Fatal("failed to start balancer: ", "error", err.Error())
	}

	// launch signal listener without waiting group incrementation
	go startSignalListener()

	// wait till balancer finish working
	blnc.Wait()

	logger.Debug("Closing GoLoBa")
}

func startSignalListener() {
	logger.Debug("starting signal listener")
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt)

	// just wait for signal - no need to save it
	<-sch
	logger.Debug("interrupt signal received - preparing to exit")
	blnc.Stop()
}
