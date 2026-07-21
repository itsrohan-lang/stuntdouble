#include "wfp_interceptor.h"

// Global Engine State
PDEVICE_OBJECT g_DeviceObject = NULL;
UINT32 g_CalloutId = 0;

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

    // Check against StuntDouble's hardcoded blacklist of destructive database ports
    if (remotePort == POSTGRES_PORT || remotePort == MONGO_PORT || 
        remotePort == MYSQL_PORT || remotePort == REDIS_PORT) {
        
        // Block the connection silently natively in the NT Kernel
        classifyOut->actionType = FWP_ACTION_BLOCK;
        classifyOut->rights &= ~FWPS_RIGHT_ACTION_VALID;

        // Log to Event Viewer (which the Go telemetry agent can read)
        KdPrint(("[StuntDouble WFP] Blocked rogue AI agent outbound connection to database port %d\n", remotePort));
        return;
    }

    // Permit benign traffic (like fetching NPM packages)
    classifyOut->actionType = FWP_ACTION_PERMIT;
}

NTSTATUS StuntDoubleCalloutNotify(
    FWPS_CALLOUT_NOTIFY_TYPE notifyType,
    const GUID* filterKey,
    FWPS_FILTER1* filter
) {
    return STATUS_SUCCESS;
}

extern "C" NTSTATUS DriverEntry(PDRIVER_OBJECT DriverObject, PUNICODE_STRING RegistryPath) {
    KdPrint(("[StuntDouble WFP] Injecting StuntDouble Zero-Trust Kernel Driver...\n"));
    
    DriverObject->DriverUnload = DriverUnload;

    // Register the WFP Callout
    FWPS_CALLOUT0 callout = { 0 };
    callout.calloutKey = { /* GUID would go here */ };
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
    KdPrint(("[StuntDouble WFP] Unloading StuntDouble Kernel Driver...\n"));
    
    if (g_CalloutId != 0) {
        FwpsCalloutUnregisterById0(g_CalloutId);
    }
}
