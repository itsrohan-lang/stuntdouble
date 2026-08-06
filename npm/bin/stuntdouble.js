#!/usr/bin/env node

const fs = require('fs');
const { spawn } = require('child_process');
const os = require('os');
const { fetchBinary, binaryPath, DownloadError } = require('../scripts/download');

function run() {
  const child = spawn(binaryPath, process.argv.slice(2), { stdio: 'inherit' });

  child.on('exit', (code, signal) => {
    // Preserve signal-death as a shell-conventional 128+n exit rather than
    // reporting success when the core was killed.
    if (signal) process.exit(128 + (os.constants.signals[signal] || 0));
    process.exit(code ?? 1);
  });
  child.on('error', (err) => {
    console.error(`Error executing StuntDouble core: ${err.message}`);
    process.exit(1);
  });
}

// The core binary is fetched per-platform rather than bundled, so it is absent
// until the postinstall runs. npm 11.19+ gates install scripts behind explicit
// approval, so "absent" is the normal first-run state, not an error — fetch it
// here instead of failing.
if (fs.existsSync(binaryPath)) {
  run();
} else {
  console.error('stuntdouble: core binary not present, fetching it now (one time)...');
  fetchBinary({ log: (m) => console.error(m) })
    .then(run)
    .catch((err) => {
      const detail = err instanceof DownloadError ? err.message : err.stack || err.message;
      console.error(
        `\nstuntdouble: could not obtain the core binary.\n${detail}\n\n` +
          `Alternatively install the CLI directly:\n` +
          `  curl -fsSL https://raw.githubusercontent.com/itsrohan-lang/stuntdouble/main/install.sh | sh\n`
      );
      process.exit(1);
    });
}
