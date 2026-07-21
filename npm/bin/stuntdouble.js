#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');

// Execute the bundled Go binary in the same bin folder
const binaryPath = path.resolve(__dirname, 'sd');
const args = process.argv.slice(2);

const child = spawn(binaryPath, args, { stdio: 'inherit' });

child.on('exit', (code) => process.exit(code));
child.on('error', (err) => {
  console.error(`Error executing StuntDouble core: ${err.message}`);
  process.exit(1);
});
