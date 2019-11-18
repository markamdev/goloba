package main

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	forwaderBufferSize = 2048
)

func startConnectionListener(ctx *context) {
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
	// save listener to context
	ctx.listener = ln

	for activityFlag {
		conn, err := ln.Accept()
		if err != nil {
			if activityFlag == true {
				// print error message only if still active ...
				logger.Println("Failed to accept incoming connection: ", err)
			}
			// ... otherwise just exit
			break
		}

		// there should be additional waiting flag incr. for each connection
		ctx.locker.Add(1)
		go handleConnection(ctx, conn)
	}
}

// internal only functions below - no input param verification
func handleConnection(ctx *context, incoming net.Conn) {
	// remember to decrement main wait group at return
	defer ctx.locker.Done()

	// get IP address only (WARNING: it's not working properly for IPv6 connections, even on localhost)
	source := strings.Split(incoming.RemoteAddr().String(), ":")[0]
	logger.Println("Accepted connection from: ", source)

	// fetch destination (redirection) address
	dest, err := ctx.balance.GetDestinationForSource(source)
	if err == nil {
		logger.Println("... redirecting to: ", dest)
	} else {
		logger.Println("... failed to get redirect address - exiting")
		incoming.Close()
		return
	}

	// open new connection for data forwarding
	redirect, err := net.Dial("tcp", dest)
	if err != nil {
		logger.Println("Failed to connect to redirect address - exiting")
		incoming.Close()
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

	// wait till both forwarding routines finish working
	locker.Wait()
}

func forwardRoutine(in, out net.Conn, lck *sync.WaitGroup) {
	// remember to inform connection handler about finished forwarding
	defer lck.Done()

	buffer := make([]byte, forwaderBufferSize)
	for activityFlag {
		in.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
		n, err := in.Read(buffer)

		if err != nil {
			neterr, ok := err.(net.Error)
			// reading error occured - close related output socket to completely break connection
			if ok && neterr.Timeout() {
				// ignore network error
				continue
			}
			out.Close()
			return
		}
		if n == 0 {
			continue
		}
		out.Write(buffer[:n])
	}
	in.Close()
}
