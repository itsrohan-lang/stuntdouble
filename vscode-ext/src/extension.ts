import * as vscode from 'vscode';
import { exec } from 'child_process';

let statusBarItem: vscode.StatusBarItem;
let isSandboxActive = false;

export function activate(context: vscode.ExtensionContext) {
    console.log('StuntDouble extension is now active!');

    const toggleCommandId = 'stuntdouble.toggleSandbox';
    let disposable = vscode.commands.registerCommand(toggleCommandId, () => {
        isSandboxActive = !isSandboxActive;
        updateStatusBarItem();
        
        if (isSandboxActive) {
            vscode.window.showInformationMessage('StuntDouble: Hardware Sandbox ACTIVATED 🛡️');
            
            // Execute the CLI to start the warden in the background
            exec('stuntdouble warden', (error) => {
                if (error) console.error(`Warden execution error: ${error.message}`);
            });
        } else {
            vscode.window.showWarningMessage('StuntDouble: Sandbox DEACTIVATED. Run agents at your own risk.');
        }
    });

    context.subscriptions.push(disposable);

    // Create a status bar item at the bottom right
    statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
    statusBarItem.command = toggleCommandId;
    context.subscriptions.push(statusBarItem);
    
    updateStatusBarItem();
    statusBarItem.show();
}

function updateStatusBarItem() {
    if (isSandboxActive) {
        statusBarItem.text = `$(shield) StuntDouble: ON`;
        statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
        statusBarItem.tooltip = "Sandbox is intercepting AI network calls via eBPF";
    } else {
        statusBarItem.text = `$(stop) StuntDouble: OFF`;
        statusBarItem.backgroundColor = undefined;
        statusBarItem.tooltip = "Sandbox is disabled. Agents have full host access.";
    }
}

export function deactivate() {}
