package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

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
func (sdc *StuntDockerClient) SpawnIsolatedAgent(ctx context.Context, agentCmd []string, mountDir string, envImage string) error {
	if envImage == "" {
		envImage = "node:20-alpine" // Default to Node.js
	}
	
	fmt.Printf(">> [Native Engine] Pulling %s image...\n", envImage)
	
	reader, err := sdc.cli.ImagePull(ctx, "docker.io/library/"+envImage, image.PullOptions{})
	if err != nil {
		// If docker.io fails, fallback to bare string (in case they passed a full URI)
		reader, err = sdc.cli.ImagePull(ctx, envImage, image.PullOptions{})
		if err != nil {
			return err
		}
	}
	io.Copy(os.Stdout, reader)

	fmt.Println(">> [Stunt Layer] Injecting Keploy proxy sidecar...")

	// 1. Start the Keploy proxy sidecar in the background
	sidecarName := "stunt-keploy-sidecar-" + filepath.Base(mountDir)
	sidecarArgs := []string{
		"run", "-d", "--rm",
		"--name", sidecarName,
		"--cap-add=NET_ADMIN", // Keploy requires network capabilities to intercept traffic
		"-p", "16789:16789",
		"ghcr.io/keploy/keploy:v2",
	}
	
	sidecarCmd := exec.CommandContext(ctx, "docker", sidecarArgs...)
	if err := sidecarCmd.Run(); err != nil {
		return fmt.Errorf("failed to inject keploy sidecar: %w", err)
	}

	// Ensure the sidecar is cleaned up after the agent finishes
	defer func() {
		fmt.Println(">> [Stunt Layer] Tearing down Keploy sidecar...")
		exec.Command("docker", "kill", sidecarName).Run()
	}()

	fmt.Printf(">> [Native Engine] Spawning %s agent with --cap-drop=ALL attached to sidecar network...\n", envImage)
	
	// 2. Start the Agent container, attaching its network namespace to the sidecar
	args := []string{
		"run", "-it", "--rm",
		"--cap-drop=ALL",
		"--memory=2g",
		"--cpus=1.0",
		"-e", "ANTHROPIC_API_KEY",
		"-e", "OPENAI_API_KEY",
		fmt.Sprintf("--network=container:%s", sidecarName),
		"-v", fmt.Sprintf("%s:/workspace", mountDir),
		"-w", "/workspace",
		envImage,
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

	fmt.Print("\033[?25h\033[0m")
	return nil
}
