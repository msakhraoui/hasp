package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go tailCowrie()
	go tailSnort()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	fmt.Println("shutting down:", sig)
}
