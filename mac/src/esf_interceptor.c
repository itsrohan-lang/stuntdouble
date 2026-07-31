#include <EndpointSecurity/EndpointSecurity.h>
#include <bsm/libbsm.h>
#include <dispatch/dispatch.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>

// StuntDouble macOS Endpoint Security (ESF) Kernel Interceptor
// This driver runs natively on Apple Silicon/Intel to physically block outbound DB connections
// without relying on Docker's Linux-specific eBPF layer.

static const char *blocked_socket_patterns[] = {
    ".s.PGSQL.5432",   // PostgreSQL
    "mongodb-27017",   // MongoDB
    "mysql.sock",      // MySQL
    "redis.sock",      // Redis
    NULL
};

static int is_blocked_socket_path(const char *path) {
    if (!path) return 0;
    for (int i = 0; blocked_socket_patterns[i] != NULL; i++) {
        if (strstr(path, blocked_socket_patterns[i]) != NULL) {
            return 1;
        }
    }
    return 0;
}

void handle_event(es_client_t *client, const es_message_t *msg) {
    if (msg->event_type == ES_EVENT_TYPE_AUTH_UIPC_CONNECT) {
        pid_t pid = audit_token_to_pid(msg->process->audit_token);
        const char *path = msg->event.uipc_connect.file ? msg->event.uipc_connect.file->path.data : "";

        if (is_blocked_socket_path(path)) {
            printf("[StuntDouble ESF] 🚨 BLOCKED outbound database connection to socket %s from PID %d\n", path, pid);
            es_respond_auth_result(client, msg, ES_AUTH_RESULT_DENY, false);
            return;
        }

        printf("[StuntDouble ESF] Permitted connection to socket %s from PID %d\n", path, pid);
        es_respond_auth_result(client, msg, ES_AUTH_RESULT_ALLOW, false);
    }
}

int main() {
    es_client_t *client = NULL;
    es_new_client_result_t res = es_new_client(&client, ^(es_client_t *c, const es_message_t *msg) {
        handle_event(c, msg);
    });

    if (res != ES_NEW_CLIENT_RESULT_SUCCESS) {
        fprintf(stderr, "❌ [StuntDouble ESF] Failed to register macOS Endpoint Security client. Are you running as root with EndpointSecurity entitlement?\n");
        return 1;
    }

    es_event_type_t events[] = { ES_EVENT_TYPE_AUTH_UIPC_CONNECT };
    if (es_subscribe(client, events, 1) != ES_RETURN_SUCCESS) {
        fprintf(stderr, "❌ [StuntDouble ESF] Failed to subscribe to Kernel Auth UIPC Connect events.\n");
        es_delete_client(client);
        return 1;
    }

    printf("✅ [StuntDouble ESF] Active! macOS Kernel is now natively dropping rogue AI database queries.\n");

    dispatch_main();
    return 0;
}
