#!/bin/bash
cd "$(dirname "$0")"
REMOTE=$(git remote | head -1)
if [ -z "$REMOTE" ]; then
    echo "Error: No git remotes found!"
    exit 1
fi
git fetch "$REMOTE"
git checkout master
git reset --hard "$REMOTE/master"
git clean -fd
chmod +x "$0"