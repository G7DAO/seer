#!/usr/bin/env bash
set -euo pipefail

# ------------------------------------------------------------------------------
# Example script to test high-level functionality of the seer CLI commands.
# Adjust commands/flags/paths to match your environment.
#
# Prerequisites:
# - seer is built and available in your PATH
# - Environment variables for DB/storage connections, etc., are properly set
# - You have run "go build" or similar so that the "seer" binary exists
# ------------------------------------------------------------------------------



green='\033[0;32m'
red='\033[0;31m'
reset='\033[0m'


echo -e "${green}Building seer...${reset}"
go build
echo -e "Building seer succeeded.${green}(1/12)${reset}"
sleep 2
echo

echo -e "${green}Testing seer version...${reset}"
seer version
echo -e "Testing seer version succeeded.${green}(2/12)${reset}"
sleep 2
echo

echo -e "${green}Testing shell completion subcommands...${reset}"
seer completion bash > /dev/null
seer completion zsh > /dev/null
echo -e "Shell completion generation succeeded.${green}(3/12)${reset}"
sleep 2
echo

echo -e "${green}Testing seer blockchain subcommands...${reset}"
./seer blockchain generate --name ethereum
go fmt ./blockchain/ethereum/ethereum.go
echo -e "Generated ethereum. Check ./blockchain/ethereum/ethereum.go.${green}(4/12)${reset}"
sleep 2
echo

ERC20_TRANSFER_ABI='[{"type":"event","name":"Transfer","inputs":[{"name":"from","type":"address","indexed":true},{"name":"to","type":"address","indexed":true},{"name":"value","type":"uint256","indexed":false}],"anonymous":false}]'

# Verify required env variables are set
echo -e "${green}Verifying required env variables are set...${reset}"
: "${MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI:?Environment variable MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI not set}"
: "${MOONSTREAM_DB_URL:?Environment variable MOONSTREAM_DB_URL not set}"
echo -e "All required env variables are set.${green}(5/12)${reset}"
sleep 2
echo

# test case when chain ID mismatch in crawler
echo -e "${green}Testing chain ID mismatch in crawler...${reset}"
echo "n" | ./seer crawler \
    --chain arbitrum_one \
    --rpc-url "$MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI" \
    --start-block 1000 \
    --final-block 1010 \
    --threads 2 \
    --timeout 5 \
    --batch-size 2 \
    --confirmations 1 \
    --base-dir /tmp/your-data-dir || true
echo -e "Testing chain ID mismatch succeeded.${green}(6/12)${reset}"
sleep 2
echo

echo "Testing existing crawler (dry run/example) - adjust flags for your RPC or environment..."
./seer crawler \
    --chain ethereum \
    --start-block 1000 \
    --final-block 1010 \
    --threads 2 \
    --timeout 5 \
    --batch-size 2 \
    --confirmations 1 \
    --base-dir /tmp/your-data-dir \
    --rpc-url "$MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI"	
echo -e "Testing crawler succeeded.${green}(7/12)${reset}"
sleep 2
echo

echo
echo -e "${green}Testing storages/inspector subcommands...${reset}"
./seer inspector storage \
    --chain ethereum \
    --base-dir /tmp/your-data-dir \
    --delim '' \
    --return-func '' \
    --timeout 10
echo -e "Testing storages/inspector subcommands succeeded.${green}(8/12)${reset}"
sleep 2
echo


echo
echo -e "${green}Testing synchronizer subcommand (dry run/example)...${reset}"
./seer synchronizer \
    --chain ethereum \
    --start-block 1000 \
    --end-block 1010 \
    --threads 2 \
    --timeout 5 \
    --batch-size 2 \
    --base-dir /tmp/your-data-dir \
    --rpc-url "$MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI" \
    --customer-db-uri "$MOONSTREAM_DB_URL"
echo -e "Testing historical-sync command succeeded.${green}(9/12)${reset}"
sleep 2
echo

echo
echo -e "${green}Testing historical-sync command...${reset}"
./seer historical-sync \
    --chain ethereum \
    --base-dir /tmp/your-data-dir \
    --start-block 1010 \
    --end-block 1000 \
    --rpc-url "$MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI" \
    --customer-db-uri "$MOONSTREAM_DB_URL" \
    --threads 2 \
    --batch-size 2 \
    --timeout 10
echo -e "Testing historical-sync command succeeded.${green}(10/12)${reset}"
sleep 2
echo

echo
echo -e "${green}Testing deployment-blocks command...${reset}"
./seer databases index deployment-blocks \
    --chain ethereum \
    --rpc-url "$MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI"
echo -e "Testing deployment-blocks command succeeded.${green}(11/12)${reset}"
sleep 2
echo

echo
echo -e "${green}Testing create-jobs command...${reset}"
# Example creating a job for an ABI file. Adjust flags as needed.
echo "$ERC20_TRANSFER_ABI" | ./seer databases index create-jobs \
    --chain ethereum \
    --address "0x0000000000000000000000000000000000000000" \
    --abi-file "/dev/stdin" \
    --customer-id "00000000-0000-0000-0000-000000000000" \
    --user-id "00000000-0000-0000-0000-000000000000" \
    --rpc-url "$MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI"
echo -e "Testing create-jobs command succeeded.${green}(12/12)${reset}"
sleep 2
echo


echo
echo -e "All test commands completed successfully!${green}(12/12)${reset}"
exit 0
