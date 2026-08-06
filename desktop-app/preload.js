// Preload: the only bridge between the renderer and Node.
//
// Runs with Node available but in an isolated context, so the renderer sees
// exactly the four functions below and nothing else — no ipcRenderer, no
// require, no process. Keep this surface as small as it is.

const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('stuntdouble', {
  // Ask the main process to launch a sandbox. The name is re-validated in the
  // main process; the check here is only to fail fast with a clear message.
  startSandbox: (agent) => ipcRenderer.send('start-sandbox', agent),

  // Subscribe to streamed output. Handlers receive only the payload — the
  // IpcRendererEvent is deliberately not forwarded, since it carries `sender`
  // and would hand the renderer a route back to the ipc surface.
  onOutput: (handler) => {
    const wrapped = (_event, chunk) => handler(chunk);
    ipcRenderer.on('sandbox-output', wrapped);
    return () => ipcRenderer.removeListener('sandbox-output', wrapped);
  },

  onStatus: (handler) => {
    const wrapped = (_event, status) => handler(status);
    ipcRenderer.on('sandbox-status', wrapped);
    return () => ipcRenderer.removeListener('sandbox-status', wrapped);
  },
});
