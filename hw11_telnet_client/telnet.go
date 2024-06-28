package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *telnetClient) Send() error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		_, err := c.conn.Write([]byte(scanner.Text() + "\n"))
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stderr, "\nEOF received, exiting...")
	return nil
}

func (c *telnetClient) Receive() error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	_, err := io.Copy(c.out, c.conn)
	return err
}
