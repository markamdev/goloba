package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"

	"github.com/markamdev/goloba/pkg/balancer"
	"github.com/markamdev/goloba/pkg/utils"
)

var (
	loggingFile = flag.String("log-file", "", "Output file for logs")
	help        = flag.Bool("h", false, "Print help screen")
	port        = flag.Int("port", 8060, "GoLoBa listening port")
	servers     = flag.String("targets", "", "List of comma separated target servers")
)

var blnc *balancer.Balancer

func main() {
	utils.SetupLogger()
	logrus.Debugln("GoLoBa - simple Go Load Balancer (for TCP traffic) v.", currentVersion)

	flag.Parse()

	if loggingFile == nil {
		loggingFile = new(string)
	}
	if len(*loggingFile) == 0 {
		*loggingFile = stdOutPath
	}

	// if requested print help and exit
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// open (or create if not exists) output file for logs
	logFile, err := os.OpenFile(*loggingFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	// prepare balancer
	blnc = balancer.New()
	err = blnc.Init(balancer.Configuration{
		Port:    uint(*port),
		Servers: strings.Split(strings.Trim(*servers, "\""), ","),
	})
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

	logrus.Debug("Closing GoLoBa")
}

func fatalAtStart(msg string, er error) {
	// leave message about fatal error on console
	fmt.Println("Fatal error occured - see logs for details")
	// log error and exit
	if er != nil {
		logrus.Fatalln(msg, er.Error())
	} else {
		logrus.Fatalln(msg)
	}
}

func startSignalListener() {
	logrus.Debug("Starting signal listener")
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt)

	// just wait for signal - no need to save it
	<-sch
	logrus.Debugln("Interrupt signal received - preparing to exit")
	blnc.Stop()
}
