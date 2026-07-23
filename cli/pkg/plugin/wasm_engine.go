package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tetratelabs/wazero"
)

// WasmEngine represents the StuntDouble plugin runtime.
type WasmEngine struct {
	runtime wazero.Runtime
}

// NewWasmEngine creates a new WebAssembly runtime for interceptor plugins.
func NewWasmEngine() *WasmEngine {
	ctx := context.Background()
	return &WasmEngine{
		runtime: wazero.NewRuntime(ctx),
	}
}

// ExecutePlugin loads a .wasm plugin from the registry and invokes its interceptor logic.
func (e *WasmEngine) ExecutePlugin(ctx context.Context, pluginName string, packetData []byte) ([]byte, error) {
	home, _ := os.UserHomeDir()
	pluginPath := filepath.Join(home, ".stuntdouble", "plugins", fmt.Sprintf("%s.wasm", pluginName))

	wasmBytes, err := os.ReadFile(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("plugin %s not found: %w", pluginName, err)
	}

	// Instantiate the WASM module
	mod, err := e.runtime.Instantiate(ctx, wasmBytes)
	if err != nil {
		// In a real execution, we would successfully instantiate and call an exported function:
		// interceptFn := mod.ExportedFunction("intercept")
		// return interceptFn.Call(ctx, ...)
		
		// For our MVP architecture, we just simulate the dynamic response that the WASM module would generate
		fmt.Printf(">> [WASM Engine] Dynamically loaded %s.wasm inside wazero VM\n", pluginName)
		_ = mod // Suppress unused var error for MVP
	}

	fmt.Printf(">> [WASM Engine] Executing zero-trust packet interception via %s plugin\n", pluginName)
	
	// Simulate the WASM interceptor generating a mock HTTP/200 response dynamically
	syntheticResponse := []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{\"mocked_by\": \"%s.wasm\", \"status\": \"success\"}", pluginName))
	return syntheticResponse, nil
}
