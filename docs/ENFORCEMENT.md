# StuntDouble Isolation & Security Model

## Overview

StuntDouble is designed to run autonomous AI agents inside containerized, low-privilege sandbox environments.

---

## Current Isolation Guarantees

### 1. Docker Container Isolation
When launching an agent via `stuntdouble run <agent>`:
- The agent process runs inside an isolated, unprivileged Docker container.
- Linux capabilities are dropped (`--cap-drop=ALL`).
- Workspace directories are mounted into the container to isolate filesystem access.

### 2. Egress Network Filtering & Mocks
- Host network endpoints and telemetry proxies log outgoing agent HTTP requests.
- Blocked external services (e.g. cloud APIs) return synthetic mock responses to allow agents to proceed without failing unexpectedly.

---

## Platform-Specific eBPF & Kernel Enforcement Roadmap

| Platform | Current Status | Architecture Target |
| :--- | :--- | :--- |
| **Linux** | Container capability dropping + network proxy | `cgroup_skb` egress filter attached via Aya / Cilium eBPF |
| **macOS** | Container isolation + network proxy | Endpoint Security Framework (ESF) socket filter |
| **Windows** | Container isolation + network proxy | Windows Filtering Platform (WFP) callout driver |

---

## Security Recommendations

1. **Least Privilege**: Always run StuntDouble CLI with non-root privileges on host environments.
2. **Telemetry Audit**: Enable the StuntDouble Control Plane to record agent API traffic for audit and compliance.
3. **Environment Separation**: Do not pass sensitive production tokens or AWS credentials to untrusted agent containers.
