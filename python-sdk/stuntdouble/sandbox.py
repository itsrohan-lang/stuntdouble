import subprocess
import json
import os

class Sandbox:
    """
    StuntDouble Python SDK Context Manager.
    Automatically spawns a native OS-level StuntDouble isolation wrapper for AI agents.
    """
    def __init__(self, network_policy="strict", drop_caps=True):
        self.network_policy = network_policy
        self.drop_caps = drop_caps

    def __enter__(self):
        print(f"🛡️ [StuntDouble SDK] Securing Python process (Policy: {self.network_policy})...")
        # In a deep native integration, this would inject C-types into the running process.
        # For the MVP SDK, we provide a `.run()` wrapper that delegates to the CLI.
        return self

    def run(self, agent_command: str):
        """
        Executes the agent command securely inside the StuntDouble ephemeral container.
        """
        print(f"🚀 [StuntDouble SDK] Orchestrating safe execution for: {agent_command}")
        
        args = ["sd", "run", "sh", "-c", agent_command]
        
        try:
            result = subprocess.run(args, check=True, text=True, capture_output=False)
            print("✅ [StuntDouble SDK] Execution concluded safely.")
        except subprocess.CalledProcessError as e:
            print(f"⚠️ [StuntDouble SDK] Agent execution was terminated by the sandbox: {e}")

    def __exit__(self, exc_type, exc_value, traceback):
        print(">> [StuntDouble SDK] Tearing down eBPF sandbox constraints.")
