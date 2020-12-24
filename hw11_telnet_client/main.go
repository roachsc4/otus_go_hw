package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	timeout := flag.Duration("timeout", 10000000, "use it to specify dial timeout")
	flag.Parse()
	positinalArgs := flag.Args()

	if len(positinalArgs) != 2 {
		log.Fatal("Correct usage: telnet host port [--timeout=2s]")
	}

	host := positinalArgs[0]
	port := positinalArgs[1]
	address := net.JoinHostPort(host, port)
	log.Println(address)
	client := NewTelnetClient(
		address,
		*timeout,
		os.Stdin,
		os.Stdout)

	err := client.Connect()
	defer client.Close()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Connected to %s...", address)
	sigintChannel := make(chan os.Signal, 1)
	doneCh := make(chan int)

	signal.Notify(sigintChannel, syscall.SIGINT)

	go func() {
		<-sigintChannel
		fmt.Println("Got SIGINT")
		doneCh <- 3
	}()

	go func() {
		log.Println("Start receiving")
		err := client.Receive()
		if err != nil {
			log.Println("Error during receive: ", err)
		}
		log.Println("Stop receiving")
		doneCh <- 1
	}()

	go func() {
		log.Println("Start sending")
		err := client.Send()
		if err != nil {
			log.Println("Error during send: ", err)
		}
		log.Println("Stop sending")
		doneCh <- 2
	}()

	log.Println(<-doneCh)
}
