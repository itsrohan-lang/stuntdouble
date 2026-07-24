#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
#include <linux/tcp.h>
#include <bpf/bpf_helpers.h>

/* 
 * StuntDouble eBPF Kernel Probe
 * Intercepts outbound egress traffic from AI Agent Docker containers (cgroup_skb).
 * Returns TC_ACT_SHOT (drop) if traffic violates the zero-trust policy.
 */

SEC("cgroup_skb/egress")
int stuntdouble_egress_filter(struct __sk_buff *skb) {
    // Read the IP protocol from the socket buffer
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;

    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return 1; // Pass allowed

    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return 1; // Only inspect IPv4

    struct iphdr *ip = data + sizeof(*eth);
    if ((void *)(ip + 1) > data_end)
        return 1;

    // Only inspect TCP traffic (e.g. HTTP, Postgres)
    if (ip->protocol != IPPROTO_TCP)
        return 1;

    struct tcphdr *tcp = (void *)ip + sizeof(*ip);
    if ((void *)(tcp + 1) > data_end)
        return 1;

    __u16 dest_port = bpf_ntohs(tcp->dest);

    // ZERO-TRUST BLOCKED PORTS
    // 5432: PostgreSQL
    // 27017: MongoDB
    // 3306: MySQL
    // 6379: Redis
    // If the agent tries to hit a database directly without Keploy, kill it at the kernel level.
    if (dest_port == 5432 || dest_port == 27017 || dest_port == 3306 || dest_port == 6379) {
        bpf_printk("[StuntDouble] 🚨 Agent attempted illegal DB connection on port %d! PACKET DROPPED.\n", dest_port);
        
        // Return 0 to drop the packet natively in the Linux kernel
        return 0;
    }

    // Pass traffic otherwise
    return 1;
}

char _license[] SEC("license") = "GPL";
