#!/bin/bash
set -e
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
cd "$SCRIPT_DIR"
if command -v pkill > /dev/null 2>&1; then
    pkill -f electron 2>/dev/null || true
elif command -v killall > /dev/null 2>&1; then
    killall electron 2>/dev/null || true
fi
npm config set registry https://registry.npmmirror.com/
npm install
npm run build
