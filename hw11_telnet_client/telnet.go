package main

import (
	"bufio"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	io.Closer
	Connect() error
	Receive() error
	Send() error
}

type Client struct {
	connection net.Conn
	address    string
	timeout    time.Duration
	in         io.Reader
	out        io.Writer
}

func (client *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", client.address, client.timeout)
	if err != nil {
		return err
	}
	client.connection = conn
	return nil
}

func (client *Client) Send() error {
	scanner := bufio.NewScanner(client.in)
	for scanner.Scan() {
		_, err := client.connection.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (client *Client) Receive() error {
	scanner := bufio.NewScanner(client.connection)
	for scanner.Scan() {
		_, err := client.out.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (client *Client) Close() error {
	if client.connection != nil {
		return client.connection.Close()
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
