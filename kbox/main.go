package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func printUsage() {
	fmt.Println("KBox - Container Management CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  kbox <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  images   - List available images")
	fmt.Println("  connect  - Setup and start a kernel container (requires connection_file and tag)")
	fmt.Println("  status   - Check if the container is running")
	fmt.Println("  ensure   - Check if running, start if stopped (Auto-fix)")
	fmt.Println("  start    - Start the container")
	fmt.Println("  stop     - Stop the container")
	fmt.Println("  restart  - Restart the container")
	fmt.Println("  logs     - Get last logs of the container")
	fmt.Println("\nExample:")
	fmt.Println("  kbox images")
	fmt.Println("  kbox connect ./kernel-connection.json my-image:latest")
	fmt.Println("  kbox status my_nginx_container")
	fmt.Println("  kbox start my_nginx_container")
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

	if cmd == "images" {
		images, err := client.ListImages()
		if err != nil {
			fmt.Printf("Error listing images: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Available Images:")
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Printf("%-20s %-20s %-10s\n", "IMAGE ID", "REPOSITORY", "TAG")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, img := range images {
			repo := "<none>"
			tag := "<none>"
			if len(img.RepoTags) > 0 {
				repo = img.RepoTags[0]
			}
			fmt.Printf("%-20s %-20s %-10s\n", img.ID[:12], repo, tag)
		}
		return
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

	if len(os.Args) < 3 {
		fmt.Printf("Command %s requires a container ID\n", cmd)
		printUsage()
		os.Exit(1)
	}

	containerID := os.Args[2]

	switch cmd {

	case "status":
		running, err := client.IsRunning(containerID)
		if err != nil {
			fmt.Printf("Error checking status: %v\n", err)
			os.Exit(1)
		}
		if running {
			fmt.Printf("Container %s is RUNNING\n", containerID)
		} else {
			fmt.Printf("Container %s is STOPPED\n", containerID)
		}

	case "ensure":
		// This action ensures the container is running. 
		// If it's stopped, it starts it.
		running, err := client.IsRunning(containerID)
		if err != nil {
			fmt.Printf("Error checking status: %v\n", err)
			os.Exit(1)
		}
		if !running {
			fmt.Printf("Container %s is stopped. Starting it now...\n", containerID)
			err := client.StartContainer(containerID)
			if err != nil {
				fmt.Printf("Error starting container: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Container %s started successfully\n", containerID)
		} else {
			fmt.Printf("Container %s is already running. No action needed.\n", containerID)
		}

	case "start":
		err := client.StartContainer(containerID)
		if err != nil {
			fmt.Printf("Error starting container: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Container %s started successfully\n", containerID)

	case "stop":
		err := client.StopContainer(containerID)
		if err != nil {
			fmt.Printf("Error stopping container: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Container %s stopped successfully\n", containerID)

	case "restart":
		err := client.RestartContainer(containerID)
		if err != nil {
			fmt.Printf("Error restarting container: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Container %s restarted successfully\n", containerID)

	case "logs":
		logs, err := client.GetLogs(containerID)
		if err != nil {
			fmt.Printf("Error fetching logs: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Logs for %s:\n---\n%s\n---\n", containerID, logs)

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}
