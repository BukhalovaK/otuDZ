package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for connecting to the server")

	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Wrong arguments\n")
	}

	address := net.JoinHostPort(args[0], args[1])

	telnetClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := telnetClient.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer telnetClient.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		if err := telnetClient.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := telnetClient.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
	case <-time.After(timeout + 1*time.Second):
	}

	telnetClient.Close()

	wg.Wait()
}
