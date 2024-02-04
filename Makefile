.PHONY: build clean

build: seer

seer:
	go build .

clean:
	rm -f ownable-erc-721
	rm -f seer

rebuild: clean build

examples/ownable-erc-721/OwnableERC721.go: rebuild
	./seer evm generate --abi fixtures/OwnableERC721.json --bytecode fixtures/OwnableERC721.bin --cli --package main --struct OwnableERC721 --output examples/ownable-erc-721/OwnableERC721.go

ownable-erc-721: examples/ownable-erc-721/OwnableERC721.go
	go build ./examples/ownable-erc-721
