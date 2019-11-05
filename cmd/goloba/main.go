package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/markamdev/goloba/pkg/balancer"
)

var (
	logFile string = "goloba.log"
	logger  *log.Logger
)

type context struct {
	locker  *sync.WaitGroup
	cfg     config
	active  bool
	balance *balancer.Balancer
}

func main() {
	fmt.Println("GoLoBa - simple Go Load Balancer (for TCP traffic)")

	// separate option for "help" flag
	var help bool
	var confFile string
	// read options
	flag.StringVar(&confFile, "f", "goloba.conf", "Path to configuration file")
	flag.BoolVar(&help, "h", false, "Print help message")
	flag.StringVar(&logFile, "l", "goloba.log", "Output file for application logs")
	flag.Parse()

	// open (or create if not exists) output file for logs
	logFile, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	// if requested print help and exit
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

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
		active:  true,
		balance: blnc,
	}

	// increment locker flag and launch listener
	glbCtx.locker.Add(1)
	go startListener(&glbCtx)

	// wait till all child finished
	glbCtx.locker.Wait()
}

func reportFailure(msg string) {
	fmt.Println("Fatal error occured - check logfile: ", logFile)
	logger.Fatalln(msg)
}
