// Package ebpf is intended to provide kernel-level network egress filtering
// for sandboxed agent containers.
//
// It is not implemented on any platform. AttachInterceptor returns
// ErrUnsupported everywhere. Callers must treat that error as fatal: without a
// working interceptor the sandbox applies no network restrictions at all, and
// proceeding would give the caller a false sense of containment.
//
// The intended Linux implementation is sketched in bpf_prog.c (a
// cgroup_skb/egress filter that drops traffic to a configurable port set). It
// is not compiled into the binary.
package ebpf

import "errors"

// ErrUnsupported reports that no network enforcement backend is available on
// this platform.
var ErrUnsupported = errors.New("ebpf: kernel-level network enforcement is not implemented; see docs/ENFORCEMENT.md")
