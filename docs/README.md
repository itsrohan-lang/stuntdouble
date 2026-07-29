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

## 🛠️ Command Reference

The StuntDouble CLI (`sd`) runs AI agents in a restricted Docker container.

> Network egress filtering is **not implemented**. The sandbox provides container
> isolation only. See [ENFORCEMENT.md](./ENFORCEMENT.md) for the full boundary.

### `sd init`
Writes a `.stuntdouble.yaml` config file in the current directory.
- **Usage**: Run this at the root of the codebase where the agent will operate.

### `sd run <agent> [args...]`
The primary command. Runs the agent inside a Docker container with
`--cap-drop=ALL`, a 2g memory cap, a 1.0 CPU cap, and only the working directory
mounted at `/workspace`. Snapshots the workspace first so `sd rewind` can undo it.
- **Arguments**: `<agent>` — `claude`, `sh`, `bash`, or any npm package name.
  Anything after the agent name is forwarded to it verbatim.
- **Options**:
  - `--allow-unenforced-network`: Required. Acknowledges that egress filtering is
    not active. `sd run` refuses to start without it.
  - `--env, -e`: Base runtime image (default: `node:20-alpine`).
- **Example**: `sd run --allow-unenforced-network --env python:3.11-alpine claude`

### `sd rewind`
Restores the workspace to the snapshot taken at the last `sd run`.

### `sd record <agent>`
Records database and API traffic with a Keploy sidecar to generate mocks. Keploy
runs privileged, with `--pid=host` and `--net=host`.

### `sd stats`
Shows the local run count from `.stuntdouble.telemetry.json`.

### `sd serve` / `sd monitor`
`serve` starts the local telemetry API; `monitor` is a terminal view of the
control plane.

### `sd swarm`
Spawns multiple agent containers on a shared Docker network.

### `sd ci`
Generates a GitHub Actions workflow that runs agents in the sandbox.

## 🌐 Other components in this repo

These exist but are unproven — see [PLAN.md](../PLAN.md) for status:

- **Python SDK** (`python-sdk/`): `Sandbox` context manager that shells out to `sd run`.
- **Kubernetes Operator** (`k8s-operator/`): `StuntDoublePolicy` CRD. Policies are
  distributed and displayed, not enforced.
- **Terraform Provider** (`terraform-provider-stuntdouble/`): manages control-plane
  policy documents.

---
*Generated for the StuntDouble ecosystem.*
