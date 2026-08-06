const { app, BrowserWindow, ipcMain, Tray, Menu, Notification, shell } = require('electron');
const path = require('path');
const { exec, spawn } = require('child_process');
const http = require('http');

let tray = null;
let win = null;
let lastSeenLogId = 0;

function createWindow() {
  // Reuse the existing window — the tray menu can fire this repeatedly.
  if (win && !win.isDestroyed()) {
    win.show();
    win.focus();
    return win;
  }

  win = new BrowserWindow({
    width: 900,
    height: 700,
    title: "StuntDouble Desktop",
    webPreferences: {
      // The renderer displays audit-log text and agent stdout, both of which
      // are ultimately agent-controlled. Node must not be reachable from it:
      // with nodeIntegration on, a single innerHTML sink anywhere in this UI
      // would be remote code execution on the host.
      nodeIntegration: false,
      contextIsolation: true,
      sandbox: true,
      preload: path.join(__dirname, 'preload.js'),
    }
  });

  // This window only ever shows the bundled local UI. Anything that tries to
  // navigate it elsewhere, or open a child window, is a bug or an attack.
  win.webContents.on('will-navigate', (event) => event.preventDefault());
  win.webContents.setWindowOpenHandler(({ url }) => {
    if (url.startsWith('https://')) shell.openExternal(url);
    return { action: 'deny' };
  });

  win.on('closed', () => { win = null; });

  win.loadFile('index.html');
  return win;
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

// Agent names are passed to `sd run <agent>`. spawn() is used without a shell,
// so this is not a shell-injection surface — but an unchecked value starting
// with "-" would be parsed by sd as a flag rather than an agent name, and one
// containing a path separator could point at another executable. Require a
// plain leading-alphanumeric identifier.
const AGENT_PATTERN = /^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$/;

// IPC handler to start sandbox
ipcMain.on('start-sandbox', (event, agent) => {
  if (typeof agent !== 'string' || !AGENT_PATTERN.test(agent)) {
    event.reply('sandbox-status', {
      success: false,
      output: `\n[Rejected: "${String(agent).slice(0, 64)}" is not a valid agent name]\n`,
    });
    return;
  }

  console.log("Starting sandbox for", agent);

  // Try using the local compiled binary if we are running in the repo
  const localSdPath = path.join(__dirname, '..', 'cli', 'sd');

  // Run the command using spawn to stream output in real-time
  // Use --remote to offload to StuntDouble Cloud, preventing slow local Docker pulls
  const child = spawn(localSdPath, ['run', agent, '--remote']);

  // Fall back to sd on PATH when the repo-local binary is not present. The
  // fallback replaces the failed child rather than running alongside it, so
  // output is never streamed from both.
  let usedFallback = false;
  child.on('error', () => {
    usedFallback = true;
    const fallback = spawn('sd', ['run', agent, '--remote']);
    fallback.on('error', (err) => {
      event.reply('sandbox-status', {
        success: false,
        output: `\n[Could not launch sd: ${err.message}]\n`,
      });
    });
    streamChild(fallback, event);
  });

  streamChild(child, event, () => usedFallback);
});

function streamChild(child, event, suppressed = () => false) {
  child.stdout.on('data', (data) => {
    event.reply('sandbox-output', data.toString());
  });

  child.stderr.on('data', (data) => {
    event.reply('sandbox-output', data.toString());
  });

  child.on('close', (code) => {
    // A child that failed to spawn also emits 'close'; the fallback reports
    // for it.
    if (suppressed()) return;
    event.reply('sandbox-status', {
      success: code === 0,
      output: `\n[Process exited with code ${code}]\n`
    });
  });
}
