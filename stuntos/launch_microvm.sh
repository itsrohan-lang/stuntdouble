#!/bin/bash

# StuntOS MicroVM Launcher (Powered by AWS Firecracker)
# Boot an isolated hardware-level Virtual Machine in <150ms for an AI agent.

KERNEL_IMAGE="vmlinux.bin"
ROOTFS_IMAGE="rootfs.ext4"
API_SOCKET="/tmp/firecracker.socket"
AGENT_NAME=$1

if [ -z "$AGENT_NAME" ]; then
    echo "Usage: ./launch_microvm.sh <agent_name>"
    exit 1
fi

echo "🔥 Initializing StuntOS Hardware Isolation for agent: $AGENT_NAME"

# Clean up any existing socket
rm -f $API_SOCKET

# Start Firecracker in the background
firecracker --api-sock $API_SOCKET &
FC_PID=$!
sleep 0.1

echo "⚙️ Configuring MicroVM CPU and Memory..."
curl -X PUT --unix-socket $API_SOCKET \
    http://localhost/machine-config \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{
        "vcpu_count": 2,
        "mem_size_mib": 512,
        "smt": false
    }' > /dev/null 2>&1

echo "🐧 Loading StuntOS Kernel..."
curl -X PUT --unix-socket $API_SOCKET \
    http://localhost/boot-source \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{
        "kernel_image_path": "'$(pwd)'/'$KERNEL_IMAGE'",
        "boot_args": "console=ttyS0 reboot=k panic=1 pci=off"
    }' > /dev/null 2>&1

echo "📦 Mounting Read-Only Root Filesystem..."
curl -X PUT --unix-socket $API_SOCKET \
    http://localhost/drives/rootfs \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{
        "drive_id": "rootfs",
        "path_on_host": "'$(pwd)'/'$ROOTFS_IMAGE'",
        "is_root_device": true,
        "is_read_only": true
    }' > /dev/null 2>&1

echo "🚀 Booting StuntOS..."
curl -X PUT --unix-socket $API_SOCKET \
    http://localhost/actions \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{
        "action_type": "InstanceStart"
    }' > /dev/null 2>&1

echo "✅ Bare-Metal Sandbox Active!"
echo "Type 'exit' to kill the MicroVM."
wait $FC_PID
