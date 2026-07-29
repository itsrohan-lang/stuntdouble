//go:build darwin

package ebpf

// Network enforcement on macOS would require an Endpoint Security (ESF) system
// extension. The helper binary referenced by earlier versions of this file
// (mac/src/esf_interceptor) is not built, signed, or shipped by this
// repository, so there is nothing to attach to.

// Interceptor would hold a handle to the ESF helper process.
type Interceptor struct{}

// AttachInterceptor is not implemented on macOS. It returns ErrUnsupported so
// callers refuse to run workloads that would otherwise be unrestricted.
func AttachInterceptor(cgroupPath string) (*Interceptor, error) {
	return nil, ErrUnsupported
}

// Detach is a no-op while AttachInterceptor always fails.
func (i *Interceptor) Detach() {}
