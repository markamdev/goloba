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
	stdOutPath  = "/dev/stdout"
)

var blnc *balancer.Balancer

func main() {
	fmt.Println("GoLoBa - simple Go Load Balancer (for TCP traffic)")

	// separate option for "help" flag
	var help bool
	var confFileName string
	var logFileName string
	var logStdOut bool
	// read options
	flag.StringVar(&confFileName, "f", defConfFile, "Path to configuration file")
	flag.BoolVar(&help, "h", false, "Print help message")
	flag.StringVar(&logFileName, "l", defLogFile, "Output file for application logs")
	flag.BoolVar(&logStdOut, "log-stdout", false, "Print logs on standard output instead of file")
	flag.Parse()

	if logStdOut {
		logFileName = stdOutPath
	}

	// if requested print help and exit
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// open (or create if not exists) output file for logs
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
	cfg, err := loadConfig(confFileName)
	if err != nil {
		fatalAtStart("Failed to load config: ", err)
	}

	// prepare balancer
	blnc = balancer.New()
	err = blnc.Init(balancer.Configuration{Port: cfg.Port, Servers: cfg.Servers})
	if err != nil {
		fatalAtStart("Failed to init balancer: ", err)
	}

	// start balancer
	err = blnc.Start()
	if err != nil {
		fatalAtStart("Failed to start balancer: ", err)
	}

	// launch signal listener without waiting group incrementation
	go startSignalListener()

	// wait till balancer finish working
	blnc.Wait()

	log.Println("Closing GoLoBa")
}

func fatalAtStart(msg string, er error) {
	// leave message about fatal error on console
	fmt.Println("Fatal error occured - see logs for details")
	// log error and exit
	if er != nil {
		log.Fatalln(msg, er.Error())
	} else {
		log.Fatalln(msg)
	}
}

func startSignalListener() {
	log.Println("Starting signal listener")
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt)

	// just wait for signal - no need to save it
	_ = <-sch
	log.Println("Interrupt signal received - preparing to exit")
	blnc.Stop()
}
