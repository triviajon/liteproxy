#!/usr/bin/env bash
set -euo pipefail

# Source versions
source versions.env
echo "GO_VERSION=$GO_VERSION"

# Extract go.mod version
go_mod_version=$(grep -E '^go [0-9]+\.[0-9]+(\.[0-9]+)?$' processor/go.mod | awk '{print $2}')
if [ -z "$go_mod_version" ]; then
  echo "ERROR: go.mod version not found"
  exit 1
fi
echo "go.mod version=$go_mod_version"

# Extract README go version (support both old and new formatting)
readme_go_version=$(grep -Eo 'GO_VERSION=[0-9]+\.[0-9]+(\.[0-9]+)?' README.md 2>/dev/null || true)
readme_go_version=$(echo "$readme_go_version" | head -n1 | cut -d= -f2)
if [ -z "$readme_go_version" ]; then
  echo "ERROR: README go version not found (expected 'golang version: goX.Y.Z' or 'GO_VERSION=X.Y.Z')"
  exit 1
fi
echo "README golang version=$readme_go_version"

# Compare
if [ "$GO_VERSION" != "$go_mod_version" ]; then
  echo "ERROR: versions.env GO_VERSION ($GO_VERSION) != go.mod go version ($go_mod_version)"
  exit 1
fi

if [ "$GO_VERSION" != "$readme_go_version" ]; then
  echo "ERROR: versions.env GO_VERSION ($GO_VERSION) != README golang version ($readme_go_version)"
  exit 1
fi

echo "Version sync check passed."
