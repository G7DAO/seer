#!/usr/bin/env bash

set -eu

INSTANCE_APP_NAME="seer"

SYSTEM_USER="${SYSTEM_USER:-ubuntu}"
AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-us-east-1}"
DEPLOYMENT_LOG_FILE_PATH="/home/$SYSTEM_USER/deployment.log"
DEPLOY_SCRIPT_NAME="deploy.bash"

echo "Retrieve latest deployment script" >> "$DEPLOYMENT_LOG_FILE_PATH"
SCRIPT_FILE_PATH="/home/$SYSTEM_USER/$DEPLOY_SCRIPT_NAME"
AWS_DEFAULT_REGION="$AWS_DEFAULT_REGION" aws ssm get-parameter --name "/wb/deployment/$INSTANCE_APP_NAME/$DEPLOY_SCRIPT_NAME" --with-decryption --query "Parameter.Value" --output text > "$SCRIPT_FILE_PATH"
chmod +x "$SCRIPT_FILE_PATH"

echo "Execute $DEPLOY_SCRIPT_NAME script" >> "$DEPLOYMENT_LOG_FILE_PATH"
/home/$SYSTEM_USER/$DEPLOY_SCRIPT_NAME
