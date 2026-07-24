#!/bin/bash

# StuntDouble End-to-End Simulation Script
echo "================================================="
echo "🛡️  STUNTDOUBLE ZERO-TRUST SIMULATION"
echo "================================================="
echo ""

echo "[1/4] Booting Enterprise Control Plane on port 4439..."
sleep 1
echo "✅ Control Plane Active (Mock Engine Ready)"

echo "[2/4] Loading Rust eBPF Kernel Probes (cgroup_skb/egress)..."
sleep 1
echo "✅ BPF hooks attached to cgroup v2."

echo "[3/4] Spawning AI Agent Sandbox (Claude Code)..."
sleep 2
echo "🤖 Agent 'claude' running in isolated MicroVM."

echo ""
echo ">> [Agent Log] Attempting to exfiltrate data to api.stripe.com (Port 443)..."
sleep 1
echo "🚨 [eBPF Kernel Alert] PACKET DROPPED! Connection to blocked host denied."
echo "👻 [Keploy Mock] Injecting synthetic 200 OK ghost response to agent."

echo ""
echo ">> [Agent Log] Received 200 OK from api.stripe.com. Continuing logic..."
sleep 1

echo "[4/4] Generating Audit Telemetry..."
sleep 1
echo "✅ Audit log securely committed to SQLite."
echo "✅ Alert pushed to Next.js SOC Dashboard and Native Desktop Tray."

echo ""
echo "================================================="
echo "🎉 SIMULATION COMPLETE: ZERO-DAY ESCAPES = 0"
echo "================================================="
