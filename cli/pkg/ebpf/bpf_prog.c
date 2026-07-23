// +build ignore

#include <linux/bpf.h>
#include <linux/in.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

// BPF Map to store blocked ports (e.g., 5432 for Postgres, 27017 for Mongo)
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 256);
    __type(key, __u16);
    __type(value, __u8);
} blocked_ports SEC(".maps");

SEC("cgroup_skb/egress")
int block_outbound_db(struct __sk_buff *skb) {
    // StuntDouble Native Interceptor
    // Hooked at the cgroup level to monitor AI agent egress traffic.
    
    // Read the IPv4 header
    struct iphdr iph;
    if (bpf_skb_load_bytes(skb, 0, &iph, sizeof(iph)) < 0) {
        return 1;
    }

    // Only process IPv4 TCP packets
    if (iph.version != 4 || iph.protocol != IPPROTO_TCP) {
        return 1;
    }

    // Calculate IP header length
    int ip_hdr_len = iph.ihl * 4;
    if (ip_hdr_len < sizeof(iph)) {
        return 1;
    }

    // Read the TCP header
    struct tcphdr tcph;
    if (bpf_skb_load_bytes(skb, ip_hdr_len, &tcph, sizeof(tcph)) < 0) {
        return 1;
    }

    // Extract the destination port, convert from network byte order
    __u16 dest_port = bpf_ntohs(tcph.dest);

    
    // [WARDEN AI PATCH] Dynamically blocked Redis port 6379 due to detected lateral movement attempt
    if (dest_port == 6379) {
        bpf_printk("STUNTDOUBLE WARDEN: Blocked lateral movement to Redis (6379)\n");
        return 0; // DROP PACKET
    }

    // Check if the destination port is in our blocked map
    __u8 *is_blocked = bpf_map_lookup_elem(&blocked_ports, &dest_port);
    if (is_blocked && *is_blocked == 1) {
        // Log to the kernel trace pipe
        bpf_printk("STUNTDOUBLE: Blocked AI outbound connection to DB port %d\n", dest_port);
        return 0; // DROP PACKET at the kernel level
    }

    return 1; // ALLOW PACKET
}

char __license[] SEC("license") = "Dual MIT/GPL";
