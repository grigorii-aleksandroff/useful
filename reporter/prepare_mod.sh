#!/usr/bin/env bash
set -e

if [ -f ".env" ]; then
  export $(grep -v '^#' .env | xargs)
fi

GO_MOD_PATH="go.mod"
PREFIX=""
if [ "$1" = "test" ]; then
  GO_MOD_PATH="test/go.mod"
  mkdir -p "$(dirname "$GO_MOD_PATH")"
  PREFIX="../"
fi

GO_VERSION=${GO_VERSION:-1.25.1}

echo "module git.bububla.com/kilogramix/asia/reporter.git" > "$GO_MOD_PATH"
echo "" >> "$GO_MOD_PATH"

if [ -n "$CONTRACTS_PATH" ]; then
  echo "replace git.bububla.com/kilogramix/asia/contracts.git => ${PREFIX}${CONTRACTS_PATH}" >> "$GO_MOD_PATH"
fi

if [ -n "$PLATFORM_PATH" ]; then
  echo "replace git.itechpsp.com/e46/box/platform.git => ${PREFIX}${PLATFORM_PATH}" >> "$GO_MOD_PATH"
fi

echo "" >> "$GO_MOD_PATH"
cat go.mod.require >> "$GO_MOD_PATH"

echo "✅ $GO_MOD_PATH сгенерирован"
