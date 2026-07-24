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

// IPC handler to start sandbox
ipcMain.on('start-sandbox', (event, agent) => {
  console.log("Starting sandbox for", agent);
  exec(`sd run ${agent}`, (error, stdout, stderr) => {
    event.reply('sandbox-status', {
      success: !error,
      output: stdout || stderr
    });
  });
});
