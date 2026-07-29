//go:build linux

package ebpf

// The eBPF program in bpf_prog.c is not yet compiled or loaded. Generating the
// Go bindings requires clang and kernel headers:
//
//	go generate ./pkg/ebpf
//
// which is not wired into the build yet. Until loadBpfObjects exists,
// AttachInterceptor reports that enforcement is unavailable rather than
// pretending to have attached anything.
//
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang bpf bpf_prog.c -- -I../headers

// Interceptor would hold the loaded BPF objects and the cgroup link.
type Interceptor struct{}

// AttachInterceptor is intended to load bpf_prog.c into the kernel and attach
// it to the cgroup at cgroupPath as a cgroup_skb/egress filter.
//
// It is not implemented. It returns ErrUnsupported so callers refuse to run
// workloads that would otherwise be unrestricted.
func AttachInterceptor(cgroupPath string) (*Interceptor, error) {
	return nil, ErrUnsupported
}

// Detach releases the kernel hooks. It is a no-op while AttachInterceptor
// always fails.
func (i *Interceptor) Detach() {}
