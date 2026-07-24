# 🏗️ StuntDouble Architecture

StuntDouble is designed to intercept and monitor AI agent actions at the lowest possible level (the kernel) while providing a high-level, human-readable SOC dashboard for the security team.

```mermaid
graph TD
    subgraph "VS Code Environment"
        A[AI Agent e.g. Claude Code] -->|Shell Commands| B(StuntDouble CLI)
        E[VS Code Extension] -.->|Polls| C
    end

    subgraph "Linux Host"
        B -.->|Docker Exec| F[Isolated Docker Container]
        F -->|cgroup_skb Outbound| G{Rust eBPF Probe}
    end

    subgraph "Enterprise Network"
        G -->|Block/Allow| C[Golang Control Plane :4439]
        C -->|Mock Response| K[Keploy Mock Engine]
        C <-->|Audit & Policies| DB[(SQLite / Postgres)]
    end
    
    subgraph "Security Team"
        D[Next.js Dashboard] <-->|GraphQL / REST| C
    end

    style G fill:#ef4444,stroke:#333,stroke-width:2px,color:#fff
    style C fill:#00f0ff,stroke:#333,stroke-width:2px,color:#000
    style K fill:#8a2be2,stroke:#333,stroke-width:2px,color:#fff
```

### Components
1. **VS Code Extension & CLI**: Wraps the AI agent and forces it to run inside a Docker container.
2. **Rust eBPF Probe**: Sits in the Linux kernel and intercepts `cgroup_skb` network packets before they leave the container.
3. **Golang Control Plane**: A centralized proxy that receives telemetry from the eBPF probe, enforces JSON policies, and writes to an SQLite audit database.
4. **Keploy Mock Engine**: Injects ghost responses for blocked APIs so agents don't crash.
5. **Next.js Dashboard**: A React frontend for the CTO to view live analytics and block logs.
