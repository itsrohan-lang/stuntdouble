# 📚 StuntDouble Documentation Site

This is the Next.js frontend repository for the official StuntDouble documentation site.

## 🚀 Getting Started

First, install dependencies and run the development server:

```bash
cd docs
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## 🛠️ Detailed Command Reference

The StuntDouble CLI (`sd`) provides extreme isolation for AI agents. Below is the comprehensive command manual.

### `sd init`
Initializes a new StuntDouble sandbox in the current directory.
- **Action**: Generates a `.stuntdouble.yaml` and `.stuntdouble.telemetry.json` file.
- **Usage**: Run this at the root of the codebase where the AI agent will operate.

### `sd run <agent>`
The primary orchestration command. Wraps the target AI agent inside an eBPF-secured Docker container (or remote cloud MicroVM).
- **Arguments**: `<agent>` (e.g. `claude`, `cursor`, `opendevin`, `bash`).
- **Options**:
  - `--remote, -r`: Offloads the sandbox to the StuntDouble Enterprise Cloud for execution.
  - `--env, -e`: Dynamically injects a specific base runtime image (default: `node:20-alpine`).
- **Example**: `sd run claude --env python:3.11-alpine`

### `sd daemon`
Starts the background Control Plane listener. Used heavily in CI/CD pipelines (like our GitHub Action) or Kubernetes Operators to enforce rules dynamically.
- **Options**:
  - `--mode`: Enforcement mode (`audit`, `block`, or `chaos`). Default is `block`.
  - `--policy`: Path to the `.stuntdouble.yaml` configuration.

### `sd chaos`
Activates Chaos Monkey Testing. This command actively injects simulated network drops, file-permission denials, and artificial latency to test how well the AI agent's error-recovery loop handles restricted environments.

### `sd protocol attest`
Performs a cryptographic attestation on the loaded sandbox kernel modules (via Sigstore) to guarantee that they have not been maliciously tampered with prior to agent execution.

## 🌐 Enterprise Integrations

- **Python SDK**: Import `stuntbot` in your LangChain workflows to spawn secure containers directly via Python.
- **Kubernetes Operator**: Apply `StuntDoublePolicy` CRDs directly to your cluster to enforce global network policies.
- **Terraform Provider**: Manage API keys and global RBAC policies on the Control Plane using standard IaC.

---
*Generated for the StuntDouble ecosystem.*
