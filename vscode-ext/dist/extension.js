"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = __importStar(require("vscode"));
const child_process_1 = require("child_process");
let statusBarItem;
let isSandboxActive = false;
function activate(context) {
    console.log('StuntDouble extension is now active!');
    const toggleCommandId = 'stuntdouble.toggleSandbox';
    let disposable = vscode.commands.registerCommand(toggleCommandId, () => {
        isSandboxActive = !isSandboxActive;
        updateStatusBarItem();
        if (isSandboxActive) {
            vscode.window.showInformationMessage('StuntDouble: Hardware Sandbox ACTIVATED 🛡️');
            // Execute the CLI to start the warden in the background
            (0, child_process_1.exec)('stuntdouble warden', (error) => {
                if (error)
                    console.error(`Warden execution error: ${error.message}`);
            });
        }
        else {
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
    }
    else {
        statusBarItem.text = `$(stop) StuntDouble: OFF`;
        statusBarItem.backgroundColor = undefined;
        statusBarItem.tooltip = "Sandbox is disabled. Agents have full host access.";
    }
}
function deactivate() { }
//# sourceMappingURL=extension.js.map