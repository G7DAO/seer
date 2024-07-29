#!/usr/bin/env bash

# Deployment script of monitoring services

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
SECRETS_DIR="${SECRETS_DIR:-/home/ubuntu/seer-secrets}"
PARAMETERS_ENV_PATH="${SECRETS_DIR}/monitoring.env"
SCRIPT_DIR="$(realpath $(dirname $0))"
USER_SYSTEMD_DIR="${USER_SYSTEMD_DIR:-/home/ubuntu/.config/systemd/user}"

# Service files
MONITORING_SEER_CRAWLERS_SERVICE_FILE="monitoring-seer-crawlers.service"

set -eu

echo
echo
echo -e "${PREFIX_INFO} Copy monitoring binary from GCP Storage"
gcloud storage cp gs://moonstream-binaries/prod/monitoring/monitoring "/home/ubuntu/monitoring"
chmod +x "/home/ubuntu/monitoring"
chown ubuntu:ubuntu "/home/ubuntu/monitoring"

echo
echo
echo -e "${PREFIX_INFO} Prepare secrets directory if required"
if [ ! -d "${SECRETS_DIR}" ]; then
  mkdir "${SECRETS_DIR}"
  echo -e "${PREFIX_WARN} Created new secrets directory"
fi

echo
echo
echo -e "${PREFIX_INFO} Retrieving deployment parameters"
echo "# Monitoring environment variables" > "${PARAMETERS_ENV_PATH}"

MONITORING_SECRETS_LIST=$(gcloud secrets list --filter 'labels.monitoring:true' --format 'value(name)')
for s in $MONITORING_SECRETS_LIST; do
  secret_value=$(gcloud secrets versions access latest --secret=$s)
  echo "${s}=${secret_value}" >> "${PARAMETERS_ENV_PATH}"
done

chmod 0640 "${PARAMETERS_ENV_PATH}"

echo
echo
echo -e "${PREFIX_INFO} Prepare monitoring configuration"
if [ ! -d "/home/ubuntu/.monitoring" ]; then
  mkdir -p /home/ubuntu/.monitoring
  echo -e "${PREFIX_WARN} Created monitoring configuration directory"
fi
cp "${SCRIPT_DIR}/monitoring-seer-crawlers-config.json" /home/ubuntu/.monitoring/monitoring-seer-crawlers-config.json

echo
echo
echo -e "${PREFIX_INFO} Prepare user systemd directory if required"
if [ ! -d "${USER_SYSTEMD_DIR}" ]; then
  mkdir -p "${USER_SYSTEMD_DIR}"
  echo -e "${PREFIX_WARN} Created new user systemd directory"
fi

echo
echo
echo -e "${PREFIX_INFO} Replacing existing monitoring seer crawler service definition with ${MONITORING_SEER_CRAWLERS_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${MONITORING_SEER_CRAWLERS_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${MONITORING_SEER_CRAWLERS_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${MONITORING_SEER_CRAWLERS_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${MONITORING_SEER_CRAWLERS_SERVICE_FILE}"
