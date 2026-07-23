# Stop Trusting AI Agents with Your Production Databases. Introducing StuntDouble.

*We gave LLMs terminal access. Then we realized they hallucinate. It’s time for zero-trust AI.*

If you’ve been paying attention to the AI engineering space, you’ve probably used an autonomous coding agent like Claude Code, Cursor, or OpenDevin. These tools are incredible—they can refactor massive codebases, debug complex architectures, and run shell commands on your behalf.

**But there is a terrifying flaw in the current ecosystem.**

When you give an AI agent terminal access to your local machine, it has *your* permissions. If it hallucinates a `DROP TABLE` command while trying to debug a local Postgres database, it will execute it. If it accidentally runs `rm -rf /` or decides to `curl` your AWS API keys to an untrusted server, your OS won't stop it.

We need to stop trusting agents blindly. We need a way to say: **"Write code for me, but don't blow up my computer."**

Enter **StuntDouble**.

## What is StuntDouble?

StuntDouble is a zero-trust, kernel-level execution sandbox specifically designed for autonomous AI agents. Instead of running an agent directly on your host machine, you run it through the StuntDouble CLI:

`stuntdouble run claude`

Behind the scenes, StuntDouble dynamically spins up an ephemeral, headless Docker container. It maps *only* your current working directory, locks down the CPU/RAM, and strips away all Linux root capabilities (`--cap-drop=ALL`).

But containerization alone isn't enough. AI agents still need to *believe* they are interacting with real databases and APIs to write tests and validate logic. 

## The Magic of Kernel-Level Egress Interception

This is where StuntDouble becomes magic. 

We utilize **eBPF on Linux**, **Endpoint Security (ESF) on macOS**, and **Windows Filtering Platform (WFP) on Windows** to physically intercept outbound TCP traffic at the kernel level. 

If your AI agent hallucinates and tries to aggressively migrate your local Postgres database (port 5432) or hit a live production Stripe API:
1. The StuntDouble Kernel Hook intercepts the packet before it touches the real network interface.
2. It pipes the packet to our high-performance **Wazero Plugin Engine**.
3. A community-built WebAssembly (`.wasm`) plugin reads the query and instantly generates a synthetic, mocked HTTP 200 success response.
4. The AI agent receives the mock response, believes it successfully modified the database, and continues its work—while your real database remains completely untouched.

## Time-Travel Built In

Even in a sandbox, an AI agent can mess up your local files. That's why StuntDouble takes a zero-copy git snapshot of your directory the exact millisecond before the agent executes. If the agent hallucinates and rewrites 50 files incorrectly, you simply run:

`stuntdouble rewind`

And your workspace is instantly restored.

## Deployed at Scale

StuntDouble isn't just a CLI. It's an enterprise ecosystem.
- **VS Code Extension**: Run your sandboxes directly from your IDE.
- **GitHub Actions**: Use our composite action to safely execute AI agents in your CI/CD pipelines to review PRs without exposing your build servers.
- **Kubernetes Native**: Deploy the StuntDouble DaemonSet via our official Helm Charts to wrap every AI workload in your enterprise cluster in a zero-trust eBPF shield.

## Try It Today

The era of trusting AI with raw shell access is over. Let your AI agents do their own stunts.

**Install via NPM:**
```bash
npm install -g stuntdouble-sandbox-cli
sd run claude
```

**[Check out the GitHub Repository Here]**
