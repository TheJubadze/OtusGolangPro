package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	host := flag.String("host", "localhost", "host to connect to")
	port := flag.String("port", "4242", "port to connect to")
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	address := net.JoinHostPort(*host, *port)
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	handleSignals()

	if err := client.Connect(); err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Error connecting:", err)
		return
	}
	defer func(client TelnetClient) {
		err := client.Close()
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "Error closing connection:", err)
		}
	}(client)

	_, _ = fmt.Fprintf(os.Stderr, "Connected to %s:%s", *host, *port)

	go func() {
		if err := client.Receive(); err != nil {
			_, _ = fmt.Fprint(os.Stderr, "Error receiving:", err)
		}
	}()

	if err := client.Send(); err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Error sending:", err)
	}
}

func handleSignals() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChannel
		_, _ = fmt.Fprintf(os.Stderr, "Received signal: %s\n", sig)
		os.Exit(0)
	}()
}
