.PHONY: all build-cli build-ebpf build-mac build-windows build-wasm start-soc

all: build-cli build-ebpf build-mac build-windows build-wasm

build-cli:
	@echo "🔨 Building Go CLI Wrapper..."
	cd cli && go build -o sd main.go

build-ebpf:
	@echo "🦀 Building Rust eBPF Probes (Linux)..."
	cd core-ebpf && cargo build --release

build-mac:
	@echo "🍎 Building macOS Endpoint Security Framework Interceptor..."
	clang -Wall -O2 mac/src/esf_interceptor.c -o mac/esf_interceptor -framework EndpointSecurity -framework Foundation

build-windows:
	@echo "🪟 Building Windows Filtering Platform Interceptor..."
	# Assuming MinGW or MSVC is installed
	x86_64-w64-mingw32-g++ -Wall -O2 windows/src/wfp_interceptor.cpp -o windows/wfp_interceptor.exe -lwintrust -lws2_32 -lfwpuclnt

build-wasm:
	@echo "🕸️ Building WebAssembly Edge Policies..."
	cd wasm && GOOS=js GOARCH=wasm go build -o stuntdouble.wasm main.go

start-soc:
	@echo "🐳 Starting the StuntDouble SOC (Control Plane + Dashboard) via Docker Compose..."
	docker-compose up -d

stop-soc:
	@echo "🛑 Stopping the StuntDouble SOC..."
	docker-compose down
