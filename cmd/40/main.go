package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Not-cottage-cheese-but-cottage-cheese/final-go/server"
)

func main() {
	server := server.NewServer()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-shutdown
		fmt.Println("Gracefully shutting down...")
		server.Shutdown()
	}()

	if err := server.Run(); err != nil {
		panic(err)
	}
}
