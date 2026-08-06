#!/usr/bin/env node
//
// Eagerly fetches the StuntDouble core binary at install time.
//
// This is best-effort, not load-bearing: it exits 0 even on failure, because
// bin/stuntdouble.js performs the same verified download lazily on first run.
// That matters because npm 11.19+ withholds install scripts pending
// `npm install-scripts approve`, so this script may never execute at all —
// aborting the install on a transient network error would break `npm install`
// for no benefit while the run-time path can still recover.
//
// The security property is preserved either way: an unverified binary is never
// written to disk. See scripts/download.js.
//
// Environment:
//   STUNTDOUBLE_VERSION        download a specific tag (default: package version)
//   STUNTDOUBLE_SKIP_DOWNLOAD  skip entirely (offline installs, CI image builds)

const { fetchBinary } = require('./download');

if (process.env.STUNTDOUBLE_SKIP_DOWNLOAD) {
  console.log('stuntdouble: STUNTDOUBLE_SKIP_DOWNLOAD set, skipping core download.');
  process.exit(0);
}

fetchBinary({ log: (m) => console.log(m) }).catch((err) => {
  console.warn(
    `\nstuntdouble: could not pre-fetch the core binary.\n` +
      `${err.message}\n\n` +
      `This is not fatal — it will be fetched on first run of \`sd\`.\n`
  );
  process.exit(0);
});
