#pragma once

#include <ntddk.h>
#include <fwpsk.h>
#include <fwpmk.h>

#define STUNTDOUBLE_TAG 'tnuS'

// Port Definitions for Destructive AI Targeting
#define POSTGRES_PORT 5432
#define MONGO_PORT 27017
#define MYSQL_PORT 3306
#define REDIS_PORT 6379

// Driver Entry & Unload Methods
extern "C" NTSTATUS DriverEntry(PDRIVER_OBJECT DriverObject, PUNICODE_STRING RegistryPath);
extern "C" VOID DriverUnload(PDRIVER_OBJECT DriverObject);

// WFP Callout Functions
VOID StuntDoubleCalloutClassify(
    const FWPS_INCOMING_VALUES0* inFixedValues,
    const FWPS_INCOMING_METADATA_VALUES0* inMetaValues,
    VOID* layerData,
    const VOID* classifyContext,
    const FWPS_FILTER1* filter,
    UINT64 flowContext,
    FWPS_CLASSIFY_OUT0* classifyOut
);

NTSTATUS StuntDoubleCalloutNotify(
    FWPS_CALLOUT_NOTIFY_TYPE notifyType,
    const GUID* filterKey,
    FWPS_FILTER1* filter
);
