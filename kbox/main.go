package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func printUsage() {
	fmt.Println("KBox - Jupyter Kernel Container Management CLI") 
	fmt.Println("\nUsage:")
	fmt.Println("  kbox <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  connect  - Setup and start a kernel container (requires connection_file and tag)")
	fmt.Println("\nExample:")
	fmt.Println("  kbox connect ./kernel-connection.json my-image:latest")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	client, err := NewKBoxClient()
	if err != nil {
		fmt.Printf("Error connecting to Docker: %v\n", err)
		os.Exit(1)
	}

	if cmd == "connect" {
		if len(os.Args) < 4 {
			fmt.Println("Error: 'connect' requires connection_file and docker_tag")
			fmt.Println("Usage: kbox connect <connection_file> <docker_tag>")
			os.Exit(1)
		}
		connFile := os.Args[2]
		tag := os.Args[3]
		
		// Default docker run command from connect_kbox.sh
		runCmd := []string{"python", "-m", "ipykernel_launcher", "-f", "/connection_file.json"}
		
		containerID, err := client.ConnectKBox(connFile, tag, runCmd)
		if err != nil {
			fmt.Printf("Error during connection setup: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Connection established! Container %s is now running. Waiting for shutdown signal...\n", containerID)

		// Block and wait for termination signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		<-sigChan // Block until a signal is received

		fmt.Printf("\nShutdown signal received. Stopping container %s...\n", containerID)
		if err := client.StopContainer(containerID); err != nil {
			fmt.Printf("Error stopping container: %v\n", err)
		} else {
			fmt.Println("Container stopped successfully.")
		}
		return
	}
}
