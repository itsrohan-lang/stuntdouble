import time
import requests
import json

class Warden:
    """
    The autonomous AI Warden SDK.
    Connects to the StuntDouble Control Plane to monitor agent actions and dynamically block threats.
    """
    def __init__(self, control_plane_url="http://localhost:4439"):
        self.url = control_plane_url
        self.last_seen_id = 0
        print(f"🛡️ Warden Agent initialized. Monitoring Control Plane at {self.url}")

    def fetch_recent_logs(self):
        try:
            res = requests.get(f"{self.url}/api/audit")
            if res.status_code == 200:
                return res.json()
        except requests.exceptions.ConnectionError:
            print("⚠️ Control Plane is unreachable.")
        return []

    def update_policy(self, org_id, blocked_ports, allowed_agents):
        """Deploy a new zero-trust policy dynamically"""
        payload = {
            "org_id": org_id,
            "blocked_ports": blocked_ports,
            "allowed_agents": allowed_agents,
            "strict_egress": True
        }
        try:
            res = requests.post(f"{self.url}/policy", json=payload)
            if res.status_code == 200:
                print(f"🔒 New Policy Deployed! Blocked Ports: {blocked_ports}")
                return True
        except requests.exceptions.ConnectionError:
            pass
        return False

    def watch(self, poll_interval=2.0):
        """Continuously watch for malicious actions and take automated action."""
        print("👀 Warden is now watching the audit stream...")
        while True:
            logs = self.fetch_recent_logs()
            if logs:
                latest = logs[0]
                if latest.get("id", 0) > self.last_seen_id:
                    self.last_seen_id = latest["id"]
                    
                    target = latest.get("target", "")
                    action = latest.get("action", "")
                    
                    if "EXFILTRATION" in action or "api.stripe.com" in target:
                        print(f"🚨 WARDEN DETECTED THREAT: Agent {latest.get('agent_id')} attempted {action} to {target}!")
                        print("🤖 Warden is autonomously locking down the network...")
                        # Automatically block port 443 (HTTPS) dynamically
                        self.update_policy("ent_global_warden_lockdown", [5432, 27017, 3306, 6379, 8080, 443], ["none"])
            
            time.sleep(poll_interval)
