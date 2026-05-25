#!/usr/bin/env bash
# gocode v1.0.3 — Tailwind CSS v4 scaffold

DEST="$1"
NAME="$(basename "$DEST")"
SRC="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Create folder structure
mkdir -p "$DEST/src/css" "$DEST/src/js" "$DEST/src/components" "$DEST/dist/css"

# Copy all boilerplate files
cp "$SRC/src/index.html"            "$DEST/src/index.html"
cp "$SRC/src/css/input.css"         "$DEST/src/css/input.css"
cp "$SRC/src/js/main.js"            "$DEST/src/js/main.js"
cp "$SRC/src/components/header.html" "$DEST/src/components/header.html"
cp "$SRC/src/components/footer.html" "$DEST/src/components/footer.html"
cp "$SRC/package.json"              "$DEST/package.json"
cp "$SRC/.gitignore"                "$DEST/.gitignore"

# Replace PROJECT_NAME placeholder
find "$DEST" -type f \( -name "*.html" -o -name "*.js" -o -name "*.json" -o -name "*.css" \) \
    -exec sed -i "s/PROJECT_NAME/${NAME}/g" {} +

cd "$DEST"

# Install Tailwind v4
npm install --save-dev tailwindcss @tailwindcss/cli --silent

# Initial build so CSS exists immediately
npx @tailwindcss/cli -i ./src/css/input.css -o ./dist/css/style.css --quiet 2>/dev/null || true

echo "Tailwind v4 ready — run: npm run dev"
