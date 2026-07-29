# Enforcement and the security boundary

This document states exactly what StuntDouble enforces, what it does not, and what it would
take to close the gap. It exists because the project's code, README and dashboards
previously claimed a network enforcement layer that was never built.

Referenced from `cli/pkg/ebpf.ErrUnsupported`, `core-ebpf/src/main.rs`, `sd --help`,
`sd init`, [README.md](../README.md) and [ARCHITECTURE.md](../ARCHITECTURE.md).

---

## Summary

**Network egress filtering is not implemented on any platform.** An agent running under
`sd run` has the same outbound network access as the host. `sd run` will not start without
`--allow-unenforced-network`, which is how you acknowledge that.

Everything StuntDouble does enforce is a Docker or git mechanism, not a StuntDouble one.

A previous version of this document said blocked external services "return synthetic mock
responses to allow agents to proceed". Nothing is blocked, and no such interception exists.

## What is enforced

These are real, and they are the reason to use the tool.

| Control | Mechanism | Enforced by |
| --- | --- | --- |
| No capabilities in the sandbox | `--cap-drop=ALL` | Linux kernel, via Docker |
| Memory ceiling | `--memory=2g` | cgroup v2 memory controller |
| CPU ceiling | `--cpus=1.0` | cgroup v2 CPU controller |
| Filesystem scope | only `$PWD` bind-mounted at `/workspace` | mount namespace |
| Reversible changes | pre-run tree snapshot via git plumbing | `sd rewind` |

Consequences worth being concrete about:

- A `rm -rf` inside `/workspace` is recoverable with `sd rewind`. A destructive command
  elsewhere in the container destroys the container, not your machine — but `/workspace` is
  your actual working directory, so the snapshot is the thing that saves you.
- An agent that forks endlessly or allocates without bound is capped, not stopped. The
  container gets OOM-killed rather than taking the host down.
- The snapshot is a git tree object in your own repository. It does not cover files outside
  the repo, and it is not a substitute for committing.
- `sd rewind` runs `git clean -fd`, so it also deletes untracked files you created yourself
  after the run started.

## What is not enforced

### Network egress

Nothing filters it. Specifically, an agent in the sandbox can reach:

- your local databases on the host, via the Docker bridge gateway
- cloud instance metadata endpoints (`169.254.169.254`) if the host has them
- any paid API, including with the keys described below
- the public internet, for exfiltration or for pulling code to execute

The Keploy sidecar shares the agent's network namespace, but it is a capture tool. It
observes; it does not deny.

### Policy

`.stuntdouble.yaml` and the control plane's policy document describe intent.
`enforcement_mode`, allow lists and deny lists are parsed, distributed over REST/GraphQL, and
displayed in the dashboard. **No code path applies them to a connection.** Treat the policy
file as documentation of what you would like to be true.

### Audit integrity

The control plane's SQLite audit log contains what CLI instances POST to it. There is no
signing, no append-only guarantee, and no server-side corroboration. A client that lies is
believed. It is useful for "what did I run last Tuesday", not for "prove the agent did not
do X".

### Credentials

`ANTHROPIC_API_KEY` and `OPENAI_API_KEY` are forwarded into the agent container, because the
agent cannot call its model without them. Anything executing in that container can read them
from its own environment, and — see above — can send them anywhere.

Mitigation is operational, not technical: use scoped keys, set spend limits, rotate them.

### `sd record`

Keploy's capture mode needs `--pid=host` and `--net=host` and runs with `NET_ADMIN`. That is
a broad grant on your machine. It is opt-in and separate from `sd run`; do not treat an
`sd record` session as sandboxed.

### The Helm chart

`charts/stuntdouble/templates/daemonset.yaml` requests `privileged: true`, `hostPID`,
`hostNetwork`, a writable `/sys/fs/cgroup`, and a `docker.sock` hostPath mount. Those
permissions were provisioned for the eBPF engine that does not exist, so today they buy
nothing and grant node-root. **Do not deploy the chart as-is.**

## Threat model

| Adversary | Contained? |
| --- | --- |
| Agent that makes a mistake — deletes files, botches a refactor, loops | Yes — `sd rewind`, resource limits |
| Agent that runs an unexpected build step or installs a package | Partly — confined to the container filesystem, but it can fetch anything |
| Agent that reads your source and sends it somewhere | **No** |
| Agent that uses your API keys for its own purposes | **No** |
| Agent that connects to production because it found the URL in `.env` | **No** |
| Deliberately hostile or prompt-injected agent aiming at the host | **No** — assume container escape is in scope for a determined attacker |

StuntDouble reduces the cost of agent *error*. It is not a boundary against agent *intent*.

## Closing the gap

The intended design is a `cgroup_skb/egress` BPF program attached to the agent container's
cgroup, dropping packets whose destination is not permitted by policy.

`cli/pkg/ebpf/bpf_prog.c` is a plausible starting point: it parses the IPv4 header, reads the
TCP destination port with `bpf_skb_load_bytes`, looks it up in a `BPF_MAP_TYPE_HASH` named
`blocked_ports`, and returns `0` (drop) on a hit. It is marked `// +build ignore` and no build
step invokes `bpf2go`, so it has never been compiled or run. It is untested code, not a
working program.

Work required, in order:

1. **Generate bindings.** Wire the `//go:generate ... bpf2go` directive in `loader_linux.go`
   into the build, with clang and kernel headers available. Vendor the generated object so
   users do not need a toolchain.
2. **Load and attach.** Replace the `ErrUnsupported` stub with `LoadBpfObjects` +
   `link.AttachCgroup`, targeting **the container's** cgroup path. Not `/sys/fs/cgroup/` —
   that is the root cgroup and would filter the entire host. Resolve the path from the
   container ID.
3. **Populate from policy.** Write the parsed deny list into `blocked_ports` at attach time,
   and update the map when policy changes.
4. **Read back drops.** Expose a real blocked-connection count via a counter map, so telemetry
   can report what was actually denied. Every "blocked" figure this project ever displayed was
   fabricated; do not re-introduce one that is not read from a map.
5. **Fail closed.** If the program cannot attach, `sd run` must refuse. That is the current
   behaviour and should be preserved.
6. **Handle the escape hatches.** Destination-port filtering is bypassed by a proxy on an
   allowed port, and says nothing about DNS or IPv6. Decide whether the goal is default-deny
   by address, and document what remains bypassable.

macOS and Windows have no equivalent. `cgroup_skb` is Linux-only; under Docker Desktop the
agent container runs inside a Linux VM, so a Linux implementation would apply to the
container but would not filter anything on the host side. Native macOS EndpointSecurity and
Windows Filtering Platform backends were previously *claimed* — the loaders shelled out to
`mac/src/esf_interceptor` and a `StuntDoubleWFP` service, and neither has ever existed. Both
now return `ErrUnsupported`.

Until step 5 is real, the honest description of this project is the one in the README: a
resource-limited, capability-dropped, snapshot-backed container runner.

## Security recommendations

1. Run the CLI as a non-root user on the host.
2. Use scoped, revocable, spend-limited API keys for any agent you sandbox.
3. Do not pass production credentials or cloud tokens into an agent container.
4. Commit before `sd run`. The snapshot is a convenience; git is the durable record.
5. Do not rely on the audit log for compliance evidence.
