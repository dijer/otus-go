package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var ErrCantConnect = errors.New("cant connect")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	connection net.Conn
	address    string
	in         io.ReadCloser
	out        io.Writer
	timeout    time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *telnetClient) Connect() error {
	connection, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return ErrCantConnect
	}

	tc.connection = connection

	return nil
}

func (tc *telnetClient) Send() error {
	_, err := io.Copy(tc.connection, tc.in)

	return err
}

func (tc *telnetClient) Receive() error {
	_, err := io.Copy(tc.out, tc.connection)

	return err
}

func (tc *telnetClient) Close() error {
	return tc.connection.Close()
}
