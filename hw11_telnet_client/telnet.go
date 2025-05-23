package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
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
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	dialer := net.Dialer{}
	t.conn, err = dialer.DialContext(ctx, "tcp", t.address)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

// Send отправляет данные из STDIN в сокет.
func (t *telnetClient) Send() error {
	reader := bufio.NewReader(t.in)
	writer := bufio.NewWriter(t.conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("send error: %w", err)
		}

		_, err = writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("send error: %w", err)
		}
		err = writer.Flush()
		if err != nil {
			return fmt.Errorf("flush error: %w", err)
		}
	}
}

// Receive получает данные из сокета и выводит их в STDOUT.
func (t *telnetClient) Receive() error {
	reader := bufio.NewReader(t.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("receive error: %w", err)
		}
		_, err = t.out.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("write to output error: %w", err)
		}
	}
}
