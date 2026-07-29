import subprocess


class Sandbox:
    """Thin wrapper around the ``sd`` CLI.

    This does nothing to the calling Python process. It shells out to
    ``sd run``, which starts the agent command in a Docker container with
    dropped capabilities, memory and CPU limits, and only the working
    directory mounted.

    Network egress filtering is not implemented, so the sandboxed command can
    still reach anything the host can reach. ``sd run`` refuses to start
    unless that is acknowledged, which is why ``allow_unenforced_network``
    has no default -- see docs/ENFORCEMENT.md.
    """

    def __init__(self, allow_unenforced_network: bool, env: str = "node:20-alpine"):
        if not allow_unenforced_network:
            raise ValueError(
                "Sandbox requires allow_unenforced_network=True. Network egress "
                "filtering is not implemented, so a command run through this SDK "
                "has unrestricted outbound network access. Pass True to "
                "acknowledge that, or do not use the sandbox."
            )
        self.env = env

    def __enter__(self):
        return self

    def run(self, agent_command: str) -> subprocess.CompletedProcess:
        """Run ``agent_command`` inside the container.

        Returns the CompletedProcess. Raises CalledProcessError if the agent
        exits non-zero, so failures are not swallowed.
        """
        return subprocess.run(
            [
                "sd",
                "run",
                "--allow-unenforced-network",
                "--env",
                self.env,
                "sh",
                "-c",
                agent_command,
            ],
            check=True,
            text=True,
        )

    def __exit__(self, exc_type, exc_value, traceback):
        return False
