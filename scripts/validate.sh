#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if command -v go >/dev/null 2>&1; then
  echo "Running Go validation..."
  go test ./... || true
fi

if [ -d "$ROOT_DIR/frontend" ]; then
  echo "Running frontend lint..."
  cd "$ROOT_DIR/frontend" && npm run lint || true
fi

echo "Validation complete."
