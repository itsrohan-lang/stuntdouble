//
// Shared, verified downloader for the StuntDouble core binary.
//
// Used by two callers:
//   - scripts/postinstall.js — eager fetch at install time (the happy path)
//   - bin/stuntdouble.js     — lazy fetch on first run, for when the postinstall
//                              never ran. npm 11.19+ withholds install scripts
//                              pending `npm install-scripts approve`, and pnpm
//                              /yarn have comparable gating, so a package that
//                              relies on postinstall alone is broken by default
//                              for a large share of users.
//
// Fails closed in both paths: a missing checksum or a mismatch aborts rather
// than leaving an unverified binary on disk. This mirrors install.sh.

const fs = require('fs');
const path = require('path');
const https = require('https');
const crypto = require('crypto');

const REPO = 'itsrohan-lang/stuntdouble';
const MAX_REDIRECTS = 5;

// release.yml publishes exactly these asset names; keep the two in sync.
const TARGETS = {
  'linux:x64': 'stuntdouble-linux-amd64',
  'linux:arm64': 'stuntdouble-linux-arm64',
  'darwin:x64': 'stuntdouble-darwin-amd64',
  'darwin:arm64': 'stuntdouble-darwin-arm64',
  'win32:x64': 'stuntdouble-windows-amd64.exe',
};

const binDir = path.resolve(__dirname, '..', 'bin');
const binaryPath = path.join(binDir, process.platform === 'win32' ? 'sd.exe' : 'sd');

class DownloadError extends Error {}

function get(url, version, redirectsLeft = MAX_REDIRECTS) {
  return new Promise((resolve, reject) => {
    https
      .get(url, { headers: { 'User-Agent': `stuntdouble-npm/${version}` } }, (res) => {
        // GitHub redirects release downloads to a CDN host; follow manually.
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          res.resume();
          if (redirectsLeft === 0) return reject(new DownloadError(`too many redirects for ${url}`));
          return resolve(get(new URL(res.headers.location, url).toString(), version, redirectsLeft - 1));
        }
        if (res.statusCode !== 200) {
          res.resume();
          return reject(new DownloadError(`HTTP ${res.statusCode} for ${url}`));
        }
        const chunks = [];
        res.on('data', (c) => chunks.push(c));
        res.on('end', () => resolve(Buffer.concat(chunks)));
        res.on('error', reject);
      })
      .on('error', reject);
  });
}

// Downloads and verifies the core binary, writing it to bin/. Resolves to the
// binary path. Throws DownloadError with an actionable message on any failure.
async function fetchBinary({ log = () => {} } = {}) {
  const pkg = require('../package.json');

  const key = `${process.platform}:${process.arch}`;
  const asset = TARGETS[key];
  if (!asset) {
    throw new DownloadError(
      `unsupported platform ${key}.\n` +
        `Supported: ${Object.keys(TARGETS).join(', ')}.\n` +
        `Build from source instead: https://github.com/${REPO}`
    );
  }

  const tag = process.env.STUNTDOUBLE_VERSION || `v${pkg.version}`;
  const base = `https://github.com/${REPO}/releases/download/${tag}`;

  log(`stuntdouble: downloading ${asset} (${tag})...`);

  let binary;
  try {
    binary = await get(`${base}/${asset}`, pkg.version);
  } catch (err) {
    throw new DownloadError(
      `download failed: ${base}/${asset}\n` +
        `${err.message}\n` +
        `If release ${tag} has no asset for this platform, install from source or set ` +
        `STUNTDOUBLE_SKIP_DOWNLOAD=1 to bypass.`
    );
  }

  let sums;
  try {
    sums = await get(`${base}/SHA256SUMS`, pkg.version);
  } catch (err) {
    throw new DownloadError(
      `could not download SHA256SUMS for ${tag} (${err.message}). ` +
        `Refusing to install an unverified binary.`
    );
  }

  // Lines are "<sha256>  <name>" (or " *<name>" in binary mode).
  const expected = sums
    .toString('utf8')
    .split('\n')
    .map((l) => l.trim().split(/[ \t]+/))
    .filter((p) => p.length >= 2 && p[1].replace(/^\*/, '') === asset)
    .map((p) => p[0])[0];

  if (!expected) {
    throw new DownloadError(`no checksum for ${asset} in SHA256SUMS. Refusing to install.`);
  }

  const actual = crypto.createHash('sha256').update(binary).digest('hex');
  if (actual !== expected) {
    throw new DownloadError(
      `checksum mismatch for ${asset}\n` +
        `  expected: ${expected}\n` +
        `  actual:   ${actual}\n` +
        `The download may be corrupt or tampered with. Nothing was installed.`
    );
  }

  // Write via a unique temp file and rename, so a concurrent `sd` invocation
  // never observes a half-written binary.
  fs.mkdirSync(binDir, { recursive: true });
  const tmp = `${binaryPath}.${process.pid}.tmp`;
  fs.writeFileSync(tmp, binary, { mode: 0o755 });
  fs.renameSync(tmp, binaryPath);

  log(`stuntdouble: checksum verified, core installed (${tag}).`);
  return binaryPath;
}

module.exports = { fetchBinary, binaryPath, DownloadError, TARGETS };
