import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
    console.log('StuntDouble VS Code extension is now active!');

    // Command 1: Run an Agent
    let disposableRun = vscode.commands.registerCommand('stuntdouble.runAgent', async () => {
        const agentName = await vscode.window.showInputBox({
            prompt: 'Enter the AI agent to sandbox (e.g., claude, aider, sh)',
            placeHolder: 'claude'
        });

        if (agentName) {
            const terminal = vscode.window.createTerminal('StuntDouble Sandbox');
            terminal.show();
            terminal.sendText(`stuntdouble run ${agentName}`);
            vscode.window.showInformationMessage(`Locking down kernel and starting ${agentName}...`);
        }
    });

    // Command 2: Time-Travel Rewind
    let disposableRewind = vscode.commands.registerCommand('stuntdouble.rewind', () => {
        const terminal = vscode.window.createTerminal('StuntDouble Rewind');
        terminal.show();
        terminal.sendText(`stuntdouble rewind`);
        vscode.window.showInformationMessage(`Initiating zero-copy time travel...`);
    });

    // Command 3: Open Dashboard WebView
    let disposableDashboard = vscode.commands.registerCommand('stuntdouble.dashboard', () => {
        const panel = vscode.window.createWebviewPanel(
            'stuntDoubleDashboard',
            'Mission Control',
            vscode.ViewColumn.Two,
            { enableScripts: true }
        );

        // Load the local Next.js dashboard inside VS Code
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
                <iframe src="http://localhost:3000/dashboard"></iframe>
            </body>
            </html>
        `;
    });

    context.subscriptions.push(disposableRun, disposableRewind, disposableDashboard);
}

export function deactivate() {}
