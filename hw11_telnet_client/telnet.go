package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

type TCPClient struct {
	address string
	timeout time.Duration

	in  io.ReadCloser
	out io.Writer

	conn net.Conn
}

func (tc *TCPClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return fmt.Errorf("error on connect: %w", err)
	}
	tc.conn = conn
	return nil
}

func (tc *TCPClient) Close() error {
	log.Println("CLosed")
	err := tc.conn.Close()
	if err != nil {
		return fmt.Errorf("error on closing client: %w", err)
	}
	return nil
}

func (tc *TCPClient) Send() error {
	_, err := io.Copy(tc.conn, tc.in)
	if err != nil {
		return fmt.Errorf("error while sending: %w", err)
	}
	return nil
}

func (tc *TCPClient) Receive() error {
	_, err := io.Copy(tc.out, tc.conn)
	if err != nil {
		return fmt.Errorf("error while receiving: %w", err)
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TCPClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
