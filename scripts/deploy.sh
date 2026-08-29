#!/bin/bash
#
# Auto-deploy script for reminder.
# Checks for new GitHub releases and deploys automatically.
#
# Usage:
#   ./deploy.sh            - Check for new release and deploy if available
#   ./deploy.sh --force    - Force re-deploy the latest release
#   ./deploy.sh --cron     - Install a cron job to run every 5 minutes
#   ./deploy.sh --install  - Install systemd service + deploy timer
#
# By default, the binary, .env, and data are expected in the same directory
# as this script. Set REMINDER_HOME to override (e.g. REMINDER_HOME=/root).
# Stores the currently running version in .version file.

set -euo pipefail

REPO="asim/reminder"
ARCH="linux_amd64"
DEPLOY_DIR="$(cd "$(dirname "$0")" && pwd)"
# HOME_DIR is where the binary, .env, and data live — override with REMINDER_HOME
HOME_DIR="${REMINDER_HOME:-$DEPLOY_DIR}"
VERSION_FILE="$HOME_DIR/.version"
LOG_FILE="$HOME_DIR/reminder.log"
BINARY="$HOME_DIR/reminder"
SERVICE_NAME="reminder"

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

# Check if systemd is managing the service
has_systemd_service() {
    systemctl is-enabled "$SERVICE_NAME" >/dev/null 2>&1
}

deploy_version() {
    local version="$1"
    local file="reminder_${version}_${ARCH}.tar.gz"
    local url="https://github.com/$REPO/releases/download/v${version}/${file}"

    log "Downloading reminder v${version}"
    cd "$HOME_DIR"

    if ! wget -q "$url" -O "$file"; then
        log "ERROR: Failed to download $url"
        rm -f "$file"
        return 1
    fi

    log "Extracting $file"
    tar zxf "$file"

    if [ ! -f "$HOME_DIR/reminder" ]; then
        log "ERROR: Binary not found after extraction"
        rm -f "$file"
        return 1
    fi

    chmod +x "$HOME_DIR/reminder"

    if has_systemd_service; then
        log "Restarting via systemd"
        sudo systemctl restart "$SERVICE_NAME"
    else
        log "Stopping current instance"
        killall reminder 2>/dev/null || true
        sleep 3

        log "Starting reminder v${version}"
        if [ -f "$HOME_DIR/.env" ]; then
            set -a
            . "$HOME_DIR/.env"
            set +a
        fi

        nohup "$BINARY" --serve --web >> "$LOG_FILE" 2>&1 &
        disown
    fi

    echo "$version" > "$VERSION_FILE"
    rm -f "$file"

    log "Deployed reminder v${version}"
}

install_systemd() {
    # Write service unit inline — no external file needed
    sudo tee /etc/systemd/system/reminder.service >/dev/null <<EOF
[Unit]
Description=Reminder - Quran, Hadith & Names of Allah
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=$HOME_DIR
ExecStart=/bin/bash -c 'set -a; . $HOME_DIR/.env; set +a; exec $HOME_DIR/reminder --serve --web'
Restart=on-failure
RestartSec=5
StandardOutput=append:$HOME_DIR/reminder.log
StandardError=append:$HOME_DIR/reminder.log
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
EOF

    # Install deploy timer (runs every 5 minutes instead of cron)
    sudo tee /etc/systemd/system/reminder-deploy.service >/dev/null <<EOF
[Unit]
Description=Check for new reminder releases

[Service]
Type=oneshot
ExecStart=$DEPLOY_DIR/deploy.sh
StandardOutput=append:$DEPLOY_DIR/deploy.log
StandardError=append:$DEPLOY_DIR/deploy.log
EOF

    sudo tee /etc/systemd/system/reminder-deploy.timer >/dev/null <<EOF
[Unit]
Description=Check for new reminder releases every 5 minutes

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min

[Install]
WantedBy=timers.target
EOF

    sudo systemctl daemon-reload

    # Stop any manually-running instance
    killall reminder 2>/dev/null || true
    sleep 2

    # Enable and start
    sudo systemctl enable --now reminder
    sudo systemctl enable --now reminder-deploy.timer

    log "Installed and started:"
    log "  reminder.service        - starts on boot, restarts on crash"
    log "  reminder-deploy.timer   - checks for updates every 5 min"
    log ""
    log "Useful commands:"
    log "  sudo systemctl status reminder         - check status"
    log "  sudo journalctl -u reminder -f         - follow logs"
    log "  sudo systemctl restart reminder         - restart"
    log "  sudo systemctl stop reminder            - stop"
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

# Handle flags
case "${1:-}" in
    --install)
        install_systemd
        exit 0
        ;;
    --cron)
        install_cron
        exit 0
        ;;
esac

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
