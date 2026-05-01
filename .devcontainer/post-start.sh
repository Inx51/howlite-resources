#!/usr/bin/env bash
set -euo pipefail

PACKAGES=(
  "libzmq3-dev"
  "libczmq-dev"
)

MISSING=()
for package in "${PACKAGES[@]}"; do
  if ! dpkg -s "$package" >/dev/null 2>&1; then
    MISSING+=("$package")
  fi
done

if [[ ${#MISSING[@]} -eq 0 ]]; then
  echo "All packages are already installed."
  exit 0
fi

echo "Installing: ${MISSING[*]}"
sudo apt-get update
sudo apt-get install -y "${MISSING[@]}"
