# 🗺 Roadmap

Status of the project, honestly stated. A box is checked only when the feature works
end to end.

Every item below was previously marked complete, including a custom hypervisor, physical
hardware clusters, and a protocol built into foundational model weights. None of that was
built. The checkboxes have been reset to reflect the code.

---

## Phase 1: Container isolation

**Goal:** one command to run an agent in a restricted container.

* [x] CLI project structure (Go + cobra)
* [x] `sd init` generates `.stuntdouble.yaml` in the format the parser actually accepts
* [x] `sd run <agent>` spawns the agent and pipes I/O
* [x] `--cap-drop=ALL`, memory and CPU limits, scoped workspace mount
* [ ] **Network egress filtering.** Not implemented on any platform. This is the gap that
      matters most: it is the project's premise. `cli/pkg/ebpf.AttachInterceptor` returns
      `ErrUnsupported` and `sd run` requires `--allow-unenforced-network` to proceed.

## Phase 2: Mocking

**Goal:** let agents exercise code paths that hit external services, safely.

* [x] Keploy sidecar attached to the agent's network namespace
* [x] `sd record` to capture mocks
* [x] `/api/keploy/mock` returns a canned success payload
* [ ] Transparent interception — answering a blocked connection with a mock instead of
      failing. Blocked on egress filtering.

## Phase 3: Telemetry & UX

* [x] Local run counter and `sd stats`
* [x] `sd monitor` terminal view of the control plane
* [x] Dashboard fed entirely by real API responses
* [ ] Report what was blocked during a session. Nothing is blocked, so there is nothing to
      report.
* [ ] Sub-second container spin-up. Never measured; startup is dominated by image pull.

## Phase 4: Control plane

* [x] Go service aggregating run counts, with an SQLite audit log
* [x] Bearer-token auth on all endpoints; loopback bind; single-origin CORS
* [x] Policy document served over REST and GraphQL
* [ ] Policy **enforcement**. Policies are distributed and displayed, not applied.
* [ ] Tamper-evident audit log. Records are self-reported by clients today.
* [ ] Hosted/remote execution. There is no StuntDouble Cloud; the `--remote` flag that
      printed cloud provisioning messages without doing anything has been removed.

## Phase 5: Ecosystem

* [x] CI workflow generator (`sd ci`)
* [ ] IDE extensions. `vscode-extension/` exists but is unproven and unpublished.
* [ ] Plugin registry. The `sd install` command wrote an 8-byte stub file and claimed to
      have installed an interceptor; it has been removed.

## Phase 6: Multi-agent

* [x] `sd swarm` spawns multiple agent containers on a shared Docker network
* [ ] Isolated inter-agent network with synthetic external services

## Not planned

Removed from the roadmap. These were listed as complete; none exist, and none are a
sensible next step for a project whose core enforcement mechanism is unbuilt:

- **Custom hypervisor / MicroVM ("StuntOS").** The `sd os` command printed
  "your hardware is cryptographically partitioned" and slept. Removed.
- **"Stunt Boxes" physical hardware clusters.** No hardware exists.
- **"Stunt Protocol" (STP) adopted by foundational models.** The `sd protocol` command
  generated an HMAC token that nothing consumed, and claimed compliant models would refuse
  to run outside a sandbox. No model implements this. Removed.
- **Partnering with OpenAI / Anthropic / Google to build STP into model weights.** No such
  partnership was ever sought or exists.
- **Autonomous "Warden Agent" self-patching against zero-days.** `sd warden` fabricated a
  zero-day alert, slept two seconds, and appended a hardcoded port-6379 block to
  `bpf_prog.c` — a file that is never compiled. Repeated runs duplicated the block
  indefinitely. Removed, and the committed patch reverted.

---

## Next step

One thing: implement `cgroup_skb/egress` filtering on Linux and wire it to the policy
document. Everything else here is scaffolding around a control that does not exist yet.
