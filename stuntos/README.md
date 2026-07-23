# 💿 StuntOS (Bare-Metal Isolation)

StuntOS is a custom, 20MB bootable Linux Operating System explicitly designed to run untrusted AI workloads natively on bare metal hardware (Stunt Boxes).

## Build Instructions

To compile the `StuntOS.iso`, you will need a Linux machine with `buildroot` installed.

1. Download Buildroot:
```bash
wget https://buildroot.org/downloads/buildroot-2024.02.tar.gz
tar -xvf buildroot-2024.02.tar.gz
cd buildroot-2024.02
```

2. Apply the StuntOS Configuration:
```bash
cp ../stuntos/buildroot_config .config
make olddefconfig
```

3. Compile the Kernel and ISO:
```bash
make
```

The output will be an ultra-lean `.iso` file located in `output/images/rootfs.iso9660`. You can flash this to a USB drive or boot it directly in Proxmox/VMware.
