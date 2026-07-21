.PHONY: all build ebpf clean

all: ebpf build

ebpf:
	@echo ">> Compiling eBPF C-code to BPF bytecode using Cilium bpf2go..."
	@cd pkg/ebpf && go generate ./...

build:
	@echo ">> Compiling StuntDouble CLI for your architecture..."
	@go build -o sd main.go

clean:
	@rm -f sd
	@rm -f pkg/ebpf/*_bpfel.go pkg/ebpf/*_bpfel.o
	@rm -f pkg/ebpf/*_bpfeb.go pkg/ebpf/*_bpfeb.o
