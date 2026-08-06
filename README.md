<div align="center">
  <img src="https://raw.githubusercontent.com/itsrohan-lang/stuntdouble/main/docs/assets/logo.png" alt="StuntDouble Logo" width="200" height="200" />
  <h1>🛡️ StuntDouble</h1>
  <p><b>Run AI coding agents in a restricted container, and roll back what they change</b></p>

  [![NPM Version](https://img.shields.io/npm/v/stuntdouble-sandbox-cli?color=00f0ff&style=for-the-badge)](https://www.npmjs.com/package/stuntdouble-sandbox-cli)
  [![License](https://img.shields.io/badge/License-MIT-8a2be2.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
</div>

<br/>

StuntDouble runs autonomous coding agents (Claude Code, Cursor, OpenDevin, aider) inside
a Docker container with dropped capabilities and resource limits, snapshots your workspace
so you can undo what the agent did, and records what happened.

## ⚠️ What this does and does not protect against

Read this before relying on StuntDouble for anything.

**What works today:**

| Control | Mechanism |
| --- | --- |
| No new privileges in the sandbox | `--cap-drop=ALL` on the agent container |
| Bounded CPU, memory & time | `--memory=2g --cpus=1.0 --max-duration=15m` |
| Filesystem scope | only the working directory is bind-mounted at `/workspace` |
| Undo agent changes | workspace snapshot via git plumbing; `sd rewind` and `sd checkpoint` |
| Zero-Trust API Credentials | dummy key substitution proxy (`ZeroTrustProxy`) on host egress |
| Indirect Prompt Injection Guard | adversarial prompt vector detector & sanitizer via `sd guard prompt` |
| Enterprise AI Data Loss Prevention | PII & secret inspection scanner via `sd dlp scan` |
| Synthetic Offline Mocking | auto-generate synthetic DB & API mocks via `sd mock` |
| CI/CD Security & Red-Teaming | `sd verify` security gate & `sd swarm --redteam` adversarial node |
| Run history | local JSON counter, control plane + dashboard |

**What does not work:**

- **Network egress filtering is not implemented on any platform.** `cli/pkg/ebpf.AttachInterceptor` returns `ErrUnsupported` on Linux, macOS and Windows alike. The BPF program in `cli/pkg/ebpf/bpf_prog.c` is marked `// +build ignore` and no build step compiles it. `sd run` requires `--allow-unenforced-network` to proceed, which is how you acknowledge this.
- **There are no native OS interceptors.** The macOS EndpointSecurity (`mac/`) and Windows WFP (`windows/`) backends were claimed by earlier versions but have never been built, signed or shipped. Both loaders return `ErrUnsupported`.
- **Policies are advisory everywhere.** `blocked_ports`, `strict_egress` and allow/deny lists are parsed, distributed and displayed — but no code path applies them to a connection.
- **The audit ledger is only tamper-evident against direct database edits.** `/api/audit/verify` checks an HMAC-SHA256 chain, but the signing key is `STUNTDOUBLE_TOKEN` — the same bearer token every CLI client must hold to write logs. Any authorized client can therefore forge a valid chain. It detects someone editing `stuntdouble_audit.db` behind the server's back; it does not detect a client that lies.
- **`sd record` runs Keploy privileged** with `--pid=host` and `--net=host`. That is a broad grant on your host.

See [docs/ENFORCEMENT.md](./docs/ENFORCEMENT.md) for the full boundary and the plan to close it.

## 🚀 Quick start

### 1. Install the CLI

```bash
curl -sSL https://raw.githubusercontent.com/itsrohan-lang/stuntdouble/main/install.sh | bash
```

The installer verifies the binary against the `SHA256SUMS` published with the release and
aborts if it cannot. If you would rather not pipe a script to a shell, download
`install.sh`, read it, and run it yourself — or install from npm:

```bash
npm install -g stuntdouble-sandbox-cli
```

### Uninstallation

To cleanly remove StuntDouble CLI binaries and global npm packages from your system:

```bash
curl -sSL https://raw.githubusercontent.com/itsrohan-lang/stuntdouble/main/uninstall.sh | bash
```

### 2. Sandbox an agent

```bash
sd run claude --allow-unenforced-network
```

This pulls the runtime image, starts a Keploy sidecar, snapshots the workspace, and runs
the agent in a container with `--cap-drop=ALL` and the current directory mounted at
`/workspace`. Pick a different base image with `-e python:3.11-alpine`.

`--allow-unenforced-network` is required and is not a formality: without it there is no
egress filtering, and the flag is how you confirm you know that.

### 3. Undo what it did

```bash
sd rewind
```

Restores tracked files to the pre-run snapshot and removes files created since. **This
discards uncommitted work**, including your own changes made after the run started.

## 🏢 Optional: control plane and dashboard

Only needed if you want run telemetry aggregated across machines.

```bash
export STUNTDOUBLE_TOKEN=$(openssl rand -hex 32)
cd control-plane && go run .
```

The control plane binds to `127.0.0.1:4439` and requires this bearer token on every
endpoint except `/api/health`. It refuses to start without one.

```bash
cd dashboard
npm install
NEXT_PUBLIC_STUNTDOUBLE_TOKEN=$STUNTDOUBLE_TOKEN npm run dev
```

Open `http://localhost:3000`. Every number shown is read from the control plane; when it is
unreachable, cards show `—` rather than a placeholder value.

## 👻 Keploy mocks

`sd record <command>` runs your command under Keploy to capture outbound calls as
replayable mocks, so an agent can exercise code paths that hit a database without touching
a real one. The control plane also exposes `/api/keploy/mock`, which returns a canned
success payload for a caller that has been routed to it.

To be precise about the mechanism: Keploy does the capture and replay. StuntDouble does not
intercept a connection and transparently redirect it to a mock — that would require the
egress filtering described above.

## 📂 Repository layout

| Path | Status |
| --- | --- |
| `cli/` | Go CLI — working core with ZeroTrustProxy, checkpointing & guardrails |
| `control-plane/` | Go telemetry + policy service |
| `dashboard/`, `docs/` | Next.js dashboard and site |
| `cli/pkg/ebpf/` | Stub — every backend returns `ErrUnsupported`; the BPF program is not compiled |
| `core-ebpf/` | Rust eBPF engine placeholder — never built or run |
| `charts/`, `k8s-operator/`, `terraform-provider-stuntdouble/` | Helm chart, CRD reconciler & Terraform provider — **the chart grants host privileges for the unimplemented eBPF engine; read `charts/stuntdouble/values.yaml` before deploying** |
| `mac/`, `windows/` | Unbuilt sources for a macOS EndpointSecurity interceptor & Windows WFP driver |
| `stuntos/`, `wasm/` | WASM policy evaluator & exploratory stubs |

## 🤝 Contributing

The highest-value contribution is the one thing that would make the premise true: a working
`cgroup_skb/egress` filter behind `cli/pkg/ebpf.AttachInterceptor`. The BPF C program in
`cli/pkg/ebpf/bpf_prog.c` is a reasonable starting point but is not currently compiled.

## 📝 License

MIT.
