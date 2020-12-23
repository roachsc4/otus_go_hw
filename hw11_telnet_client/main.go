package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Correct usage: telnet host port [--timeout=2s]")
	}
	timeout := flag.Duration("timeout", 10000000, "use it to specify dial timeout")
	flag.Parse()

	log.Println(os.Args[1:])
	host := os.Args[2]
	port := os.Args[3]
	address := net.JoinHostPort(host, port)
	log.Println(address)
	client := NewTelnetClient(
		address,
		*timeout,
		os.Stdin,
		os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to %s...", address)
	sigintChannel := make(chan os.Signal, 1)

	signal.Notify(sigintChannel, syscall.SIGINT)

	go func() {
		<-sigintChannel
		fmt.Println("Got SIGINT")
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			err := client.Receive()
			if err != nil {
				log.Println(err)
				break
			}

		}
		wg.Done()
	}()

	go func() {
		for {
			err := client.Send()
			if err != nil {
				log.Println(err)
				break
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
