# 🛡️ StuntDouble

> **Zero-Trust, eBPF-powered Sandbox for Autonomous AI Agents.**

[![Build Status](https://github.com/itsrohan-lang/stuntdouble/actions/workflows/release.yml/badge.svg)](https://github.com/itsrohan-lang/stuntdouble/actions)
[![Version](https://img.shields.io/npm/v/stuntdouble-sandbox-cli.svg)](https://npmjs.org/package/stuntdouble-sandbox-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

StuntDouble is a hyper-secure execution environment that lets AI coding agents (like Claude, Cursor, and OpenDevin) write and execute code on your local machine without the risk of destroying your databases, wiping your hard drive, or exfiltrating your API keys.

## 🧠 Architecture

Please see [ARCHITECTURE.md](ARCHITECTURE.md) for detailed system design, component breakdowns, and kernel-level interception diagrams (eBPF, macOS ESF, Windows WFP).

## 🚀 Quick Start

**1. Install the CLI:**
```bash
npm install -g stuntdouble-sandbox-cli
# or via go:
go install github.com/stuntdouble/cli@latest
```

**2. Initialize your project:**
```bash
sd init
```

**3. Run an AI agent safely:**
```bash
sd run claude
```

## 🌟 Key Features

* **eBPF Kernel Level Interception:** Physically drops malicious TCP SYN packets directed at local Postgres/Mongo instances.
* **WASM Plugin Engine:** Community-driven WebAssembly plugins dynamically mock AWS, Stripe, and internal API responses.
* **Instant Time Travel:** Automatic zero-copy git snapshots allow you to revert your entire workspace instantly if the agent hallucinates.
* **Chaos Monkey Testing:** Run `sd chaos` to actively sabotage the sandbox and verify your LLM's resilience and error-recovery logic.

## 📦 Extensibility
StuntDouble ships natively as a CLI, a GitHub Action, and a VS Code Extension. It's ready to be embedded anywhere AI agents write code.

---
*Built with absolute paranoia by the StuntDouble Core Team.*
