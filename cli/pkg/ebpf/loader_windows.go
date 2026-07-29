//go:build windows

package ebpf

// Network enforcement on Windows would require a Windows Filtering Platform
// callout driver and an accompanying service. The service referenced by
// earlier versions of this file (StuntDoubleWFP) is not built, signed, or
// shipped by this repository, so there is nothing to start.

// Interceptor would hold a handle to the WFP service.
type Interceptor struct{}

// AttachInterceptor is not implemented on Windows. It returns ErrUnsupported so
// callers refuse to run workloads that would otherwise be unrestricted.
func AttachInterceptor(cgroupPath string) (*Interceptor, error) {
	return nil, ErrUnsupported
}

// Detach is a no-op while AttachInterceptor always fails.
func (i *Interceptor) Detach() {}
