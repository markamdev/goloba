package balancer

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	forwaderBufferSize = 2048
)

func (b *Balancer) startListener() {
	defer b.locker.Done()

	// start launching connection listener
	log.Print("Starting connection listener on port: ", b.port)
	for _, serv := range b.servers {
		log.Println("... adding server: ", serv)
	}

	listenAddress := ":" + strconv.Itoa(int(b.port))

	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalln("Failed to open listening socket: ", err.Error())
	}
	// save listener to context
	b.listener = ln
	acceptChan := make(chan net.Conn, 1)
	cancelChan := b.ctx.Done()

	// TODO consider incrementing WaitGroup here
	go func(chn chan net.Conn, lstn net.Listener) {
		for {
			cn, err := lstn.Accept()
			if err == nil {
				chn <- cn
			} else {
				log.Println("Connection accepting error")
				return
			}
		}
	}(acceptChan, b.listener)

	for {
		select {
		case conn := <-acceptChan:
			log.Println("Connection accepted")
			b.handleConnection(conn)
		case <-cancelChan:
			log.Println("Cancellation received")
			b.listener.Close()
			return
		}
	}
}

// internal only functions below - no input param verification
func (b *Balancer) handleConnection(incoming net.Conn) {
	// get IP address only (WARNING: it's not working properly for IPv6 connections, even on localhost)
	source := strings.Split(incoming.RemoteAddr().String(), ":")[0]
	log.Println("Accepted connection from: ", source)

	// fetch destination (redirection) address
	dest, err := b.GetDestinationForSource(source)
	if err == nil {
		log.Println("... redirecting to: ", dest)
	} else {
		log.Println("... failed to get redirect address - exiting with error:", err)
		incoming.Close()
		return
	}

	// open new connection for data forwarding
	tempdialer := net.Dialer{}
	redirect, err := tempdialer.Dial("tcp", dest)
	if err != nil {
		log.Println("Failed to connect to redirect address (", dest, ") - exiting with error:", err)
		incoming.Close()
		return
	}

	// two new goroutines will be created
	b.locker.Add(2)

	// client to server forwarder
	go b.forwardRoutine(b.ctx, incoming, redirect)
	// server to client forwarder
	go b.forwardRoutine(b.ctx, redirect, incoming)
}

func (b *Balancer) forwardRoutine(ctx context.Context, in, out net.Conn) {

	// tracking of opened connections
	b.notifyOpened(in.RemoteAddr().String(), out.RemoteAddr().String())
	defer b.notifyClosed(in.RemoteAddr().String(), out.RemoteAddr().String())
	// tracking of working goroutines
	defer b.locker.Done()

	buffer := make([]byte, forwaderBufferSize)
	// repeat as long as no error from context (so context not closed/canceled)
	for ctx.Err() == nil {
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
	log.Println("Forwarding routine closed by context")
	in.Close()
}

// getDestination returns destination server (address:port) selected by round robin algorithm
func (b *Balancer) getDestination() (string, error) {
	if !b.initialized {
		return "", errors.New("Balancer not initializer")
	}

	// in future there will be another way of destination selection
	b.currentID++
	return b.servers[b.currentID%len(b.servers)], nil
}

// notifyOpened informs balanacing algorithm, that connection to particular destination has been opened
func (b *Balancer) notifyOpened(source, destination string) error {
	if !b.initialized {
		return errors.New("Balancer not initialized")
	}

	// lock access to mapping for reading
	notifLocker.Lock()
	defer notifLocker.Unlock()

	// if destination already saved - increment counter
	if dest, ok := b.mapping[source]; ok {
		if dest.destination != destination {
			// already saved destination is other than given one
			return errors.New("Notifying different destination than already saved")
		}
		dest.counter++
		b.mapping[source] = dest
		return nil
	}

	// save new mapping information
	b.mapping[source] = destCounter{counter: 1, destination: destination}
	return nil
}

// notifyClosed informs balancing algorithm, that connection to particular destination has been closed
func (b *Balancer) notifyClosed(source, destination string) error {
	if !b.initialized {
		return errors.New("Balancer not initialized")
	}

	// lock access to mapping for reading
	notifLocker.Lock()
	defer notifLocker.Unlock()

	dest, ok := b.mapping[source]
	// if no mapping found - some error occured
	if !ok {
		return errors.New("Closing not mapped connection")
	}
	if dest.destination != destination {
		// already saved destination is other than given one
		return errors.New("Closing connection with destination different than already saved")
	}

	dest.counter--
	if dest.counter == 0 {
		delete(b.mapping, source)
	} else {
		b.mapping[source] = dest
	}
	return nil
}
