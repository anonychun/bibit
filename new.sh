#!/bin/bash
set -euo pipefail

# Check if project name is provided as argument
if [ $# -eq 0 ]; then
	echo -e "\033[31m❌ Project name required!\033[0m"
	echo ""
	echo -e "\033[1mUsage:\033[0m"
	echo -e "  \033[36mwget -qO- \033[34mhttps://raw.githubusercontent.com/anonychun/bibit/refs/heads/main/new.sh \033[33m| \033[36mbash -s \033[32m<project-name>\033[0m"
	echo ""
	echo -e "\033[1mExample:\033[0m"
	echo -e "  \033[36mwget -qO- \033[34mhttps://raw.githubusercontent.com/anonychun/bibit/refs/heads/main/new.sh \033[33m| \033[36mbash -s \033[32mgithub.com/anonychun/verification-api\033[0m"
	exit 1
fi

PROJECT="$1"

ORIGINAL_DIR=$(pwd)
TMP_DIR=$(mktemp -d)
DIR_NAME=$(basename "$PROJECT")

# Clone to tmp
git clone --depth 1 https://github.com/anonychun/bibit.git "$TMP_DIR"

# Replace project name
cd "$TMP_DIR"
grep -rl "github.com/anonychun/bibit" . | xargs perl -pi -e "s|github.com/anonychun/bibit|$PROJECT|g"

# Copy .env.sample to .env
cp .env.sample .env

# Remove unnecessary files
rm -rf .git new.sh

# Move back to original directory and copy project
cd "$ORIGINAL_DIR"
cp -r "$TMP_DIR" "$DIR_NAME"

# Clean up temporary directory
rm -rf "$TMP_DIR"

echo "✅ Project initialized in: $DIR_NAME"
