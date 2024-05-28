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

# Service files
SEER_CRAWLER_ETHEREUM_SERVICE_FILE="seer-crawler-ethereum.service"
SEER_CRAWLER_POLYGON_SERVICE_FILE="seer-crawler-polygon.service"

set -eu

if [ ! -d "${SECRETS_DIR}" ]; then
  mkdir "${SECRETS_DIR}"
  echo -e "${PREFIX_WARN} Created new secrets directory"
fi

echo
echo
echo -e "${PREFIX_INFO} Retrieving deployment parameters"
echo "# Seer environment variables" > "${PARAMETERS_ENV_PATH}"
MOONSTREAM_DB_V3_INDEXES_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_DB_V3_INDEXES_URI)
echo "MOONSTREAM_DB_V3_INDEXES_URI=${MOONSTREAM_DB_V3_INDEXES_URI}" >> "${PARAMETERS_ENV_PATH}"
MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI=${MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"
MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI=${MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"
SEER_CRAWLER_INDEXER_LABEL=$(gcloud secrets versions access latest --secret=SEER_CRAWLER_INDEXER_LABEL)
echo "SEER_CRAWLER_INDEXER_LABEL=${SEER_CRAWLER_INDEXER_LABEL}" >> "${PARAMETERS_ENV_PATH}"
SEER_CRAWLER_STORAGE_TYPE=$(gcloud secrets versions access latest --secret=SEER_CRAWLER_STORAGE_TYPE)
echo "SEER_CRAWLER_STORAGE_TYPE=${SEER_CRAWLER_STORAGE_TYPE}" >> "${PARAMETERS_ENV_PATH}"
SEER_CRAWLER_STORAGE_BUCKET=$(gcloud secrets versions access latest --secret=SEER_CRAWLER_STORAGE_BUCKET)
echo "SEER_CRAWLER_STORAGE_BUCKET=${SEER_CRAWLER_STORAGE_BUCKET}" >> "${PARAMETERS_ENV_PATH}"
SEER_CRAWLER_STORAGE_PREFIX=$(gcloud secrets versions access latest --secret=SEER_CRAWLER_STORAGE_PREFIX)
echo "SEER_CRAWLER_STORAGE_PREFIX=${SEER_CRAWLER_STORAGE_PREFIX}" >> "${PARAMETERS_ENV_PATH}"
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

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for ethereum blockchain service definition with ${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for polygon blockchain service definition with ${SEER_CRAWLER_POLYGON_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_POLYGON_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_POLYGON_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_POLYGON_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_POLYGON_SERVICE_FILE}"
