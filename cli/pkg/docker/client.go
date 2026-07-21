package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/moby/term"
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

	fmt.Println(">> [Native Engine] Spawning agent with --cap-drop=ALL via API...")
	
	// Define the strict container configuration
	resp, err := sdc.cli.ContainerCreate(ctx, &container.Config{
		Image:        "node:20-alpine",
		Cmd:          agentCmd,
		Tty:          true,
		OpenStdin:    true,
		StdinOnce:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		CapDrop: []string{"ALL"},
		Binds:   []string{fmt.Sprintf("%s:/workspace", mountDir)},
	}, nil, nil, "")
	
	if err != nil {
		return err
	}

	// Attach to the container streams before starting
	attachResp, err := sdc.cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}
	defer attachResp.Close()

	if err := sdc.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	fmt.Printf("✅ Agent spawned natively! Container ID: %s\n", resp.ID[:12])
	
	// Set the host terminal into raw mode so arrow keys and UI elements pass through cleanly
	inFd, isTerm := term.GetFdInfo(os.Stdin)
	if isTerm {
		state, err := term.SetRawTerminal(inFd)
		if err == nil {
			defer term.RestoreTerminal(inFd, state)
		}
	}

	// Stream TTY interactively
	go func() {
		io.Copy(os.Stdout, attachResp.Reader)
	}()
	go func() {
		io.Copy(attachResp.Conn, os.Stdin)
	}()

	// Wait for container to exit natively
	statusCh, errWaitCh := sdc.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errWaitCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}
