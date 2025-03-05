#!/usr/bin/env bash

# ------------------------------------------------------------------------------
# Test script for the generate-json CLI command.
#
# This script tests the following:
# 1. The CLI command generates the correct Go structs.
# 2. The CLI command creates the expected output file.
# 3. The content of the output file matches the expected JSON.
#
# Prerequisites:
# - The Go program is built and available in your PATH.
# - A sample Solidity file (test_solidity.sol) exists in the current directory.
# - `jq` is installed for JSON processing.
# ------------------------------------------------------------------------------

set -euo pipefail # Exit on error, undefined variable, or pipeline failure

# Variables
SOLIDITY_FILE="test_struct_to_json.sol"
STRUCT_NAME="MintParams"
OUTPUT_FILE="${STRUCT_NAME}.json"
EXPECTED_JSON='[
  {
    "amount0Desired": "1000000000000000000",
    "amount0Min": "900000000000000000",
    "amount1Desired": "1000000",
    "amount1Min": "900000",
    "deadline": "1695660000",
    "recipient": "0xhidd3n",
    "tickLower": -887220,
    "tickUpper": 887220,
    "token0": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
    "token1": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
  }
]'

# Cleanup function
cleanup() {
    if [[ -f "$OUTPUT_FILE" ]]; then
        rm "$OUTPUT_FILE"
        echo "Cleaned up $OUTPUT_FILE"
    fi
}

# Register cleanup function to run on script exit
trap cleanup EXIT

# Check prerequisites
if ! command -v jq &> /dev/null; then
    echo "Error: 'jq' is not installed. Please install it and try again."
    exit 1
fi

if [[ ! -f "$SOLIDITY_FILE" ]]; then
    echo "Error: Solidity file $SOLIDITY_FILE not found!"
    exit 1
fi

# Test 1: Check if the program generates the correct JSON
echo "Running Test 1: Generate JSON from Solidity file"
go run . generate-json "$SOLIDITY_FILE" "$STRUCT_NAME" <<EOF
0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2
0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
-887220
887220
1000000000000000000
1000000
900000000000000000
900000
0xhidd3n
1695660000
EOF

# Check if the CLI command succeeded
if [[ $? -ne 0 ]]; then
    echo "Error: CLI command failed!"
    exit 1
fi

# Test 2: Check if the output file was created
echo "Running Test 2: Save JSON to file"
if [[ ! -f "$OUTPUT_FILE" ]]; then
    echo "Error: Output file $OUTPUT_FILE not found!"
    exit 1
fi
echo "Output file $OUTPUT_FILE created successfully."

# Test 3: Verify the content of the output file
echo "Running Test 3: Verify JSON content"
ACTUAL_JSON=$(cat "$OUTPUT_FILE")

# Normalize JSON (remove whitespace and newlines)
NORMALIZED_EXPECTED=$(echo "$EXPECTED_JSON" | jq -c .)
NORMALIZED_ACTUAL=$(echo "$ACTUAL_JSON" | jq -c .)

if [[ "$NORMALIZED_ACTUAL" != "$NORMALIZED_EXPECTED" ]]; then
    echo "Error: JSON content does not match expected output!"
    echo "Expected:"
    echo "$NORMALIZED_EXPECTED"
    echo "Got:"
    echo "$NORMALIZED_ACTUAL"
    exit 1
fi
echo "JSON content matches expected output."

echo "All tests passed!"