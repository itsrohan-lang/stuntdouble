package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
	
	reader, err := sdc.cli.ImagePull(ctx, "docker.io/library/node:20-alpine", types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	fmt.Println(">> [Native Engine] Spawning agent with --cap-drop=ALL via API...")
	
	// Define the strict container configuration
	resp, err := sdc.cli.ContainerCreate(ctx, &container.Config{
		Image:        "node:20-alpine",
		Cmd:          agentCmd,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		CapDrop: []string{"ALL"},
		Binds:   []string{fmt.Sprintf("%s:/workspace", mountDir)},
	}, nil, nil, "")
	
	if err != nil {
		return err
	}

	if err := sdc.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	// This is where we would stream the logs and wait for completion in the future
	fmt.Printf("✅ Agent spawned natively! Container ID: %s\n", resp.ID[:12])
	
	return nil
}
