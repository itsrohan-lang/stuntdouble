# 🤝 Contributing to StuntDouble

First off, thank you for considering contributing to StuntDouble! We want to make local AI agents safe for everyone.

## Core Philosophy
1. **Developer Experience First:** StuntDouble must feel invisible. If it adds 10 seconds of latency to an agent run, developers won't use it.
2. **Zero Trust:** Never assume an AI agent is benign. Always default to the strictest isolation possible.

## Setting up your Local Environment

1. **Clone the repo:**
   ```bash
   git clone https://github.com/stuntdouble/stuntdouble.git
   cd stuntdouble
   ```

2. **Prerequisites:**
   * Docker (must be running locally)
   * Go 1.21+ (or Node.js 20+, depending on final stack choice)
   * `make`

3. **Build the CLI:**
   ```bash
   make build
   ```

## Development Workflow

* **Issue First:** Please open an issue to discuss significant architectural changes before submitting a PR.
* **Testing:** All new features must include integration tests verifying that sandbox escapes are impossible. Use `make test`.
* **Commit Standards:** We use Conventional Commits (e.g., `feat: added keploy injection`, `fix: patched volume mount escape`).

## Roadmap Contributions
We are currently focusing entirely on **Phase 1** (see `PLAN.md`). If you're looking for a place to start, checking the open issues labeled `good first issue` in the repository is the best path forward.

---
*Note: Because StuntDouble deals with security and host-isolation, all PRs modifying the container generation logic will undergo strict security review.*
