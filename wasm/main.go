//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

// PolicyEngine encapsulates the logic for evaluating network/file rules
type PolicyEngine struct {
	Mode string
}

func (pe *PolicyEngine) EvaluateRequest(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf(map[string]interface{}{
			"allowed": false,
			"error":   "missing request payload",
		})
	}

	reqPayload := args[0].String()
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(reqPayload), &req); err != nil {
		return js.ValueOf(map[string]interface{}{
			"allowed": false,
			"error":   "invalid JSON payload",
		})
	}

	// Simple simulation of policy evaluation logic in WASM
	// In a real port, this would link to github.com/itsrohan-lang/stuntdouble/cli/pkg/api
	target, ok := req["target"].(string)
	if !ok {
		return js.ValueOf(map[string]interface{}{
			"allowed": false,
			"error":   "target field missing",
		})
	}

	// Example static rule for WASM engine
	allowed := true
	if target == "api.stripe.com" && pe.Mode != "audit" {
		allowed = false
	}

	fmt.Printf("[WASM Engine] Evaluated request to %s: allowed=%v\n", target, allowed)

	return js.ValueOf(map[string]interface{}{
		"allowed": allowed,
		"target":  target,
		"engine":  "stuntdouble-wasm",
	})
}

func main() {
	c := make(chan struct{}, 0)

	engine := &PolicyEngine{Mode: "block"}

	js.Global().Set("StuntDouble", js.ValueOf(map[string]interface{}{
		"evaluate": js.FuncOf(engine.EvaluateRequest),
		"version":  "1.0.0-wasm",
	}))

	fmt.Println("StuntDouble WASM Engine initialized and mounted to global scope.")
	<-c
}
