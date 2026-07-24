//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

// PolicyEngine encapsulates the logic for evaluating network/file rules
type PolicyEngine struct {
	Mode          string
	BlockedTarget []string
}

func (pe *PolicyEngine) SetPolicy(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return false
	}
	policyJSON := args[0].String()
	var policy map[string]interface{}
	if err := json.Unmarshal([]byte(policyJSON), &policy); err != nil {
		fmt.Println("[WASM Engine] Error parsing policy:", err)
		return false
	}

	if mode, ok := policy["mode"].(string); ok {
		pe.Mode = mode
	}

	if blocked, ok := policy["blocked"].([]interface{}); ok {
		pe.BlockedTarget = []string{}
		for _, b := range blocked {
			if str, ok := b.(string); ok {
				pe.BlockedTarget = append(pe.BlockedTarget, str)
			}
		}
	}
	fmt.Println("[WASM Engine] Policy updated successfully")
	return true
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

	target, ok := req["target"].(string)
	if !ok {
		return js.ValueOf(map[string]interface{}{
			"allowed": false,
			"error":   "target field missing",
		})
	}

	allowed := true
	if pe.Mode != "audit" {
		for _, blocked := range pe.BlockedTarget {
			if target == blocked {
				allowed = false
				break
			}
		}
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

	engine := &PolicyEngine{Mode: "block", BlockedTarget: []string{}}

	js.Global().Set("StuntDouble", js.ValueOf(map[string]interface{}{
		"evaluate":  js.FuncOf(engine.EvaluateRequest),
		"setPolicy": js.FuncOf(engine.SetPolicy),
		"version":   "1.0.0-wasm",
	}))

	fmt.Println("StuntDouble WASM Engine initialized and mounted to global scope.")
	<-c
}
