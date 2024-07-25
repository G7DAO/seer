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
# Crawler
SEER_CRAWLER_ARBITRUM_ONE_SERVICE_FILE="seer-crawler-arbitrum-one.service"
SEER_CRAWLER_ARBITRUM_SEPOLIA_SERVICE_FILE="seer-crawler-arbitrum-sepolia.service"
SEER_CRAWLER_ETHEREUM_SERVICE_FILE="seer-crawler-ethereum.service"
SEER_CRAWLER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE="seer-crawler-game7-orbit-arbitrum-sepolia.service"
SEER_CRAWLER_MANTLE_SEPOLIA_SERVICE_FILE="seer-crawler-mantle-sepolia.service"
SEER_CRAWLER_MANTLE_SERVICE_FILE="seer-crawler-mantle.service"
SEER_CRAWLER_POLYGON_SERVICE_FILE="seer-crawler-polygon.service"
SEER_CRAWLER_XAI_SEPOLIA_SERVICE_FILE="seer-crawler-xai-sepolia.service"
SEER_CRAWLER_XAI_SERVICE_FILE="seer-crawler-xai.service"
SEER_CRAWLER_SEPOLIA_SERVICE_FILE="seer-crawler-sepolia.service"
SEER_CRAWLER_IMX_ZKEVM_SERVICE_FILE="seer-crawler-imx-zkevm.service"
SEER_CRAWLER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE="seer-crawler-imx-zkevm-sepolia.service"

# Synchronizer
SEER_SYNCHRONIZER_ETHEREUM_SERVICE_FILE="seer-synchronizer-ethereum.service"
SEER_SYNCHRONIZER_POLYGON_SERVICE_FILE="seer-synchronizer-polygon.service"
SEER_SYNCHRONIZER_ARBITRUM_ONE_SERVICE_FILE="seer-synchronizer-arbitrum-one.service"
SEER_SYNCHRONIZER_ARBITRUM_SEPOLIA_SERVICE_FILE="seer-synchronizer-arbitrum-sepolia.service"
SEER_SYNCHRONIZER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE="seer-synchronizer-game7-orbit-arbitrum-sepolia.service"
SEER_SYNCHRONIZER_MANTLE_SEPOLIA_SERVICE_FILE="seer-synchronizer-mantle-sepolia.service"
SEER_SYNCHRONIZER_MANTLE_SERVICE_FILE="seer-synchronizer-mantle.service"
SEER_SYNCHRONIZER_XAI_SEPOLIA_SERVICE_FILE="seer-synchronizer-xai-sepolia.service"
SEER_SYNCHRONIZER_XAI_SERVICE_FILE="seer-synchronizer-xai.service"
SEER_SYNCHRONIZER_SEPOLIA_SERVICE_FILE="seer-synchronizer-sepolia.service"
SEER_SYNCHRONIZER_IMX_ZKEVM_SERVICE_FILE="seer-synchronizer-imx-zkevm.service"
SEER_SYNCHRONIZER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE="seer-synchronizer-imx-zkevm-sepolia.service"

set -eu

if [ ! -d "${SECRETS_DIR}" ]; then
  mkdir "${SECRETS_DIR}"
  echo -e "${PREFIX_WARN} Created new secrets directory"
fi

echo
echo
echo -e "${PREFIX_INFO} Retrieving deployment parameters"
echo "# Seer environment variables" > "${PARAMETERS_ENV_PATH}"

MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN=$(gcloud secrets versions access latest --secret=MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN)
echo "MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN=${MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_DB_V3_INDEXES_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_DB_V3_INDEXES_URI)
echo "MOONSTREAM_DB_V3_INDEXES_URI=${MOONSTREAM_DB_V3_INDEXES_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI=${MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI=${MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI=${MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI=${MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI=${MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI=${MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_XAI_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_XAI_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_XAI_A_EXTERNAL_URI=${MOONSTREAM_NODE_XAI_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI=${MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI=${MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI=${MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_IMX_ZKEVM_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_IMX_ZKEVM_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_IMX_ZKEVM_A_EXTERNAL_URI=${MOONSTREAM_NODE_IMX_ZKEVM_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

MOONSTREAM_NODE_IMX_ZKEVM_SEPOLIA_A_EXTERNAL_URI=$(gcloud secrets versions access latest --secret=MOONSTREAM_NODE_IMX_ZKEVM_SEPOLIA_A_EXTERNAL_URI)
echo "MOONSTREAM_NODE_IMX_ZKEVM_SEPOLIA_A_EXTERNAL_URI=${MOONSTREAM_NODE_IMX_ZKEVM_SEPOLIA_A_EXTERNAL_URI}" >> "${PARAMETERS_ENV_PATH}"

echo "SEER_CRAWLER_INDEXER_LABEL=seer" >> "${PARAMETERS_ENV_PATH}"

echo "SEER_CRAWLER_STORAGE_TYPE=gcp-storage" >> "${PARAMETERS_ENV_PATH}"

