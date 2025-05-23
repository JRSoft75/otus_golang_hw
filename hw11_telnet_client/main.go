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

	fmt.Printf("Connected to %s\n", address)

	done := make(chan struct{})

	// Горутина для чтения данных из сокета и вывода в STDOUT
	go func() {
		err := client.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("\nServer closed the connection.")
			} else {
				fmt.Printf("\nReceive error: %v\n", err)
			}
		}
		close(done)
	}()

	// Горутина для записи данных из STDIN в сокет
	go func() {
		err := client.Send()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("\nEOF detected. Closing connection.")
			} else {
				fmt.Printf("\nSend error: %v\n", err)
			}
		}
		close(done)
	}()

	// Ожидание завершения работы программы
	select {
	case <-done:
		fmt.Println("Connection closed.")
	case <-ctx.Done():
		fmt.Println("\nInterrupt signal received. Exiting...")
	}
}
