package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var timeout time.Duration

func main() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
	flag.Parse()

	args := flag.Args()

	if len(args) < 2 {
		flag.Usage()
		log.Fatalln("Используйте: go-telnet [--timeout=10s] <host> <port>")
	}

	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalln(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go receive(wg, &client)
	go send(wg, &client)

	wg.Wait()
}

func send(wg *sync.WaitGroup, client *TelnetClient) {
	func() {
		defer wg.Done()
		if err := (*client).Send(); err != nil {
			log.Fatalf("Ошибка получения данных %v", err)
		}
		fmt.Fprintln(os.Stderr, "Соединение было закрыто удаленным хостом...")
	}()
}

func receive(wg *sync.WaitGroup, client *TelnetClient) {
	func() {
		defer wg.Done()
		if err := (*client).Receive(); err != nil {
			log.Fatalf("Ошибка чтения %v:", err)
		}
		fmt.Fprintln(os.Stderr, "Соединение закрыто")
	}()
}
