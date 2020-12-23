package main

import (
	"bufio"
	"fmt"
	"io"
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

	inReader *bufio.Reader
	out      io.Writer

	conn       net.Conn
	connReader *bufio.Reader
}

func (tc *TCPClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	tc.conn = conn
	tc.connReader = bufio.NewReader(conn)
	return nil
}

func (tc *TCPClient) Close() error {
	err := tc.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (tc *TCPClient) Send() error {
	data, err := tc.inReader.ReadBytes('\n')
	if err != nil {
		return err
	}
	_, err = tc.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (tc *TCPClient) Receive() error {
	data, connErr := tc.connReader.ReadBytes('\n')
	_, outWriteErr := tc.out.Write(data)
	if connErr != nil {
		return fmt.Errorf("error during receiving data: %w", connErr)
	}
	if outWriteErr != nil {
		return fmt.Errorf("error during outwriting data: %w", outWriteErr)
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TCPClient{
		address:  address,
		timeout:  timeout,
		inReader: bufio.NewReader(in),
		out:      out,
	}
}
