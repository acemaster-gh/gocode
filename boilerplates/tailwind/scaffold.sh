#!/usr/bin/env bash
# =============================================================================
# gocode v1.0.3 — Tailwind CSS CLI scaffold
# Creates proper src/dist structure with Tailwind installed locally
# =============================================================================

DEST="$1"
NAME="$(basename "$DEST")"
SRC="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Copy boilerplate structure
cp -r "$SRC/." "$DEST/"

# Replace PROJECT_NAME placeholder everywhere
find "$DEST" -type f \( -name "*.html" -o -name "*.js" -o -name "*.json" -o -name "*.css" \) | while read -r file; do
    sed -i "s/PROJECT_NAME/${NAME}/g" "$file"
done

cd "$DEST"

# Install Tailwind CSS locally
npm install --save-dev tailwindcss@latest --silent

# Run initial build so dist/css/style.css exists from the start
npx tailwindcss -i ./src/css/input.css -o ./dist/css/style.css --quiet 2>/dev/null || true

echo "Tailwind CSS CLI scaffolded. Run: npm run dev"
