const { app, BrowserWindow, ipcMain, Tray, Menu, Notification } = require('electron');
const path = require('path');
const { exec, spawn } = require('child_process');
const http = require('http');

let tray = null;
let lastSeenLogId = 0;

function createWindow() {
  const win = new BrowserWindow({
    width: 900,
    height: 700,
    title: "StuntDouble Desktop",
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false
    }
  });

  win.loadFile('index.html');
}

function createTray() {
  // Use a generic icon or the official logo if available
  tray = new Tray(path.join(__dirname, '..', 'docs', 'assets', 'logo.png'));
  const contextMenu = Menu.buildFromTemplate([
    { label: 'StuntDouble Engine: Active', type: 'normal', enabled: false },
    { type: 'separator' },
    { label: 'Open Dashboard', click: () => createWindow() },
    { label: 'Quit', click: () => app.quit() }
  ]);
  tray.setToolTip('StuntDouble Zero-Trust Sandbox');
  tray.setContextMenu(contextMenu);
}

function startAuditPolling() {
  setInterval(() => {
    http.get('http://localhost:4439/api/audit', (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          const logs = JSON.parse(data);
          if (logs && logs.length > 0) {
            const latestLog = logs[0];
            if (latestLog.id > lastSeenLogId) {
              if (lastSeenLogId !== 0 && latestLog.status.includes('Blocked')) {
                new Notification({
                  title: '🚨 StuntDouble Security Alert',
                  body: `Agent '${latestLog.agent_id}' was blocked from accessing '${latestLog.target}'!`,
                  icon: path.join(__dirname, '..', 'docs', 'assets', 'logo.png')
                }).show();
              }
              lastSeenLogId = latestLog.id;
            }
          }
        } catch (e) {
          // ignore
        }
      });
    }).on('error', () => {
      // ignore
    });
  }, 3000);
}

app.whenReady().then(() => {
  createWindow();
  createTray();
  startAuditPolling();
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

// IPC handler to start sandbox
ipcMain.on('start-sandbox', (event, agent) => {
  console.log("Starting sandbox for", agent);
  
  // Try using the local compiled binary if we are running in the repo
  const localSdPath = path.join(__dirname, '..', 'cli', 'sd');
  
  // Run the command using spawn to stream output in real-time
  // Use --remote to offload to StuntDouble Cloud, preventing slow local Docker pulls
  const child = spawn(localSdPath, ['run', agent, '--remote']);

  // Fallback if localSdPath fails to spawn
  child.on('error', (err) => {
    const fallback = spawn('sd', ['run', agent, '--remote']);
    streamChild(fallback, event);
  });

  streamChild(child, event);
});

function streamChild(child, event) {
  child.stdout.on('data', (data) => {
    event.reply('sandbox-output', data.toString());
  });

  child.stderr.on('data', (data) => {
    event.reply('sandbox-output', data.toString());
  });

  child.on('close', (code) => {
    event.reply('sandbox-status', {
      success: code === 0,
      output: `\n[Process exited with code ${code}]\n`
    });
  });
}
