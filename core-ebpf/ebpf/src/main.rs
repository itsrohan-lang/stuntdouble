#![no_std]
#![no_main]

use aya_ebpf::{
    macros::cgroup_skb,
    programs::SkBuffContext,
};

/// This is the actual Kernel-space eBPF probe that will run natively inside the Linux Kernel.
/// It intercepts outgoing network packets at the cgroup level before they reach the network interface.
#[cgroup_skb]
pub fn stuntdouble_egress_guard(ctx: SkBuffContext) -> i32 {
    // Return 1 to ALLOW the packet, 0 to DROP (Block) the packet.
    
    // TODO: Parse the packet headers (IPv4/IPv6 + TCP/UDP).
    // Extract the destination IP and Port.
    // Query the eBPF Map (shared with user-space Control Plane) to see if this port is blocked.
    // If blocked -> return 0 (DROP).
    // Else -> return 1 (ALLOW).

    1 
}

#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    unsafe { core::hint::unreachable_unchecked() }
}
