const { app, BrowserWindow, ipcMain } = require('electron');
const path = require('path');
const { exec } = require('child_process');

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

app.whenReady().then(createWindow);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

const { spawn } = require('child_process');

// IPC handler to start sandbox
ipcMain.on('start-sandbox', (event, agent) => {
  console.log("Starting sandbox for", agent);
  
  // Try using the local compiled binary if we are running in the repo
  const localSdPath = path.join(__dirname, '..', 'cli', 'sd');
  
  // Run the command using spawn to stream output in real-time
  const child = spawn(localSdPath, ['run', agent]);

  // Fallback if localSdPath fails to spawn
  child.on('error', (err) => {
    const fallback = spawn('sd', ['run', agent]);
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
