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
    if [ "$BLOCKCHAIN" != "ethereum" ] && [ "$BLOCKCHAIN" != "polygon" ] && [ "$BLOCKCHAIN" != "mantle" ] && [ "$BLOCKCHAIN" != "mantle_sepolia" ] && [ "$BLOCKCHAIN" != "sepolia" ]; then
      ./seer blockchain generate -n $BLOCKCHAIN --side-chain
      echo "Generated interface for side-chain blockchain $BLOCKCHAIN"
    else
      ./seer blockchain generate -n $BLOCKCHAIN
      echo "Generated interface for blockchain $BLOCKCHAIN"
    fi
  fi
done
