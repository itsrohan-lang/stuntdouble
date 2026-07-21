const vscode = require('vscode');
const { exec } = require('child_process');

function activate(context) {
    console.log('StuntDouble Sandbox is now active in the background.');
    
    // Automatically inject StuntDouble context when the extension loads
    exec('stuntdouble init', (err, stdout) => {
        if (!err) {
            vscode.window.showInformationMessage('🔒 StuntDouble Sandbox Enabled for all workspace agents.');
        }
    });
}

function deactivate() {}

module.exports = { activate, deactivate };
