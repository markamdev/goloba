package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/markamdev/goloba/pkg/balancer"
)

const (
	defLogFile  = "goloba.log"
	defConfFile = "goloba.conf"
)

var blnc *balancer.Balancer

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

	// configure logging utility
	log.SetOutput(logFile)
	log.SetPrefix("[GoLoBa] ")
	log.SetFlags(log.LstdFlags)
	log.Println("Starting GoLoBa ...")

	// load config from file
	cfg, err := loadConfig(confFile)
	if err != nil {
		log.Fatalln("Failed to load config: ", err.Error())
	}

	// prepare balancer
	blnc = balancer.New()
	err = blnc.Init(balancer.Configuration{Port: cfg.Port, Servers: cfg.Servers})
	if err != nil {
		log.Fatalln("Failed to init balancer: ", err.Error())
	}

	// start balancer
	err = blnc.Start()
	if err != nil {
		log.Fatalln("Failed to start balancer: ", err.Error())
	}

	// launch signal listener without waiting group incrementation
	go startSignalListener()

	// wait till balancer finish working
	blnc.Wait()

	log.Println("Closing GoLoBa")
}

func reportFailure(msg string) {
	log.Fatalln(msg)
}

func startSignalListener() {
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt)

	// just wait for signal - no need to save it
	_ = <-sch
	log.Println("Interrupt signal received - preparing to exit")
	blnc.Stop()
}
