package main

import (
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	forwaderBufferSize = 2048
)

func startListener(ctx *context) {
	defer ctx.locker.Done()

	// check input params first
	if ctx == nil {
		reportFailure("Internal error - invalid context")
	}

	if len(ctx.cfg.Servers) == 0 {
		reportFailure("No servers added - exiting!")
	}

	if ctx.cfg.Port == 0 {
		reportFailure("Invalid port number")
	}

	// start launching connection listener
	logger.Print("Starting listener on port: ", ctx.cfg.Port)
	for _, serv := range ctx.cfg.Servers {
		logger.Print("... adding server ", serv)
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

	// TODO add waiting for all forwarders to finish (additional wait group) or just close connections
}

// internal only functions below - no input param verification
func handleConnection(ctx *context, incoming net.Conn) {
	defer incoming.Close()

	// get IP address only (WARNING: it's not working properly for IPv6 connections, even on localhost)
	source := strings.Split(incoming.RemoteAddr().String(), ":")[0]
	logger.Println("Accepted connection from: ", source)

	// fetch destination (redirection) address
	dest, err := ctx.balance.GetDestinationForSource(source)
	if err == nil {
		logger.Println("... redirecting to: ", dest)
	} else {
		logger.Println("... failed to get redirect address - exiting")
		return
	}

	// open new connection for data forwarding
	redirect, err := net.Dial("tcp", dest)
	if err != nil {
		logger.Println("Failed to connect to redirect address - exiting")
		return
	}
	ctx.balance.NotifyOpened(source, dest)
	defer ctx.balance.NotifyClosed(source, dest)

	var locker sync.WaitGroup
	locker.Add(2)

	// client to server forwarder
	go forwardRoutine(incoming, redirect, &locker)
	// server to client forwarder
	go forwardRoutine(redirect, incoming, &locker)

	// wait for subroutines to finish
	locker.Wait()
}

func forwardRoutine(in, out net.Conn, lock *sync.WaitGroup) {
	defer lock.Done()
	buffer := make([]byte, forwaderBufferSize)
	for {
		n, err := in.Read(buffer)
		if err != nil {
			out.Close()
			return
		}
		if n == 0 {
			continue
		}
		out.Write(buffer[:n])
	}

}
