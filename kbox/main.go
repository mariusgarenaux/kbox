package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

func printUsage() {
	fmt.Println("KBox - Jupyter Kernel Container Management CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  kbox <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  start   - Setup and start a kernel container (requires connection_file and tag)")
	fmt.Println("  install - Create a Jupyter kernelspec for a Docker image from a built-in name or kernelspec directory")
	fmt.Println("  list    - List all built-in kernels available for installation")
	fmt.Println("\nExample:")
	fmt.Println("  kbox start ./kernel-connection.json my-image:latest --argv python -m ipykernel_launcher -f /connection_file.json")
	fmt.Println("  kbox install my-image:latest python3 [new_kernel_name]")
	fmt.Println("  kbox install my-image:latest /path/to/kernelspec [new_kernel_name]")
	fmt.Println("  kbox list")
}

type BuiltInKernel struct {
	DisplayName string   `json:"display_name"`
	Language    string   `json:"language"`
	Argv        []string `json:"argv"`
}

var knownKernels = map[string]BuiltInKernel{
	"python3": {
		DisplayName: "Python 3",
		Language:    "python",
		Argv:        []string{"python", "-m", "ipykernel_launcher", "-f", "{connection_file}"},
	},
	"python": {
		DisplayName: "Python 3",
		Language:    "python",
		Argv:        []string{"python", "-m", "ipykernel_launcher", "-f", "{connection_file}"},
	},
	"bash": {
		DisplayName: "Bash",
		Language:    "bash",
		Argv:        []string{"python", "-m", "bash_kernel", "-f", "{connection_file}"},
	},
	"pydantic_ai": {
		DisplayName: "Pydantic AI Agent",
		Language:    "text",
		Argv:        []string{"python", "-m", "pydantic_ai_kernel", "-f", "{connection_file}"},
	},
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

	if cmd == "start" {
		if len(os.Args) < 5 {
			fmt.Println("Error: 'start' requires connection_file, docker_tag and --argv")
			fmt.Println("Usage: kbox start <connection_file> <docker_tag> --argv <docker_run_cmd...>")
			os.Exit(1)
		}
		connFile := os.Args[2]
		tag := os.Args[3]

		// Find the index of the --argv flag
		argvIdx := -1
		for i := 4; i < len(os.Args); i++ {
			if os.Args[i] == "--argv" {
				argvIdx = i
				break
			}
		}

		if argvIdx == -1 {
			fmt.Println("Error: '--argv' flag is required")
			os.Exit(1)
		}

		if argvIdx+1 >= len(os.Args) {
			fmt.Println("Error: '--argv' requires at least one argument")
			os.Exit(1)
		}

		// Load all parameters after --argv into a list of strings
		dockerRunCmd := os.Args[argvIdx+1:]

		var containerID string
		containerID, err = client.ConnectKBox(connFile, tag, dockerRunCmd)
		if err != nil {
			fmt.Printf("Error during connection setup: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Connection established! Container %s is now running. Waiting for shutdown signal...\n", containerID)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Printf("\nShutdown signal received. Stopping container %s...\n", containerID)
		if err := client.StopContainer(containerID); err != nil {
			fmt.Printf("Error stopping container: %v\n", err)
		} else {
			fmt.Println("Container stopped successfully.")
		}
	} else if cmd == "install" {
		if len(os.Args) < 4 {
			fmt.Println("Error: 'install' requires docker_tag and kernelspec path")
			fmt.Println("Usage: kbox install <docker_tag> <path_to_kernelspec_dir> [new_kernel_name]")
			os.Exit(1)
		}
		tag := os.Args[2]
		kernelInput := os.Args[3]

		var sourceKernelspec map[string]interface{}
		var kernelName string

		// Check if the input is a known built-in kernel
		if builtIn, ok := knownKernels[kernelInput]; ok {
			kernelName = kernelInput
			sourceKernelspec = map[string]interface{}{
				"display_name": builtIn.DisplayName,
				"language":     builtIn.Language,
				"argv":         builtIn.Argv,
			}
		} else {
			// If not a built-in name, it must be a path
			if !filepath.IsAbs(kernelInput) && !strings.Contains(kernelInput, "/") {
				fmt.Println("Error: 'install' requires a known built-in kernel name or a path to a kernelspec directory")
				fmt.Println("Usage: kbox install <docker_tag> <built_in_name|path_to_kernelspec_dir> [new_kernel_name]")
				os.Exit(1)
			}

			// Verify the path exists and is a directory
			if info, err := os.Stat(kernelInput); err == nil && info.IsDir() {
				kernelName = filepath.Base(kernelInput)
				kernelJsonPath := filepath.Join(kernelInput, "kernel.json")
				if data, err := os.ReadFile(kernelJsonPath); err == nil {
					json.Unmarshal(data, &sourceKernelspec)
				}
			} else {
				fmt.Printf("Error: path %s is not a directory or does not exist\n", kernelInput)
				os.Exit(1)
			}
		}

		// Determine base kernels directory based on OS for the destination
		baseKernelsDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "jupyter", "kernels")
		if runtime.GOOS == "darwin" {
			baseKernelsDir = filepath.Join(os.Getenv("HOME"), "Library", "Jupyter", "kernels")
		}

		// Create a unique name for the kbox wrapped kernel to avoid overwriting the source
		kboxKernelName := "kbox-" + kernelName
		if len(os.Args) >= 5 {
			kboxKernelName = os.Args[4]
		}

		execPath, err := os.Executable()
		if err != nil {
			fmt.Printf("Error getting executable path: %v\n", err)
			os.Exit(1)
		}

		// Default argv if no source kernel was found
		argv := []string{execPath, "start", "{connection_file}", tag, "--argv", "python", "-m", "ipykernel_launcher", "-f", "/connection_file.json"}
		displayName := fmt.Sprintf("KBox - %s", kernelName)
		language := "python"

		// If we found a source kernel, wrap its argv instead
		if sourceKernelspec != nil {
			if srcArgv, ok := sourceKernelspec["argv"].([]interface{}); ok {
				// Convert []interface{} to []string and replace {connection_file} with /connection_file.json
				wrappedArgv := []string{execPath, "start", "{connection_file}", tag, "--argv"}
				for _, a := range srcArgv {
					val := fmt.Sprintf("%v", a)
					val = strings.ReplaceAll(val, "{connection_file}", "/connection_file.json")
					wrappedArgv = append(wrappedArgv, val)
				}
				argv = wrappedArgv
			} else if srcArgv, ok := sourceKernelspec["argv"].([]string); ok {
				// Handle case where it's already a string slice (from built-in)
				wrappedArgv := []string{execPath, "start", "{connection_file}", tag, "--argv"}
				for _, a := range srcArgv {
					val := strings.ReplaceAll(a, "{connection_file}", "/connection_file.json")
					wrappedArgv = append(wrappedArgv, val)
				}
				argv = wrappedArgv
			}
			if dn, ok := sourceKernelspec["display_name"].(string); ok {
				displayName = fmt.Sprintf("KBox - %s", dn)
			}
			if lang, ok := sourceKernelspec["language"].(string); ok {
				language = lang
			}
		}

		newKernelspec := map[string]interface{}{
			"argv":         argv,
			"display_name": displayName,
			"language":     language,
		}

		kernelspecData, err := json.MarshalIndent(newKernelspec, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling kernelspec: %v\n", err)
			os.Exit(1)
		}

		kernelspecDir := filepath.Join(baseKernelsDir, kboxKernelName)
		if err := os.MkdirAll(kernelspecDir, 0755); err != nil {
			fmt.Printf("Error creating kernelspec directory: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(filepath.Join(kernelspecDir, "kernel.json"), kernelspecData, 0644); err != nil {
			fmt.Printf("Error writing kernel.json: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Jupyter Kernelspec created for '%s' at %s\n", kboxKernelName, kernelspecDir)
	} else if cmd == "list" {
		fmt.Println("Available built-in kernels:")
		for name, kernel := range knownKernels {
			fmt.Printf("  - %s (%s): %s\n    Argv: %v\n", name, kernel.Language, kernel.DisplayName, kernel.Argv)
		}
	} else {
		printUsage()
	}

}

