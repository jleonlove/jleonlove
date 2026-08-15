#!/usr/bin/env bash
set -euo pipefail
node_expected=22.16.0
npm_expected=10.9.2
[[ "$(node -v)" == "v${node_expected}" ]] || { echo "wrong node"; exit 20; }
[[ "$(npm -v)" == "${npm_expected}" ]] || { echo "wrong npm"; exit 21; }
rm -f package-lock.json
npm install --package-lock-only --ignore-scripts --no-audit --no-fund
npm ci --ignore-scripts --no-audit --no-fund
npm run qualify
