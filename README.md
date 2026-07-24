<div align="center">
  <img src="https://raw.githubusercontent.com/itsrohan-lang/stuntdouble/main/docs/assets/logo.png" alt="StuntDouble Logo" width="200" height="200" />
  <h1>🛡️ StuntDouble</h1>
  <p><b>Zero-Trust eBPF Sandbox & Control Plane for Autonomous AI Agents</b></p>
  
  [![NPM Version](https://img.shields.io/npm/v/stuntdouble-sandbox-cli?color=00f0ff&style=for-the-badge)](https://www.npmjs.com/package/stuntdouble-sandbox-cli)
  [![License](https://img.shields.io/badge/License-MIT-8a2be2.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
</div>

<br/>

StuntDouble is a military-grade, zero-trust sandbox architecture designed explicitly for running AI coding agents (like Claude Engineer, OpenDevin, and Cursor) safely. It intercepts, audits, and optionally mocks system-level events (Network out, File writes, Exec) using native Linux **eBPF (Extended Berkeley Packet Filter)** before the AI can cause any harm.

---

## ✨ Features

* **🦀 Native eBPF Kernel Probes**: Written in Rust, it drops straight into your Linux Kernel to intercept rogue `cgroup_skb` network packets before they even reach the network interface.
* **🌐 CTO Control Plane**: A centralized Golang proxy that streams telemetry and enforces JSON-based global Enterprise RBAC policies across thousands of AI instances.
* **📊 Next.js SOC Dashboard**: Real-time Threat Vector analysis, live audit ledgers, and dynamic policy editors with a sleek hacker-aesthetic UI.
* **👻 Keploy API Mocking**: Instead of hard-crashing your AI agent when it tries to hit a blocked API (like Stripe or AWS), StuntDouble serves a mocked `200 OK` JSON ghost response so the agent thinks it succeeded and keeps coding!
* **💻 VS Code Extension**: Provides a seamless IDE integration. If the enterprise policy blocks your agent from exfiltrating data, you get a real-time VS Code alert natively in your editor.

---

## 🏗️ Architecture

The full architecture diagram and component breakdown can be found in [ARCHITECTURE.md](./ARCHITECTURE.md).

---

## 🚀 Quick Start

### 1. Install the CLI
Install the globally available StuntDouble CLI wrapper via NPM:
```bash
npm install -g stuntdouble-sandbox-cli
```

### 2. Start the Enterprise Control Plane
The Golang proxy handles all telemetry, database auditing, and policy enforcement.
```bash
git clone https://github.com/itsrohan-lang/stuntdouble.git
cd stuntdouble/control-plane
go run main.go
```
*(This automatically boots up the Rust eBPF engine in the background!)*

### 3. Launch the SOC Dashboard
Monitor your agents in real-time.
```bash
cd stuntdouble/dashboard
npm install
npm run dev
```
Open `http://localhost:3000` to access the CTO dashboard.

---

## 🛡️ Sandbox an Agent

To run an AI agent inside the zero-trust environment, simply prefix the command with `sd run`:

```bash
sd run claude-code
```

The StuntDouble CLI will automatically provision an ephemeral, locked-down Docker container, map the current working directory, and pipe standard I/O directly into the container while the eBPF kernel probes monitor all syscalls.

---

## 👻 Keploy Ghost Responses

If your global policy blocks access to `api.stripe.com`, your agent doesn't just crash. StuntDouble intercepts the `CONNECT` request and routes it to the `/api/keploy/mock` endpoint on the Control Plane, returning a synthetic response:

```json
{
  "status": "success",
  "mocked_by": "StuntDouble-Keploy-Integration",
  "data": {
    "status": "created",
    "amount": "0.00"
  }
}
```
The AI continues reasoning as if the API call succeeded!

---

## 📝 License
MIT License. Built for secure AI innovation.
