#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v npm >/dev/null 2>&1; then
  echo "npm is required but not installed."
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "Go is required but not installed."
  exit 1
fi

npm install
cd "$ROOT_DIR/frontend" && npm install
cd "$ROOT_DIR"

go mod tidy

if [ -d "$ROOT_DIR/prisma" ]; then
  cd "$ROOT_DIR/prisma" && npx prisma generate || true
fi

echo "Bootstrap complete."
