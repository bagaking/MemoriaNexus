#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status.

# This script reads the go.mod file and copies local modules to the appropriate locations in the Docker build context.

# Set the working directory to the script's location
cd "$(dirname "$0")"
cd ..

# Read the replace directives from go.mod
replaces=$(grep -E 'replace .* => \.\./' go.mod || true)

# Copy each local module to the Docker build context
if [ -z "$replaces" ]; then
    echo "No local module replacements found in go.mod"
else
  while IFS= read -r line; do
      # Extract the module path and local path
      module_path=$(echo "$line" | awk '{print $2}')
      local_path=$(echo "$line" | awk '{print $4}')

      # Create the target directory
      target_dir="./deployment/docker-build-context/$(basename "$local_path")"

      echo "local_path: $local_path"
      echo "target_dir: $target_dir"
      mkdir -p "$target_dir"

      # Copy the local module to the target directory, excluding .git directory
      rsync -av --exclude='.git' "$local_path/" "$target_dir/"
  done <<< "$replaces"
fi

cp go.mod ./deployment/docker-build-context/
cp go.sum ./deployment/docker-build-context/

# Replace the local module paths in go.mod
cd ./deployment/docker-build-context
if [ -z "$replaces" ]; then
    echo "No local module replacements found in go.mod"
else
  while IFS= read -r line; do
      module_path=$(echo "$line" | awk '{print $2}')
      local_path=$(echo "$line" | awk '{print $4}')
      new_path="./$(basename "$local_path")"
      echo "-> Replacing $module_path with $new_path"
      go mod edit -replace="${module_path}=${new_path}"
  done <<< "$replaces"
fi

echo "Local modules copied successfully."