package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// StuntDockerClient wraps the native Docker SDK
type StuntDockerClient struct {
	cli *client.Client
}

// NewClient initializes a native connection to the host Docker daemon
func NewClient() (*StuntDockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker daemon: %w", err)
	}
	return &StuntDockerClient{cli: cli}, nil
}

// SpawnIsolatedAgent creates a heavily restricted container natively via API
func (sdc *StuntDockerClient) SpawnIsolatedAgent(ctx context.Context, agentCmd []string, mountDir string) error {
	fmt.Println(">> [Native Engine] Pulling node:20-alpine image...")
	
	reader, err := sdc.cli.ImagePull(ctx, "docker.io/library/node:20-alpine", image.PullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	fmt.Println(">> [Native Engine] Spawning agent with --cap-drop=ALL via CLI proxy stream...")
	
	args := []string{
		"run", "-it", "--rm",
		"--cap-drop=ALL",
		"-v", fmt.Sprintf("%s:/workspace", mountDir),
		"-w", "/workspace",
		"node:20-alpine",
	}
	args = append(args, agentCmd...)

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Print("\033[?25h\033[0m") // Force restore cursor visibility and color
		return fmt.Errorf("agent execution failed: %w", err)
	}

	// Force restore cursor visibility and color just in case the agent crashed or exited abruptly
	fmt.Print("\033[?25h\033[0m")
	return nil
}
