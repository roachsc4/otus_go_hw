package main

import (
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

	in  io.Reader
	out io.Writer

	conn net.Conn
}

func (tc *TCPClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	tc.conn = conn
	return nil
}

func (tc *TCPClient) Close() error {
	log.Println("CLosed")
	err := tc.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (tc *TCPClient) Send() error {
	_, err := io.Copy(tc.conn, tc.in)
	return err
}

func (tc *TCPClient) Receive() error {
	_, err := io.Copy(tc.out, tc.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TCPClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
