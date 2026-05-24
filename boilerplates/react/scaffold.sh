#!/usr/bin/env bash
DEST="$1"
NAME="$(basename "$DEST")"
TMP=$(mktemp -d)
cd "$TMP"
npm create vite@latest "$NAME" -- --template react --yes 2>/dev/null
cp -r "$NAME/." "$DEST/"
rm -rf "$TMP"
cd "$DEST"
npm install --silent
