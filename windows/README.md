# StuntDouble WFP Kernel Driver (Windows)

Because eBPF is native to Linux, StuntDouble uses the **Windows Filtering Platform (WFP)** to achieve the exact same kernel-level network interception on Windows host machines.

## Architecture
The `wfp_interceptor.cpp` is a native Windows NT Kernel mode driver. It hooks into the `FWPS_LAYER_ALE_AUTH_CONNECT_V4` layer. Every single time an AI agent (running natively or inside Docker) tries to open an outbound TCP connection, the Windows Kernel pauses the packet and asks this driver for permission. 

If the agent tries to hit port `5432` (Postgres) or `27017` (Mongo), this driver responds with `FWP_ACTION_BLOCK`, silently dropping the packet at the lowest possible level of the OS before it ever reaches your actual database.

## Compilation Instructions
*Note: This driver must be compiled on a Windows machine using the Microsoft MSVC compiler.*

1. Install **Visual Studio 2022** with the "Desktop development with C++" workload.
2. Install the **Windows Driver Kit (WDK)**.
3. Open the `Visual Studio Developer Command Prompt`.
4. Compile the sys file:
   ```cmd
   cl.exe /kernel /I"include" /c src\wfp_interceptor.cpp
   link.exe /subsystem:native /driver /entry:DriverEntry wfp_interceptor.obj fwpuclnt.lib fwpkclnt.lib
   ```
5. To install the compiled `.sys` driver on a test machine, you must enable Test Signing mode (`bcdedit /set testsigning on`) because this is a kernel-mode driver!
