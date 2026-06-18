#!/bin/bash
set -e

APP="ivr-atljanseva"
DEPLOY_DIR="/opt/$APP"
SERVICE_FILE="deploy/$APP.service"

echo "Building binary..."
go build -o "$APP" .

echo "Stopping service..."
sudo systemctl stop "$APP" 2>/dev/null || true

echo "Copying files..."
sudo mkdir -p "$DEPLOY_DIR"
sudo cp "$APP" "$DEPLOY_DIR/"
sudo cp -r audio "$DEPLOY_DIR/"
sudo cp .env "$DEPLOY_DIR/" 2>/dev/null || echo "Warning: .env not found"

echo "Installing systemd service..."
sudo cp "$SERVICE_FILE" /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable "$APP"
sudo systemctl start "$APP"

echo "Deployed! Status:"
sudo systemctl status "$APP" --no-pager