package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type KBoxClient struct {
	cli *client.Client
	ctx context.Context
}

func NewKBoxClient() (*KBoxClient, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &KBoxClient{cli: cli, ctx: ctx}, nil
}

func (k *KBoxClient) IsRunning(containerID string) (bool, error) {
	inspect, err := k.cli.ContainerInspect(k.ctx, containerID)
	if err != nil {
		return false, err
	}
	return inspect.State.Running, nil
}

func (k *KBoxClient) StartContainer(containerID string) error {
	return k.cli.ContainerStart(k.ctx, containerID, container.StartOptions{})
}

func (k *KBoxClient) StopContainer(containerID string) error {
	return k.cli.ContainerStop(k.ctx, containerID, container.StopOptions{})
}

func (k *KBoxClient) RestartContainer(containerID string) error {
	return k.cli.ContainerRestart(k.ctx, containerID, container.StopOptions{})
}

func (k *KBoxClient) GetLogs(containerID string) (string, error) {
	out, err := k.cli.ContainerLogs(k.ctx, containerID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err
	}
	defer out.Close()
	
	buf := make([]byte, 4096)
	n, err := out.Read(buf)
	if err != nil && n == 0 {
		return "", err
	}
	return string(buf[:n]), nil
}

func (k *KBoxClient) ListImages() ([]image.Summary, error) {
	images, err := k.cli.ImageList(k.ctx, image.ListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (k *KBoxClient) ConnectKBox(connectionFile string, dockerTag string, dockerRun []string) (string, error) {
	// 1. Ensure connectionFile is an absolute path
	absConnFile, err := filepath.Abs(connectionFile)
	if err != nil {
		return "", fmt.Errorf("error resolving absolute path for connection file: %w", err)
	}

	
	// Read and modify connection file
	data, err := os.ReadFile(absConnFile)
	if err != nil {
		return "", fmt.Errorf("error reading connection file: %w", err)
	}
	
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}
	
	// Helper function to append logs to kbox_debug.log
	logDebug := func(msg string) {
		f, err := os.OpenFile("kbox_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString(msg + "\n")
		}
	}

	// Change host IP to 0.0.0.0 for port-forwarding
	config["ip"] = "0.0.0.0"
	modifiedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}
	
	logDebug(fmt.Sprintf("DEBUG: Modified config : %s", modifiedData))
	// Update connection file directly to preserve inode for Docker mounts
	if err := os.WriteFile(absConnFile, modifiedData,  0644); err != nil {
		// logDebug(fmt.Sprintf("ERROR: while writing connection file to %s", absConnFile))
		return "", fmt.Errorf("error writing connection file: %w", err)
	}
	
	// logDebug("DEBUG: Connection file successfully updated to 0.0.0.0")
	
	// Extract ports and key
	shellPort := fmt.Sprintf("%v", config["shell_port"])
	iopubPort := fmt.Sprintf("%v", config["iopub_port"])
	stdinPort := fmt.Sprintf("%v", config["stdin_port"])
	hbPort := fmt.Sprintf("%v", config["hb_port"])
	
	// Handle missing key by using a fallback or the identity
	key := fmt.Sprintf("%v", config["key"])
	if key == "<nil>" {
		// logDebug("Missing Key in connection file")
		key = "kbox-kernel"
	}
	
	containerName := "kernel-" + key
	
	// 2. Setup Port Bindings Correctly
	portBinds := nat.PortMap{
		nat.Port(shellPort + "/tcp"): []nat.PortBinding{{HostPort: shellPort}},
		nat.Port(iopubPort + "/tcp"): []nat.PortBinding{{HostPort: iopubPort}},
		nat.Port(stdinPort + "/tcp"): []nat.PortBinding{{HostPort: stdinPort}},
		nat.Port(hbPort + "/tcp"):    []nat.PortBinding{{HostPort: hbPort}},
	}
	
	cwd, _ := os.Getwd()
	logDebug(fmt.Sprintf("docker Run : %s", dockerRun))
	// 3. Create Container
	resp, err := k.cli.ContainerCreate(k.ctx, &container.Config{
		Image: dockerTag,
		Cmd:   dockerRun,
		// Cmd:   []string{"cat", "/connection_file.json"},
		// Cmd:   []string{"cat", "/notebook/bonjour.md"},
		// Cmd:   []string{"tail", "-f", "/dev/null"},
		Tty:   true,
		WorkingDir: "/notebook",
	}, &container.HostConfig{
		PortBindings: portBinds,
		AutoRemove:   true,
		Binds: []string{
			fmt.Sprintf("%s:/connection_file.json", absConnFile),
			fmt.Sprintf("%s:/notebook", cwd),
		},
	}, nil, nil, containerName)

	// logDebug(fmt.Sprintf("Container Response : %s", resp))
	if err != nil {
		// logDebug("Error while creating the container")
		return "", fmt.Errorf("error creating container: %w", err)
	}
	
	// 4. Start Container
	if err := k.cli.ContainerStart(k.ctx, resp.ID, container.StartOptions{}); err != nil {
		// logDebug(fmt.Sprintf("Error while starting the container : %s", err))
		return "", fmt.Errorf("error starting container: %w", err)
	}

	// logDebug("Successfully started the container and verified it is ready")
	return resp.ID, nil
}
