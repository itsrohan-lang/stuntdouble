# 🛡️ StuntDouble

> **The Ultimate Zero-Trust, eBPF-Powered Enterprise Sandbox for Autonomous AI Agents.**

[![Build Status](https://github.com/itsrohan-lang/stuntdouble/actions/workflows/release.yml/badge.svg)](https://github.com/itsrohan-lang/stuntdouble/actions)
[![Version](https://img.shields.io/npm/v/stuntdouble-sandbox-cli.svg)](https://npmjs.org/package/stuntdouble-sandbox-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

StuntDouble is a hyper-secure, enterprise-grade execution environment that allows AI coding agents (like Claude, Cursor, OpenDevin, and SWE-agent) to write and execute code on your local machine or cloud infrastructure *without* the risk of destroying databases, wiping hard drives, or exfiltrating API keys.

By leveraging **eBPF** (Linux), **Endpoint Security Framework** (macOS), and **WebAssembly (WASM)**, StuntDouble strictly enforces granular network and filesystem policies dynamically.

---

## 🌟 The StuntDouble Ecosystem

Over the course of its development, StuntDouble has evolved from a simple CLI into a massive global ecosystem designed for DevOps, Security, and Engineering teams:

- 💻 **StuntDouble CLI**: The native command-line interface to spin up instant, isolated AI workspaces locally.
- 🏢 **Control Plane (Go)**: A centralized REST/GraphQL API for CTOs and Security Admins to enforce global Role-Based Access Control (RBAC) policies across all sandboxes.
- 🐍 **Python SDK (`stunt-python`)**: Programmatically spawn StuntDouble sandboxes directly from Python scripts (perfect for LangChain and LlamaIndex integrations).
- 🐙 **StuntBot (GitHub App)**: Automatically runs AI agents on Pull Requests inside secure, ephemeral microVMs.
- ☸️ **Kubernetes Operator (CRDs)**: Native `StuntDoublePolicy` Kubernetes resources that sync dynamically to your cluster.
- ☁️ **Terraform Provider**: Manage your StuntDouble Control Plane policies globally using Infrastructure-as-Code (IaC).
- 🚀 **GitHub Action**: Wrap your entire CI/CD pipeline inside the StuntDouble eBPF sandbox with a single YAML step (`itsrohan-lang/stuntdouble-action`).
- 🖥️ **Desktop App**: A beautiful Electron/React GUI to manage your local sandboxes without touching a terminal.
- 📊 **Prometheus/Grafana**: Built-in metrics exporter (`/metrics`) to observe global sandbox behaviors and blocked network requests.
- 🌐 **WebAssembly Engine**: Run StuntDouble policies directly inside the browser or on Cloudflare Workers edge runtimes.
- 🐧 **StuntOS**: A custom Linux distribution built from scratch via Buildroot, hyper-optimized specifically to run StuntDouble containers at bare-metal speeds.

---

## 🚀 Quick Start

### 1. Install the CLI
```bash
npm install -g stuntdouble-sandbox-cli
# or via go:
go install github.com/itsrohan-lang/stuntdouble/cli@latest
```

### 2. Initialize your project
```bash
sd init
```
This generates a `.stuntdouble.yaml` policy file in your repository where you can define network rules, allowed agents, and file access paths.

### 3. Run an AI agent securely
```bash
sd run claude
```
The agent is now securely jailed inside the StuntDouble environment!

---

## 🛠️ Detailed CLI Commands

The `sd` (StuntDouble) CLI is the core orchestration tool for managing your sandboxes.

| Command | Description | Flags / Options |
|---------|-------------|-----------------|
| `sd init` | Initializes a new StuntDouble project by generating the default `.stuntdouble.yaml` policy file and telemetry state. | |
| `sd run <agent>` | Spawns a highly restricted Docker container and executes the specified AI agent (e.g., `claude`, `bash`) wrapped in eBPF hooks. | `--remote` (Runs in cloud microVM), `--env <image>` (Specify custom Docker runtime image) |
| `sd daemon` | Runs the StuntDouble background daemon. Used primarily by the GitHub Action and Kubernetes Operator to listen for policy updates. | `--mode <audit\|block\|chaos>`, `--policy <file>` |
| `sd chaos` | Activates Chaos Monkey Testing. Actively injects simulated network failures and file permission errors to test how resilient your AI agent's error-recovery logic is. | |
| `sd protocol attest` | Triggers a cryptographic attestation process (via Sigstore/Cosign) to guarantee the loaded sandbox kernel modules have not been tampered with. | |

### Advanced Usage Examples

**Run Claude inside a Python 3.11 environment with Cloud Sync enabled:**
```bash
sd run claude --env python:3.11-alpine --remote
```

**Start the Daemon in Audit-only mode for CI/CD:**
```bash
sudo sd daemon --mode audit --policy .stuntdouble.yaml
```

---

## 🏗️ Architecture Deep Dive

Please see [ARCHITECTURE.md](ARCHITECTURE.md) for detailed system design, component breakdowns, and kernel-level interception diagrams.

### How it works:
1. **Snapshot**: When you run `sd run`, StuntDouble takes a zero-copy git snapshot of your working directory.
2. **Jail**: It spawns an isolated Docker container with zero egress network capability (by default).
3. **Hook**: eBPF/ESF hooks are injected directly into the kernel cgroups to monitor read/write/execute syscalls.
4. **Mock**: External API calls made by the AI (like AWS or Stripe) are routed to WebAssembly plugins that return safe, mock JSON data.
5. **Revert**: If the agent hallucinates and destroys the codebase, the user can instantly time-travel back to the pre-run snapshot.

---

## 🤝 Contributing & Funding

StuntDouble is an open-source project designed to protect the future of AI engineering. 

To support the development, please check out the **Sponsor** button at the top of the repository, or visit my Ko-fi directly:
💖 [Support me on Ko-fi](https://ko-fi.com/rdx463)

---
*Built with absolute paranoia by the StuntDouble Core Team.*
