#!/usr/bin/env bash

# Deployment script

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
AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-us-east-1}"
APP_DIR="${APP_DIR:-/home/ubuntu/seer}"
SECRETS_DIR="${SECRETS_DIR:-/home/ubuntu/seer-secrets}"
PARAMETERS_ENV_PATH="${SECRETS_DIR}/app.env"
SCRIPT_DIR="$(realpath $(dirname $0))"
USER_SYSTEMD_DIR="${USER_SYSTEMD_DIR:-/home/ubuntu/.config/systemd/user}"

set -eu

if [ ! -d "${SECRETS_DIR}" ]; then
  mkdir "${SECRETS_DIR}"
  echo -e "${PREFIX_WARN} Created new secrets directory"
fi

echo
echo
echo -e "${PREFIX_INFO} Retrieving deployment parameters"
echo "# Seer environment variables" > "${PARAMETERS_ENV_PATH}"

# List of secrets to retrieve
SECRETS=(
  "MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN"
  "MOONSTREAM_DB_V3_INDEXES_URI"
  "MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_GAME7_TESTNET_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_XAI_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_IMX_ZKEVM_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_IMX_ZKEVM_SEPOLIA_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_B3_A_EXTERNAL_URI"
  "MOONSTREAM_NODE_B3_SEPOLIA_A_EXTERNAL_URI"
  "SEER_CRAWLER_STORAGE_BUCKET"
)

for SECRET in "${SECRETS[@]}"; do
  VALUE=$(gcloud secrets versions access latest --secret="${SECRET}")
  echo "${SECRET}=${VALUE}" >> "${PARAMETERS_ENV_PATH}"
done

echo "SEER_CRAWLER_INDEXER_LABEL=seer" >> "${PARAMETERS_ENV_PATH}"
echo "SEER_CRAWLER_STORAGE_TYPE=gcp-storage" >> "${PARAMETERS_ENV_PATH}"
echo "SEER_CRAWLER_STORAGE_PREFIX=prod" >> "${PARAMETERS_ENV_PATH}"
chmod 0640 "${PARAMETERS_ENV_PATH}"

echo
echo
echo -e "${PREFIX_INFO} Build seer binary"
EXEC_DIR=$(pwd)
cd "${APP_DIR}"
HOME=/home/ubuntu /usr/local/go/bin/go build -o "${APP_DIR}/seer" .
chmod +x "${APP_DIR}/seer"
chown ubuntu:ubuntu "${APP_DIR}/seer"
cd "${EXEC_DIR}"

echo
echo
echo -e "${PREFIX_INFO} Prepare user systemd directory"
if [ ! -d "${USER_SYSTEMD_DIR}" ]; then
  mkdir -p "${USER_SYSTEMD_DIR}"
  echo -e "${PREFIX_WARN} Created new user systemd directory"
fi

# Function to update service files
update_service() {
  local SERVICE_FILE=$1
  if [ -f "${SCRIPT_DIR}/${SERVICE_FILE}" ]; then
    echo
    echo
    echo -e "${PREFIX_INFO} Replacing existing service with ${SERVICE_FILE}"
    chmod 644 "${SCRIPT_DIR}/${SERVICE_FILE}"
    cp "${SCRIPT_DIR}/${SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SERVICE_FILE}"
    XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
    XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SERVICE_FILE}"
  else
    echo -e "${PREFIX_WARN} Service file ${SERVICE_FILE} not found, skipping"
  fi
}

# Function to update historical synchronizer services and timers
update_historical_service() {
  local SERVICE_FILE=$1
  local TIMER_FILE=$2
  if [ -f "${SCRIPT_DIR}/${SERVICE_FILE}" ] && [ -f "${SCRIPT_DIR}/${TIMER_FILE}" ]; then
    echo
    echo
    echo -e "${PREFIX_INFO} Replacing existing historical synchronizer with ${SERVICE_FILE} and ${TIMER_FILE}"
    chmod 644 "${SCRIPT_DIR}/${SERVICE_FILE}" "${SCRIPT_DIR}/${TIMER_FILE}"
    cp "${SCRIPT_DIR}/${SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SERVICE_FILE}"
    cp "${SCRIPT_DIR}/${TIMER_FILE}" "${USER_SYSTEMD_DIR}/${TIMER_FILE}"
    XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
    XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${TIMER_FILE}"
  else
    echo -e "${PREFIX_WARN} Service or timer file ${SERVICE_FILE} or ${TIMER_FILE} not found, skipping"
  fi
}

# Update crawler services
for SERVICE_FILE in "${SCRIPT_DIR}"/seer-crawler-*.service; do
  [ -e "$SERVICE_FILE" ] || continue
  SERVICE_FILE=$(basename "${SERVICE_FILE}")
  update_service "${SERVICE_FILE}"
done

# Update synchronizer services
for SERVICE_FILE in "${SCRIPT_DIR}"/seer-synchronizer-*.service; do
  [ -e "$SERVICE_FILE" ] || continue
  SERVICE_FILE=$(basename "${SERVICE_FILE}")
  update_service "${SERVICE_FILE}"
done

# Update historical synchronizer services and timers
for SERVICE_FILE in "${SCRIPT_DIR}"/seer-historical-synchronizer-*.service; do
  [ -e "$SERVICE_FILE" ] || continue
  SERVICE_FILE=$(basename "${SERVICE_FILE}")
  TIMER_FILE="${SERVICE_FILE%.service}.timer"
  update_historical_service "${SERVICE_FILE}" "${TIMER_FILE}"
done
