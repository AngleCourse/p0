package main

import (
	"bufio"
	"fmt"
	"github.com/cmu440/p0"
	"os"
)

const defaultPort = 20159

func main() {
	// Initialize the server.
	server := p0.New()
	if server == nil {
		fmt.Println("New() returned a nil server. Exiting...")
		return
	}

	// Start the server and continue listening for client connections in the background.
	if err := server.Start(defaultPort); err != nil {
		fmt.Printf("MultiEchoServer could not be started: %s\n", err)
		return
	}

	fmt.Printf("Started MultiEchoServer on port %d...\n", defaultPort)

	// Block forever.
	go func() {
		reader := bufio.NewReader(os.Stdin)
		text, _, _ := reader.ReadLine()
		if string(text) == "quit" {
			server.Close()
			os.Exit(0)
		}
	}()
	select {}
}