SEER_CRAWLER_STORAGE_BUCKET=$(gcloud secrets versions access latest --secret=SEER_CRAWLER_STORAGE_BUCKET)
echo "SEER_CRAWLER_STORAGE_BUCKET=${SEER_CRAWLER_STORAGE_BUCKET}" >> "${PARAMETERS_ENV_PATH}"

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

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Arbitrum One blockchain service definition with ${SEER_CRAWLER_ARBITRUM_ONE_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_ARBITRUM_ONE_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_ARBITRUM_ONE_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_ARBITRUM_ONE_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_ARBITRUM_ONE_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Arbitrum Sepolia blockchain service definition with ${SEER_CRAWLER_ARBITRUM_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_ARBITRUM_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_ARBITRUM_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_ARBITRUM_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_ARBITRUM_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Ethereum blockchain service definition with ${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_ETHEREUM_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Game7 Orbit Arbitrum Sepolia blockchain service definition with ${SEER_CRAWLER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Mantle Sepolia blockchain service definition with ${SEER_CRAWLER_MANTLE_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_MANTLE_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_MANTLE_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_MANTLE_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_MANTLE_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Mantle blockchain service definition with ${SEER_CRAWLER_MANTLE_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_MANTLE_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_MANTLE_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_MANTLE_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_MANTLE_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for polygon blockchain service definition with ${SEER_CRAWLER_POLYGON_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_POLYGON_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_POLYGON_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_POLYGON_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_POLYGON_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Xai Sepolia blockchain service definition with ${SEER_CRAWLER_XAI_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_XAI_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_XAI_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_XAI_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_XAI_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Xai blockchain service definition with ${SEER_CRAWLER_XAI_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_XAI_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_XAI_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_XAI_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_XAI_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Sepolia blockchain service definition with ${SEER_CRAWLER_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Immutable zkEvm blockchain service definition with ${SEER_CRAWLER_IMX_ZKEVM_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_IMX_ZKEVM_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_IMX_ZKEVM_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_IMX_ZKEVM_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_IMX_ZKEVM_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer crawler for Immutable zkEvm Sepolia blockchain service definition with ${SEER_CRAWLER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_CRAWLER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_CRAWLER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_CRAWLER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_CRAWLER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}"

# Synchronizers

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Ethereum blockchain service definition with ${SEER_SYNCHRONIZER_ETHEREUM_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_ETHEREUM_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_ETHEREUM_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_ETHEREUM_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_ETHEREUM_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Polygon blockchain service definition with ${SEER_SYNCHRONIZER_POLYGON_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_POLYGON_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_POLYGON_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_POLYGON_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_POLYGON_SERVICE_FILE}"



echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Arbitrum One blockchain service definition with ${SEER_SYNCHRONIZER_ARBITRUM_ONE_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_ARBITRUM_ONE_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_ARBITRUM_ONE_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_ARBITRUM_ONE_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_ARBITRUM_ONE_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Arbitrum Sepolia blockchain service definition with ${SEER_SYNCHRONIZER_ARBITRUM_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_ARBITRUM_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_ARBITRUM_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_ARBITRUM_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_ARBITRUM_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Game7 Orbit Arbitrum Sepolia blockchain service definition with ${SEER_SYNCHRONIZER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_GAME7_ORBIT_ARBITRUM_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Mantle Sepolia blockchain service definition with ${SEER_SYNCHRONIZER_MANTLE_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_MANTLE_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_MANTLE_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_MANTLE_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_MANTLE_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Mantle blockchain service definition with ${SEER_SYNCHRONIZER_MANTLE_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_MANTLE_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_MANTLE_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_MANTLE_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_MANTLE_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Xai Sepolia blockchain service definition with ${SEER_SYNCHRONIZER_XAI_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_XAI_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_XAI_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_XAI_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_XAI_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Xai blockchain service definition with ${SEER_SYNCHRONIZER_XAI_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_XAI_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_XAI_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_XAI_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_XAI_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Sepolia blockchain service definition with ${SEER_SYNCHRONIZER_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_SEPOLIA_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Immutable zkEvm blockchain service definition with ${SEER_SYNCHRONIZER_IMX_ZKEVM_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_IMX_ZKEVM_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_IMX_ZKEVM_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_IMX_ZKEVM_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_IMX_ZKEVM_SERVICE_FILE}"

echo
echo
echo -e "${PREFIX_INFO} Replacing existing seer synchronizer for Immutable zkEvm Sepolia blockchain service definition with ${SEER_SYNCHRONIZER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}"
chmod 644 "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}"
cp "${SCRIPT_DIR}/${SEER_SYNCHRONIZER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}" "${USER_SYSTEMD_DIR}/${SEER_SYNCHRONIZER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}"
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user daemon-reload
XDG_RUNTIME_DIR="/run/user/1000" systemctl --user restart "${SEER_SYNCHRONIZER_IMX_ZKEVM_SEPOLIA_SERVICE_FILE}"
