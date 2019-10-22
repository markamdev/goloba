package main

import (
	"net"
	"strconv"
)

func startListener(ctx *context) {
	defer ctx.locker.Done()
	logger.Print("Starting listener on port: ", ctx.cfg.Port)
	for _, serv := range ctx.cfg.Servers {
		logger.Print("... adding server ", serv)
	}

	if len(ctx.cfg.Servers) == 0 {
		reportFailure("No servers added - exiting!")
	}

	if ctx.cfg.Port == 0 {
		reportFailure("Invalid port number")
	}
	listenAddress := ":" + strconv.Itoa(int(ctx.cfg.Port))

	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		reportFailure(err.Error())
	}

	for ctx.active {
		conn, err := ln.Accept()
		if err != nil {
			logger.Println("Failed to accept incoming connection: ", err)
			continue
		}

		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx *context, conn net.Conn) {
	defer conn.Close()
	// debug code
	logger.Println("Accepted connection: ", conn.RemoteAddr().String())
	srv, err := ctx.balance.GetServer()
	if err == nil {
		logger.Println("... redirecting to: ", srv)
	} else {
		logger.Println(".. FAILED TO REDIRECT")
	}
	// ---------
}
