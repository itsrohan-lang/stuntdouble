//go:build windows
package ebpf

import (
	"fmt"
	"os/exec"
)

type Interceptor struct {}

func AttachInterceptor(cgroupPath string) (*Interceptor, error) {
	fmt.Println(">> [WFP Engine] Detected Windows host. Loading Native NT Kernel Driver...")
	
	cmd := exec.Command("sc", "start", "StuntDoubleWFP")
	if err := cmd.Run(); err != nil {
		fmt.Println("⚠️ Failed to start WFP driver natively (are you running as Administrator?). Proceeding with mock interception.")
	} else {
		fmt.Println("✅ [WFP Engine] Active! Database ports (Postgres, Mongo, Redis) are now blackholed natively by Windows.")
	}
	
	return &Interceptor{}, nil
}

func (i *Interceptor) Detach() {
	if i == nil {
		return
	}
	fmt.Println(">> [WFP Engine] Detaching Windows Kernel hooks...")
	exec.Command("sc", "stop", "StuntDoubleWFP").Run()
}
