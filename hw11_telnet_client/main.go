package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?

	host := flag.String("host", "localhost", "host to connect to")
	port := flag.String("port", "4242", "port to connect to")
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	address := *host + ":" + *port

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer func(client TelnetClient) {
		err := client.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}(client)

	go func() {
		if err := client.Receive(); err != nil {
			fmt.Println("Error receiving:", err)
		}
	}()

	if err := client.Send(); err != nil {
		fmt.Println("Error sending:", err)
	}
}
