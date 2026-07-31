#include <ntddk.h>
#include <fwpsk.h>
#include <fwpmk.h>
#include <initguid.h>

#define POSTGRES_PORT 5432
#define MONGO_PORT    27017
#define MYSQL_PORT    3306
#define REDIS_PORT    6379

// {4A7E2B10-3D9C-4E8A-9B1F-2C3D4E5F6A7B}
DEFINE_GUID(
    STUNTDOUBLE_CALLOUT_GUID,
    0x4a7e2b10, 0x3d9c, 0x4e8a, 0x9b, 0x1f, 0x2c, 0x3d, 0x4e, 0x5f, 0x6a, 0x7b
);

// Global Engine State
PDEVICE_OBJECT g_DeviceObject = NULL;
UINT32 g_CalloutId = 0;

extern "C" VOID DriverUnload(PDRIVER_OBJECT DriverObject);

// The core filtering function triggered by the Windows NT Kernel for every network packet
VOID StuntDoubleCalloutClassify(
    const FWPS_INCOMING_VALUES0* inFixedValues,
    const FWPS_INCOMING_METADATA_VALUES0* inMetaValues,
    VOID* layerData,
    const VOID* classifyContext,
    const FWPS_FILTER1* filter,
    UINT64 flowContext,
    FWPS_CLASSIFY_OUT0* classifyOut
) {
    // We only care about outbound IPv4 TCP connections (Layer: FWPS_LAYER_ALE_AUTH_CONNECT_V4)
    if (inFixedValues->layerId != FWPS_LAYER_ALE_AUTH_CONNECT_V4) {
        classifyOut->actionType = FWP_ACTION_PERMIT;
        return;
    }

    // Extract the remote destination port the AI agent is trying to hit
    UINT16 remotePort = inFixedValues->incomingValue[FWPS_FIELD_ALE_AUTH_CONNECT_V4_IP_REMOTE_PORT].value.uint16;

    // Check against StuntDouble's blacklist of database ports
    if (remotePort == POSTGRES_PORT || remotePort == MONGO_PORT || 
        remotePort == MYSQL_PORT || remotePort == REDIS_PORT) {
        
        // Block the connection silently natively in the NT Kernel
        classifyOut->actionType = FWP_ACTION_BLOCK;
        classifyOut->rights &= ~FWPS_RIGHT_ACTION_VALID;

        KdPrint(("[StuntDouble WFP] Blocked rogue AI agent outbound connection to database port %d\n", remotePort));
        return;
    }

    // Permit benign traffic
    classifyOut->actionType = FWP_ACTION_PERMIT;
}

NTSTATUS StuntDoubleCalloutNotify(
    FWPS_CALLOUT_NOTIFY_TYPE notifyType,
    const GUID* filterKey,
    FWPS_FILTER1* filter
) {
    UNREFERENCED_PARAMETER(notifyType);
    UNREFERENCED_PARAMETER(filterKey);
    UNREFERENCED_PARAMETER(filter);
    return STATUS_SUCCESS;
}

extern "C" NTSTATUS DriverEntry(PDRIVER_OBJECT DriverObject, PUNICODE_STRING RegistryPath) {
    UNREFERENCED_PARAMETER(RegistryPath);
    KdPrint(("[StuntDouble WFP] Injecting StuntDouble Zero-Trust Kernel Driver...\n"));
    
    DriverObject->DriverUnload = DriverUnload;

    // Register the WFP Callout
    FWPS_CALLOUT0 callout = { 0 };
    callout.calloutKey = STUNTDOUBLE_CALLOUT_GUID;
    callout.classifyFn = StuntDoubleCalloutClassify;
    callout.notifyFn = StuntDoubleCalloutNotify;

    NTSTATUS status = FwpsCalloutRegister0(DriverObject, &callout, &g_CalloutId);
    if (!NT_SUCCESS(status)) {
        KdPrint(("[StuntDouble WFP] Failed to register WFP Callout. Agent traffic is NOT protected.\n"));
        return status;
    }

    KdPrint(("[StuntDouble WFP] Active! Database ports are now natively blackholed.\n"));
    return STATUS_SUCCESS;
}

extern "C" VOID DriverUnload(PDRIVER_OBJECT DriverObject) {
    UNREFERENCED_PARAMETER(DriverObject);
    KdPrint(("[StuntDouble WFP] Unloading StuntDouble Kernel Driver...\n"));
    
    if (g_CalloutId != 0) {
        FwpsCalloutUnregisterById0(g_CalloutId);
    }
}
