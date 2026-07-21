# 🎬 StuntDouble

<p align="center">
  <em>Child-proof your AI Agents. Let them run YOLO mode, safely.</em>
</p>

<p align="center">
  <img src="https://img.shields.io/npm/v/stuntdouble-sandbox-cli?color=00f0ff&label=npm" alt="NPM Version">
  <img src="https://img.shields.io/badge/status-production-8a2be2?style=flat-square" alt="Status">
  <img src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" alt="License">
</p>

## Overview

**StuntDouble** is a 1-click execution sandbox and network mocker for AI coding agents. 

Autonomous agents (like Claude Code, Cursor, and Opencode) are incredibly powerful, but running them locally comes with extreme anxiety. One bad prompt can lead to deleted databases, leaked AWS keys, or borked system configurations. 

StuntDouble solves this by acting as the ultimate zero-trust orchestration glue. With a single command, it wraps your agent in a locked-down Docker MicroVM, intercepts destructive network calls via eBPF mocks, and governs all execution using the Universal Stunt Protocol (STP).

## 🚀 Key Features

* **1-Click Safe Mode:** Stop writing manual `docker-compose` wrappers. Just run `stuntdouble run claude`.
* **eBPF Interceptors:** Intercepts destructive database/network calls and safely mocks them so the agent *thinks* it succeeded, while your host system remains untouched.
* **StuntNet Swarms:** Orchestrate entire teams of agents (`sd swarm qa-agent dev-agent`) inside a virtual intranet isolated from the real web.
* **Time-Travel Rewind:** If an agent deletes your workspace, instantly rewind the ZFS snapshot to 5 minutes ago (`sd rewind 5`).
* **The Warden:** An autonomous AI defender that monitors network traffic and generates zero-day eBPF patches on the fly to prevent agent escapes (`sd warden`).
* **Universal Governance (STP):** Provides an HTTP server and cryptographic sandbox attestations so compliant foundational LLMs refuse to execute unless isolated (`sd protocol`).

## 🛠 Installation

StuntDouble is officially published to the global NPM registry.

```bash
npm install -g stuntdouble-sandbox-cli
# or run without installing:
npx stuntdouble-sandbox-cli init
```

## 🎮 Quick Start

1. **Initialize the sandbox** in your project (generates `.stuntdouble.yaml` and safety rules):
```bash
stuntdouble init
```

2. **Run your agent securely:**
```bash
stuntdouble run claude
```

3. **Check safety telemetry:**
```bash
stuntdouble stats
```

## 📚 Documentation & Next.js Site

We have a beautiful documentation landing page built with Next.js and Tailwind CSS.
To run the documentation site locally:
```bash
cd docs
npm install
npm run dev
```
Then navigate to [http://localhost:3000](http://localhost:3000).

For a deeper dive into how StuntDouble works under the hood:
* [Architecture Blueprint](./ARCHITECTURE.md)
* [Development Plan & Phases](./PLAN.md)
* [Contributing Guidelines](./CONTRIBUTING.md)

## License
MIT License
