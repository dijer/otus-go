package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout duration")
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatalln("host and port required!")
		return
	}

	host, port := flag.Arg(0), flag.Arg(1)

	if host == "" || port == "" {
		log.Fatal("host or port empty")
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	defer client.Close()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("...Connected to %s!\n", address)

	go func() {
		defer stop()

		if err := client.Send(); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}
	}()

	go func() {
		defer stop()

		if err := client.Receive(); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}
	}()

	<-ctx.Done()
}
