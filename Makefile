.PHONY: build clean

build: seer

seer:
	go build .

clean:
	rm -f ownable-erc-721 diamond-cut-facet
	rm -f seer

rebuild: clean build

examples/ownable-erc-721/OwnableERC721.go: rebuild
	./seer evm generate --abi fixtures/OwnableERC721.json --bytecode fixtures/OwnableERC721.bin --cli --package main --struct OwnableERC721 --output examples/ownable-erc-721/OwnableERC721.go --includemain

ownable-erc-721: examples/ownable-erc-721/OwnableERC721.go
	go build ./examples/ownable-erc-721

examples/diamond-cut-facet/DiamondCutFacet.go: rebuild
	./seer evm generate --abi fixtures/DiamondCutFacetABI.json --bytecode fixtures/DiamondCutFacet.bin --cli --package main --struct DiamondCutFacet --output examples/diamond-cut-facet/DiamondCutFacet.go --includemain
