import * as vscode from 'vscode';
import * as http from 'http';

let isStealthModeActive = false;

export function activate(context: vscode.ExtensionContext) {
    console.log('StuntDouble VS Code extension is now active!');

    // Command: Run an Agent
    let disposableRun = vscode.commands.registerCommand('stuntdouble.runAgent', async () => {
        const agentName = await vscode.window.showInputBox({
            prompt: 'Enter the AI agent to sandbox (e.g., claude, aider, sh)',
            placeHolder: 'claude'
        });

        if (agentName) {
            const terminal = vscode.window.createTerminal('StuntDouble Sandbox');
            terminal.show();
            // Changed from stuntdouble to sd
            terminal.sendText(`sd run ${agentName}`);
            vscode.window.showInformationMessage(`Locking down kernel and starting ${agentName}...`);
        }
    });

    // Command: Toggle Stealth Terminal Interceptor
    let disposableStealth = vscode.commands.registerCommand('stuntdouble.toggleStealth', () => {
        isStealthModeActive = !isStealthModeActive;
        if (isStealthModeActive) {
            vscode.window.showInformationMessage('🛡️ StuntDouble Stealth Mode ACTIVATED. All new terminals will be intercepted and sandboxed.');
        } else {
            vscode.window.showInformationMessage('⚠️ StuntDouble Stealth Mode DEACTIVATED.');
        }
    });

    // Event: Intercept Terminal Creation
    vscode.window.onDidOpenTerminal((terminal: vscode.Terminal) => {
        if (isStealthModeActive && terminal.name !== 'StuntDouble Sandbox') {
            // Wait a split second for the shell to initialize, then inject the sandbox wrapper
            setTimeout(() => {
                terminal.sendText(`sd run bash`);
                vscode.window.showWarningMessage(`StuntDouble intercepted terminal: ${terminal.name}`);
            }, 500);
        }
    });

    // Command: Open Dashboard WebView
    let disposableDashboard = vscode.commands.registerCommand('stuntdouble.dashboard', () => {
        const panel = vscode.window.createWebviewPanel(
            'stuntDoubleDashboard',
            'Mission Control',
            vscode.ViewColumn.Two,
            { enableScripts: true }
        );

        panel.webview.html = `
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <style>
                    body, html { margin: 0; padding: 0; height: 100%; border: none; }
                    iframe { width: 100%; height: 100%; border: none; }
                </style>
            </head>
            <body>
                <iframe src="http://localhost:3000/"></iframe>
            </body>
            </html>
        `;
    });

    let lastSeenLogId = 0;

    // Background polling for new blocks from Control Plane
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
                                vscode.window.showWarningMessage(`🚨 StuntDouble Blocked Agent! The agent '${latestLog.agent_id}' attempted to access '${latestLog.target}' but was intercepted by enterprise policy.`);
                            }
                            lastSeenLogId = latestLog.id;
                        }
                    }
                } catch (e) {
                    // Control Plane might not be running yet
                }
            });
        }).on('error', () => {
            // Silently ignore if Control Plane is not running
        });
    }, 3000);

    context.subscriptions.push(disposableRun, disposableStealth, disposableDashboard);
}

export function deactivate() {}
