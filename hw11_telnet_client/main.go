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
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Correct usage: telnet host port [--timeout=2s]\n")
	}
	flag.Parse()
	flag.Usage()
	positionalArgs := flag.Args()

	host := positionalArgs[0]
	port := positionalArgs[1]
	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(
		address,
		*timeout,
		os.Stdin,
		os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()
	log.Printf("Connected to %s...", address)
	sigintChannel := make(chan os.Signal, 1)
	doneCh := make(chan int)

	signal.Notify(sigintChannel, syscall.SIGINT)

	go func() {
		<-sigintChannel
		_, err := os.Stderr.WriteString("Got SIGINT")
		if err != nil {
			log.Fatal(err)
		}
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
