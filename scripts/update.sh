#!/bin/bash
set -euo pipefail

# Klikku Deploy Script
# Triggered by GitHub webhook on push to main

PROJECT_DIR="/root/hermes/klikku"
cd "$PROJECT_DIR"

echo "=== Pulling latest code ==="
git pull origin main

echo "=== Building and deploying ==="
docker compose build --no-cache backend frontend
docker compose up -d --force-recreate

echo "=== Cleaning up ==="
docker system prune -f

echo "=== Done at $(date) ==="
