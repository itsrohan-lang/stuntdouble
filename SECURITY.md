# Security Policy

## Supported Versions

Currently, StuntDouble is in active beta and only the latest major version is supported for security updates.

| Version | Supported          |
| ------- | ------------------ |
| 3.x.x   | :white_check_mark: |
| < 3.0.0 | :x:                |

## Reporting a Vulnerability

StuntDouble is a security product, meaning we take vulnerability reports extremely seriously.

If you discover a vulnerability (such as a sandbox escape, eBPF bypass, or Control Plane exploit), **please do not open a public issue.** 

Instead, please email us directly at **security@stuntdouble.io** or message us securely.

We will acknowledge your email within 24 hours, and work to issue a patch and a CVE advisory as quickly as possible. 

### What to include in your report
- The agent used (e.g. Claude Code, Aider, OpenDevin)
- The exact OS and kernel version
- A proof-of-concept (PoC) bash script or repo that demonstrates the sandbox escape
- Logs from the Control Plane (if applicable)

We deeply appreciate the community's help in keeping StuntDouble the absolute safest way to run AI code!
