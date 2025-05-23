package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Println("Usage: go-telnet [--timeout=duration] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}
	defer func(client TelnetClient) {
		_ = client.Close()
	}(client)

	fmt.Fprintf(os.Stderr, "Connected to %s\n", address)

	sendDone := make(chan struct{})
	receiveDone := make(chan struct{})

	// Горутина для чтения данных из сокета и вывода в STDOUT
	go func() {
		err := client.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stderr, "\nServer closed the connection.")
			} else {
				fmt.Fprintf(os.Stderr, "\nReceive error: %v\n", err)
			}
			return
		}
		close(receiveDone)
	}()

	// Горутина для записи данных из STDIN в сокет
	go func() {
		err := client.Send()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stderr, "\nEOF detected. Closing connection.")
			} else {
				fmt.Fprintf(os.Stderr, "\nSend error: %v\n", err)
			}
			return
		}
		close(sendDone)
	}()

	// Ожидание завершения работы программы
	select {
	case <-sendDone:
		// Ждем завершения receive
		<-receiveDone
		fmt.Fprintln(os.Stderr, "Send and receive completed.")
	case <-receiveDone:
		// Ждем завершения send
		<-sendDone
		_, err := fmt.Fprintln(os.Stderr, "Receive and send completed.")
		if err != nil {
			return
		}
	case <-ctx.Done():
		_, err := fmt.Fprintln(os.Stderr, "\nInterrupt signal received. Exiting...")
		if err != nil {
			return
		}
	}
}
