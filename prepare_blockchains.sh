#!/usr/bin/env sh

# Delete all .pb.go and interfaces for each blockchain, re-generate it from .proto
set -e

PROTO_FILES=$(find blockchain/ -name '*.proto')
for PROTO in $PROTO_FILES; do
  protoc --go_out=. --go_opt=paths=source_relative $PROTO
  echo "Regenerated from proto: $PROTO"
done

BLOCKCHAIN_NAMES_RAW=$(find blockchain/ -maxdepth 1 -type d | cut -f2 -d '/')
for BLOCKCHAIN in $BLOCKCHAIN_NAMES_RAW; do
  if [ "$BLOCKCHAIN" != "" ] && [ "$BLOCKCHAIN" != "common" ]; then
    case "$BLOCKCHAIN" in
      "starknet")
        ./seer blockchain generate -n $BLOCKCHAIN -t STARK
        echo "Generated interface for STARK blockchain $BLOCKCHAIN"
        ;;
      "ethereum" | "polygon" | "mantle" | "mantle_sepolia" | "sepolia" | "imx_zkevm" | "imx_zkevm_sepolia")
        ./seer blockchain generate -n $BLOCKCHAIN -t EVM
        echo "Generated interface for EVM blockchain $BLOCKCHAIN"
        ;;
      *)
        ./seer blockchain generate -n $BLOCKCHAIN -t EVM --side-chain
        echo "Generated interface for EVM side-chain blockchain $BLOCKCHAIN"
        ;;
    esac
  fi
done
