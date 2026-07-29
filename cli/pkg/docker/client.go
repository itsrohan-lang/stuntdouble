package docker

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

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

	// 1. Start the Keploy proxy sidecar in the background. Use a unique name so
	// two repos with the same basename can run concurrently, and avoid publishing
	// Keploy's port on the host. The agent joins the sidecar network namespace
	// directly via --network=container:<name>.
	sidecarName, err := newSidecarName(mountDir)
	if err != nil {
		return fmt.Errorf("generating sidecar name: %w", err)
	}
	sidecarArgs := []string{
		"run", "-d", "--rm",
		"--name", sidecarName,
		"--cap-add=NET_ADMIN", // Keploy requires network capabilities to intercept traffic
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

func newSidecarName(mountDir string) (string, error) {
	suffixBytes := make([]byte, 4)
	if _, err := rand.Read(suffixBytes); err != nil {
		return "", err
	}
	base := sanitizeContainerNamePart(filepath.Base(mountDir))
	return fmt.Sprintf("stunt-keploy-%s-%s", base, hex.EncodeToString(suffixBytes)), nil
}

func sanitizeContainerNamePart(input string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(input) {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '.', r == '_', r == '-':
			b.WriteRune(r)
		case unicode.IsSpace(r):
			b.WriteRune('-')
		default:
			b.WriteRune('-')
		}
	}

	name := strings.Trim(b.String(), ".-_")
	if name == "" {
		return "workspace"
	}
	if len(name) > 40 {
		name = strings.TrimRight(name[:40], ".-_")
		if name == "" {
			return "workspace"
		}
	}
	return name
}
