package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// SpawnSwarm orchestrates multiple agents inside an isolated Docker bridge network
func (sdc *StuntDockerClient) SpawnSwarm(ctx context.Context, agents []string, mountDir string) error {
	networkName := "stuntnet-01"

	// 1. Create the isolated Virtual Network
	fmt.Printf(">> [Native Engine] Provisioning isolated bridge network: %s\n", networkName)
	netResp, err := sdc.cli.NetworkCreate(ctx, networkName, network.CreateOptions{
		Driver:   "bridge",
		Internal: true, // completely block external internet access for the swarm
	})
	if err != nil {
		// Ignore error if network already exists
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("failed to create StuntNet: %w", err)
		}
	}
	defer func() {
		fmt.Printf(">> [Native Engine] Tearing down StuntNet: %s\n", networkName)
		sdc.cli.NetworkRemove(context.Background(), netResp.ID)
	}()

	var containerIDs []string

	// 2. Spawn each agent as a headless microservice inside the Swarm
	for i, agent := range agents {
		fmt.Printf(">> [Native Engine] Spawning swarm node %d: %s\n", i+1, agent)
		
		agentCmdStr := fmt.Sprintf("npx -y %s", agent)
		
		resp, err := sdc.cli.ContainerCreate(ctx, &container.Config{
			Image: "node:20-alpine",
			Cmd:   []string{"sh", "-c", agentCmdStr},
			Labels: map[string]string{
				"stuntdouble.swarm": "true",
			},
		}, &container.HostConfig{
			CapDrop: []string{"ALL"},
			Binds:   []string{fmt.Sprintf("%s:/workspace", mountDir)},
			NetworkMode: container.NetworkMode(networkName),
		}, nil, nil, fmt.Sprintf("stunt-node-%d", i+1))

		if err != nil {
			return fmt.Errorf("failed to spawn node %s: %w", agent, err)
		}

		if err := sdc.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			return fmt.Errorf("failed to start node %s: %w", agent, err)
		}

		containerIDs = append(containerIDs, resp.ID)
	}

	fmt.Println("\n✅ StuntNet Swarm is Live!")
	fmt.Printf("🔒 Agents can communicate via http://stunt-node-X inside the sandbox.\n")
	fmt.Println("⚠️ External internet access is completely hard-blocked by Docker network rules.")
	
	// Wait for context cancellation (Ctrl+C)
	fmt.Println("\nPress Ctrl+C to safely terminate the swarm and wipe the network...")
	<-ctx.Done()

	// 3. Tear down the swarm
	fmt.Println("\n>> [Native Engine] Terminating all swarm nodes...")
	for _, id := range containerIDs {
		sdc.cli.ContainerRemove(context.Background(), id, container.RemoveOptions{Force: true})
	}

	return nil
}
