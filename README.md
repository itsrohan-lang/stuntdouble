# 🎬 StuntDouble

<p align="center">
  <em>Child-proof your AI Agents. Let them run YOLO mode, safely.</em>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/status-in%20development-orange?style=flat-square" alt="Status">
  <img src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" alt="License">
</p>

## Overview

**StuntDouble** is a 1-click execution sandbox and network mocker for AI coding agents. 

Autonomous agents (like Claude Code, Cursor, and Opencode) are incredibly powerful, but running them locally comes with extreme anxiety. One bad prompt can lead to deleted databases, leaked AWS keys, or borked system configurations. 

StuntDouble solves this by acting as the ultimate DX (Developer Experience) orchestration glue. With a single command, it:
1. Wraps your agent in an ephemeral, locked-down Docker microVM.
2. Injects a powerful eBPF mocking layer (via Keploy).
3. Intercepts destructive database/network calls and safely mocks them so the agent *thinks* it succeeded, while your host system remains untouched.

## 🚀 Features

* **1-Click Safe Mode:** Stop writing manual `docker-compose` wrappers. Just run `sd claude`.
* **"Panic Mode" Network Blocking:** Automatically block outbound calls to common local databases (Postgres, Mongo) unless explicitly allowed.
* **The "Stunt" Layer:** Seamlessly record real traffic and mock database responses so agents can run integration tests without destroying data.

## 🛠 Installation

*(Coming Soon - Project is in Phase 1 Planning)*
```bash
npm install -g stuntdouble
# or
brew install stuntdouble
```

## 🎮 Usage

Initialize a safe environment in your project:
```bash
stuntdouble init
```
This generates a `.stuntdouble.yaml` configuration file outlining network policies.

Run your agent securely:
```bash
stuntdouble run claude
```
*(Your agent is now running inside a strictly isolated microVM, with network calls mocked!)*

## 📚 Documentation

For a deeper dive into how StuntDouble works under the hood, please refer to the following documents:
* [Architecture Blueprint](./ARCHITECTURE.md)
* [Development Plan & Phases](./PLAN.md)
* [Contributing Guidelines](./CONTRIBUTING.md)

## License
MIT License
