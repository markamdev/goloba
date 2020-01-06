package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/markamdev/goloba/pkg/balancer"
)

const (
	defLogFile  = "goloba.log"
	defConfFile = "goloba.conf"
)

var (
	logger       *log.Logger
	activityFlag bool
)

type context struct {
	locker   *sync.WaitGroup
	cfg      config
	balance  *balancer.Balancer
	listener net.Listener
}

func main() {
	fmt.Println("GoLoBa - simple Go Load Balancer (for TCP traffic)")

	// separate option for "help" flag
	var help bool
	var confFile string
	var outputFile string
	// read options
	flag.StringVar(&confFile, "f", defConfFile, "Path to configuration file")
	flag.BoolVar(&help, "h", false, "Print help message")
	flag.StringVar(&outputFile, "l", defLogFile, "Output file for application logs")
	flag.Parse()

	// if requested print help and exit
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// open (or create if not exists) output file for logs
	logFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		pe, ok := err.(*os.PathError)
		if ok {
			fmt.Println("Failed to init logger output file: ", pe.Unwrap().Error(), " -  Exiting ...")
		} else {
			fmt.Println("Failed to init logger output file: ", err.Error(), " - Exiting")
		}
		os.Exit(1)
	}
	defer logFile.Close()

	// prepare logger instance
	logger = log.New(logFile, "[GoLoBa] ", log.LstdFlags)
	logger.Println("Starting GoLoBa ...")

	// load config from file
	cfg, err := loadConfigFile(confFile)
	if err != nil {
		reportFailure(err.Error())
	}

	// prepare balancer
	blnc := balancer.New()
	blnc.Init(balancer.Configuration{Servers: cfg.Servers, Logger: logger})

	// prepare main application context
	glbCtx := context{locker: new(sync.WaitGroup),
		cfg:     cfg,
		balance: blnc,
	}

	// set global activity flag to true
	activityFlag = true

	// increment locker flag and launch listener
	glbCtx.locker.Add(1)
	go startConnectionListener(&glbCtx)

	// launch signal listener without waiting group incrementation
	go startSignalListener(&glbCtx)

	// wait till all child finished
	glbCtx.locker.Wait()

	fmt.Println("Connection forwarding finished")
}

func reportFailure(msg string) {
	fmt.Println("Fatal error occured - check logs for details")
	logger.Fatalln(msg)
}

func startSignalListener(ctx *context) {
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt)

	// just wait for signal - no need to save it
	_ = <-sch
	logger.Println("Interrupt signal received - preparing to exit")
	activityFlag = false
	if ctx.listener != nil {
		ctx.listener.Close()
	}
}
