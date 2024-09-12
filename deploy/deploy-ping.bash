#!/usr/bin/env bash

# Deployment script of ping services

# Colors
C_RESET='\033[0m'
C_RED='\033[1;31m'
C_GREEN='\033[1;32m'
C_YELLOW='\033[1;33m'

# Logs
PREFIX_INFO="${C_GREEN}[INFO]${C_RESET} [$(date +%d-%m\ %T)]"
PREFIX_WARN="${C_YELLOW}[WARN]${C_RESET} [$(date +%d-%m\ %T)]"
PREFIX_CRIT="${C_RED}[CRIT]${C_RESET} [$(date +%d-%m\ %T)]"

# Main
SCRIPT_DIR="$(realpath $(dirname $0))"
USER_SYSTEMD_DIR="${USER_SYSTEMD_DIR:-/home/ubuntu/.config/systemd/user}"

# Service files
PING_SERVICE_FILE="ping.service"

set -eu

echo
echo
echo -e "${PREFIX_INFO} Install ping"
HOME=/home/ubuntu /usr/local/go/bin/go install github.com/kompotkot/ping@latest

echo
echo
echo -e "${PREFIX_INFO} Prepare user systemd directory if required"
if [ ! -d "${USER_SYSTEMD_DIR}" ]; then
  mkdir -p "${USER_SYSTEMD_DIR}"
  echo -e "${PREFIX_WARN} Created new user systemd directory"
fi

echo
echo
echo -e "${PREFIX_INFO} Replacing existing ping service definition with ${PING_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${PING_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${PING_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${PING_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${PING_SERVICE_FILE}"
