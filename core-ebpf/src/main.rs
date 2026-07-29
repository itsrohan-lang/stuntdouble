//! StuntDouble eBPF engine.
//!
//! This binary is a placeholder. It does not load or attach any eBPF program:
//! the bytecode it would load (`ebpf-bytecode/stuntdouble.bpf.o`) is not built
//! by this repository, and no build step produces it.
//!
//! It exits with an error rather than idling, so that anything supervising it
//! (see `control-plane/main.go`) observes a failure instead of assuming
//! enforcement is active.

use log::error;

fn main() -> Result<(), anyhow::Error> {
    env_logger::init();

    error!(
        "kernel-level network enforcement is not implemented: no compiled eBPF \
         bytecode is available to load. See docs/ENFORCEMENT.md."
    );

    anyhow::bail!("ebpf engine unimplemented")
}
