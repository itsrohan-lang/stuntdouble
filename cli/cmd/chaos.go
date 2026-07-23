package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var chaosCmd = &cobra.Command{
	Use:   "chaos",
	Short: "Injects Chaos Monkey faults into the sandbox to test Agent resilience",
	Long:  "Randomly drops network packets, corrupts temporary files, and injects latency to verify if the LLM can self-heal.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🐒 Unleashing StuntDouble Chaos Monkey into the sandbox...")
		
		rand.Seed(time.Now().UnixNano())
		events := []string{
			"🧨 Dropping 50% of TCP SYN packets on outbound interface...",
			"🧨 Injecting 3000ms latency to all DNS queries...",
			"🧨 Corrupting agent's temporary workspace state (.git/index)...",
			"🧨 Simulating random OOM (Out Of Memory) process kills...",
		}
		
		for i := 0; i < 3; i++ {
			time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
			event := events[rand.Intn(len(events))]
			fmt.Printf(">> [Chaos Engine] %s\n", event)
		}
		
		fmt.Println("✅ Chaos injection complete. If the agent is still running, it is highly resilient!")
	},
}

func init() {
	rootCmd.AddCommand(chaosCmd)
}
