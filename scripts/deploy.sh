#!/bin/bash
#
# Auto-deploy script for reminder.
# Checks for new GitHub releases and deploys automatically.
#
# Usage:
#   ./deploy.sh          - Check for new release and deploy if available
#   ./deploy.sh --force  - Force re-deploy the latest release
#   ./deploy.sh --cron   - Install a cron job to run every 5 minutes
#
# Expects a .env file in the same directory with environment variables.
# Stores the currently running version in .version file.

set -euo pipefail

REPO="asim/reminder"
ARCH="linux_amd64"
DEPLOY_DIR="$(cd "$(dirname "$0")" && pwd)"
VERSION_FILE="$DEPLOY_DIR/.version"
LOG_FILE="$DEPLOY_DIR/reminder.log"
BINARY="$DEPLOY_DIR/reminder"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"
}

get_latest_release() {
    curl -s "https://api.github.com/repos/$REPO/releases/latest" \
        | grep '"tag_name"' \
        | head -1 \
        | sed 's/.*"tag_name": *"v\([^"]*\)".*/\1/'
}

get_current_version() {
    if [ -f "$VERSION_FILE" ]; then
        cat "$VERSION_FILE"
    else
        echo ""
    fi
}

deploy_version() {
    local version="$1"
    local file="reminder_${version}_${ARCH}.tar.gz"
    local url="https://github.com/$REPO/releases/download/v${version}/${file}"

    log "Downloading reminder v${version}"
    cd "$DEPLOY_DIR"

    if ! wget -q "$url" -O "$file"; then
        log "ERROR: Failed to download $url"
        rm -f "$file"
        return 1
    fi

    log "Extracting $file"
    tar zxf "$file"

    if [ ! -f "$DEPLOY_DIR/reminder" ]; then
        log "ERROR: Binary not found after extraction"
        rm -f "$file"
        return 1
    fi

    chmod +x "$DEPLOY_DIR/reminder"

    log "Stopping current instance"
    killall reminder 2>/dev/null || true
    sleep 3

    log "Starting reminder v${version}"
    if [ -f "$DEPLOY_DIR/.env" ]; then
        set -a
        . "$DEPLOY_DIR/.env"
        set +a
    fi

    nohup "$BINARY" --serve --web >> "$LOG_FILE" 2>&1 &
    disown

    echo "$version" > "$VERSION_FILE"
    rm -f "$file"

    log "Deployed reminder v${version} (pid $!)"
}

install_cron() {
    local script="$DEPLOY_DIR/deploy.sh"
    local job="*/5 * * * * $script >> $DEPLOY_DIR/deploy.log 2>&1"

    # Remove any existing reminder deploy cron entries
    crontab -l 2>/dev/null | grep -v "$script" | crontab - 2>/dev/null || true

    # Add new entry
    (crontab -l 2>/dev/null; echo "$job") | crontab -

    log "Installed cron job: $job"
    log "Deploy logs will go to $DEPLOY_DIR/deploy.log"
}

# Handle --cron flag
if [ "${1:-}" = "--cron" ]; then
    install_cron
    exit 0
fi

FORCE=false
if [ "${1:-}" = "--force" ]; then
    FORCE=true
fi

latest=$(get_latest_release)
current=$(get_current_version)

if [ -z "$latest" ]; then
    log "ERROR: Could not determine latest release"
    exit 1
fi

if [ "$FORCE" = true ] || [ "$latest" != "$current" ]; then
    if [ -n "$current" ]; then
        log "Upgrading from v${current} to v${latest}"
    else
        log "Installing v${latest}"
    fi
    deploy_version "$latest"
else
    log "Already running v${current}, no update needed"
fi
