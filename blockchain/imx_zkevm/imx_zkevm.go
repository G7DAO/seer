package imx_zkevm

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/protobuf/proto"

	seer_common "github.com/moonstream-to/seer/blockchain/common"
	"github.com/moonstream-to/seer/indexer"
	"github.com/moonstream-to/seer/version"
)

func NewClient(url string, timeout int) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Client{rpcClient: rpcClient}, nil
}

// Client is a wrapper around the Ethereum JSON-RPC client.

type Client struct {
	rpcClient *rpc.Client
}

// Client common

// ChainType returns the chain type.
func (c *Client) ChainType() string {
	return "imx_zkevm"
}

// Close closes the underlying RPC client.
func (c *Client) Close() {
	c.rpcClient.Close()
}

// GetLatestBlockNumber returns the latest block number.
func (c *Client) GetLatestBlockNumber() (*big.Int, error) {
	var result string
	if err := c.rpcClient.CallContext(context.Background(), &result, "eth_blockNumber"); err != nil {
		return nil, err
	}

	// Convert the hex string to a big.Int
	blockNumber, ok := new(big.Int).SetString(result, 0) // The 0 base lets the function infer the base from the string prefix.
	if !ok {
		return nil, fmt.Errorf("invalid block number format: %s", result)
	}

	return blockNumber, nil
}

// BlockByNumber returns the block with the given number.
func (c *Client) GetBlockByNumber(ctx context.Context, number *big.Int, withTransactions bool) (*seer_common.BlockJson, error) {

	var rawResponse json.RawMessage // Use RawMessage to capture the entire JSON response
	err := c.rpcClient.CallContext(ctx, &rawResponse, "eth_getBlockByNumber", "0x"+number.Text(16), withTransactions)
	if err != nil {
		fmt.Println("Error calling eth_getBlockByNumber: ", err)
		return nil, err
	}

	var response_json map[string]interface{}

	err = json.Unmarshal(rawResponse, &response_json)

	delete(response_json, "transactions")

	var block *seer_common.BlockJson
	err = c.rpcClient.CallContext(ctx, &block, "eth_getBlockByNumber", "0x"+number.Text(16), true) // true to include transactions
	return block, err
}

// BlockByHash returns the block with the given hash.
func (c *Client) BlockByHash(ctx context.Context, hash common.Hash) (*seer_common.BlockJson, error) {
	var block *seer_common.BlockJson
	err := c.rpcClient.CallContext(ctx, &block, "eth_getBlockByHash", hash, true) // true to include transactions
	return block, err
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
func (c *Client) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	err := c.rpcClient.CallContext(ctx, &receipt, "eth_getTransactionReceipt", hash)
	return receipt, err
}

// Get bytecode of a contract by address.
func (c *Client) GetCode(ctx context.Context, address common.Address, blockNumber uint64) ([]byte, error) {
	var code hexutil.Bytes
	if blockNumber == 0 {
		latestBlockNumber, err := c.GetLatestBlockNumber()
		if err != nil {
			return nil, err
		}
		blockNumber = latestBlockNumber.Uint64()
	}
	err := c.rpcClient.CallContext(ctx, &code, "eth_getCode", address, "0x"+fmt.Sprintf("%x", blockNumber))
	if err != nil {
		log.Printf("Failed to get code for address %s at block %d: %v", address.Hex(), blockNumber, err)
		return nil, err
	}

	if len(code) == 0 {
		return nil, nil
	}
	return code, nil
}
func (c *Client) ClientFilterLogs(ctx context.Context, q ethereum.FilterQuery, debug bool) ([]*seer_common.EventJson, error) {
	var logs []*seer_common.EventJson
	fromBlock := q.FromBlock
	toBlock := q.ToBlock
	batchStep := new(big.Int).Sub(toBlock, fromBlock) // Calculate initial batch step

	for {
		// Calculate the next "lastBlock" within the batch step or adjust to "toBlock" if exceeding
		nextBlock := new(big.Int).Add(fromBlock, batchStep)
		if nextBlock.Cmp(toBlock) > 0 {
			nextBlock = new(big.Int).Set(toBlock)
		}

		var result []*seer_common.EventJson
		err := c.rpcClient.CallContext(ctx, &result, "eth_getLogs", struct {
			FromBlock string           `json:"fromBlock"`
			ToBlock   string           `json:"toBlock"`
			Addresses []common.Address `json:"addresses"`
			Topics    [][]common.Hash  `json:"topics"`
		}{
			FromBlock: toHex(fromBlock),
			ToBlock:   toHex(nextBlock),
			Addresses: q.Addresses,
			Topics:    q.Topics,
		})

		if err != nil {
			if strings.Contains(err.Error(), "query returned more than 10000 results") {
				// Halve the batch step if too many results and retry
				batchStep.Div(batchStep, big.NewInt(2))
				if batchStep.Cmp(big.NewInt(1)) < 0 {
					// If the batch step is too small we will skip that block
					fromBlock = new(big.Int).Add(nextBlock, big.NewInt(1))
					if fromBlock.Cmp(toBlock) > 0 {
						break
					}
					continue
				}
				continue
			} else {
				// For any other error, return immediately
				return nil, err
			}
		}

		// Append the results and adjust "fromBlock" for the next batch
		logs = append(logs, result...)
		fromBlock = new(big.Int).Add(nextBlock, big.NewInt(1))

		if debug {
			log.Printf("Fetched logs: %d", len(result))
		}

		// Break the loop if we've reached or exceeded "toBlock"
		if fromBlock.Cmp(toBlock) > 0 {
			break
		}
	}

	return logs, nil
}

// Utility function to convert big.Int to its hexadecimal representation.
func toHex(number *big.Int) string {
	return fmt.Sprintf("0x%x", number)
}

func fromHex(hex string) *big.Int {
	number := new(big.Int)
	number.SetString(hex, 0)
	return number
}

// FetchBlocksInRange fetches blocks within a specified range.
// This could be useful for batch processing or analysis.
func (c *Client) FetchBlocksInRange(from, to *big.Int, debug bool) ([]*seer_common.BlockJson, error) {
	var blocks []*seer_common.BlockJson
	ctx := context.Background() // For simplicity, using a background context; consider timeouts for production.

	for i := new(big.Int).Set(from); i.Cmp(to) <= 0; i.Add(i, big.NewInt(1)) {
		block, err := c.GetBlockByNumber(ctx, i, true)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, block)
		if debug {
			log.Printf("Fetched block number: %d", i)
		}
	}

	return blocks, nil
}

// FetchBlocksInRangeAsync fetches blocks within a specified range concurrently.
func (c *Client) FetchBlocksInRangeAsync(from, to *big.Int, debug bool, maxRequests int) ([]*seer_common.BlockJson, error) {
	var (
		blocks []*seer_common.BlockJson

		mu  sync.Mutex
		wg  sync.WaitGroup
		ctx = context.Background()
	)

	var blockNumbersRange []*big.Int
	for i := new(big.Int).Set(from); i.Cmp(to) <= 0; i.Add(i, big.NewInt(1)) {
		blockNumbersRange = append(blockNumbersRange, new(big.Int).Set(i))
	}

	sem := make(chan struct{}, maxRequests)             // Semaphore to control concurrency
	errChan := make(chan error, 1)

	for _, b := range blockNumbersRange {
		wg.Add(1)
		go func(b *big.Int) {
			defer wg.Done()

			sem <- struct{}{} // Acquire semaphore

			block, getErr := c.GetBlockByNumber(ctx, b, true)
			if getErr != nil {
				log.Printf("Failed to fetch block number: %d, error: %v", b, getErr)
				errChan <- getErr
				return
			}

			mu.Lock()
			blocks = append(blocks, block)
			mu.Unlock()

			if debug {
				log.Printf("Fetched block number: %d", b)
			}

			<-sem
		}(b)
	}

	wg.Wait()
	close(sem)
	close(errChan)

	if err := <-errChan; err != nil {
		return nil, err
	}

	return blocks, nil
}

// ParseBlocksWithTransactions parses blocks and their transactions into custom data structure.
// This method showcases how to handle and transform detailed block and transaction data.
func (c *Client) ParseBlocksWithTransactions(from, to *big.Int, debug bool, maxRequests int) ([]*ImxZkevmBlock, error) {
	var blocksWithTxsJson []*seer_common.BlockJson
	var fetchErr error
	if maxRequests > 1 {
		blocksWithTxsJson, fetchErr = c.FetchBlocksInRangeAsync(from, to, debug, maxRequests)
	} else {
		blocksWithTxsJson, fetchErr = c.FetchBlocksInRange(from, to, debug)
	}
	if fetchErr != nil {
		return nil, fetchErr
	}

	var parsedBlocks []*ImxZkevmBlock
	for _, blockAndTxsJson := range blocksWithTxsJson {
		// Convert BlockJson to Block and Transactions as required.
		parsedBlock := ToProtoSingleBlock(blockAndTxsJson)

		for _, txJson := range blockAndTxsJson.Transactions {
			txJson.BlockTimestamp = blockAndTxsJson.Timestamp

			parsedTransaction := ToProtoSingleTransaction(&txJson)
			parsedBlock.Transactions = append(parsedBlock.Transactions, parsedTransaction)
		}

		parsedBlocks = append(parsedBlocks, parsedBlock)
	}

	return parsedBlocks, nil
}

func (c *Client) ParseEvents(from, to *big.Int, blocksCache map[uint64]indexer.BlockCache, debug bool) ([]*ImxZkevmEventLog, error) {
	logs, err := c.ClientFilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: from,
		ToBlock:   to,
	}, debug)

	if err != nil {
		fmt.Println("Error fetching logs: ", err)
		return nil, err
	}

	var parsedEvents []*ImxZkevmEventLog

	for _, log := range logs {
		parsedEvent := ToProtoSingleEventLog(log)
		parsedEvents = append(parsedEvents, parsedEvent)

	}

	return parsedEvents, nil
}

func (c *Client) FetchAsProtoBlocksWithEvents(from, to *big.Int, debug bool, maxRequests int) ([]proto.Message, []indexer.BlockIndex, uint64, error) {
	blocks, err := c.ParseBlocksWithTransactions(from, to, debug, maxRequests)
	if err != nil {
		return nil, nil, 0, err
	}

	var blocksSize uint64

	blocksCache := make(map[uint64]indexer.BlockCache)

	for _, block := range blocks {
		blocksCache[block.BlockNumber] = indexer.BlockCache{
			BlockNumber:    block.BlockNumber,
			BlockHash:      block.Hash,
			BlockTimestamp: block.Timestamp,
		} // Assuming block.BlockNumber is int64 and block.Hash is string
	}

	events, err := c.ParseEvents(from, to, blocksCache, debug)
	if err != nil {
		return nil, nil, 0, err
	}

	var blocksProto []proto.Message
	var blocksIndex []indexer.BlockIndex

	for bI, block := range blocks {
		for _, tx := range block.Transactions {
			for _, event := range events {
				if tx.Hash == event.TransactionHash {
					tx.Logs = append(tx.Logs, event)
				}
			}
		}
		

		// Prepare blocks to index
		blocksIndex = append(blocksIndex, indexer.NewBlockIndex("imx_zkevm",
			block.BlockNumber,
			block.Hash,
			block.Timestamp,
			block.ParentHash,
			uint64(bI),
			"",
			0,
		))

		blocksSize += uint64(proto.Size(block))
		blocksProto = append(blocksProto, block) // Assuming block is already a proto.Message
	}

	return blocksProto, blocksIndex, blocksSize, nil
}

func (c *Client) ProcessBlocksToBatch(msgs []proto.Message) (proto.Message, error) {
	var blocks []*ImxZkevmBlock
	for _, msg := range msgs {
		block, ok := msg.(*ImxZkevmBlock)
		if !ok {
			return nil, fmt.Errorf("failed to type assert proto.Message to *ImxZkevmBlock")
		}
		blocks = append(blocks, block)
	}

	return &ImxZkevmBlocksBatch{
		Blocks: blocks,
		SeerVersion: version.SeerVersion,
	}, nil
}

func ToEntireBlocksBatchFromLogProto(obj *ImxZkevmBlocksBatch) *seer_common.BlocksBatchJson {
	blocksBatchJson := seer_common.BlocksBatchJson{
		Blocks: []seer_common.BlockJson{},
		SeerVersion: obj.SeerVersion,
	}

	for _, b := range obj.Blocks {
		var txs []seer_common.TransactionJson
		for _, tx := range b.Transactions {
			var accessList []seer_common.AccessList
			for _, al := range tx.AccessList {
				accessList = append(accessList, seer_common.AccessList{
					Address:     al.Address,
					StorageKeys: al.StorageKeys,
				})
			}
			var events []seer_common.EventJson
			for _, e := range tx.Logs {
				events = append(events, seer_common.EventJson{
					Address:          e.Address,
					Topics:           e.Topics,
					Data:             e.Data,
					BlockNumber:      fmt.Sprintf("%d", e.BlockNumber),
					TransactionHash:  e.TransactionHash,
					BlockHash:        e.BlockHash,
					Removed:          e.Removed,
					LogIndex:         fmt.Sprintf("%d", e.LogIndex),
					TransactionIndex: fmt.Sprintf("%d", e.TransactionIndex),
				})
			}
			txs = append(txs, seer_common.TransactionJson{
				BlockHash:            tx.BlockHash,
				BlockNumber:          fmt.Sprintf("%d", tx.BlockNumber),
				ChainId:              tx.ChainId,
				FromAddress:          tx.FromAddress,
				Gas:                  tx.Gas,
				GasPrice:             tx.GasPrice,
				Hash:                 tx.Hash,
				Input:                tx.Input,
				MaxFeePerGas:         tx.MaxFeePerGas,
				MaxPriorityFeePerGas: tx.MaxPriorityFeePerGas,
				Nonce:                tx.Nonce,
				V:                    tx.V,
				R:                    tx.R,
				S:                    tx.S,
				ToAddress:            tx.ToAddress,
				TransactionIndex:     fmt.Sprintf("%d", tx.TransactionIndex),
				TransactionType:      fmt.Sprintf("%d", tx.TransactionType),
				Value:                tx.Value,
				IndexedAt:            fmt.Sprintf("%d", tx.IndexedAt),
				BlockTimestamp:       fmt.Sprintf("%d", tx.BlockTimestamp),
				AccessList:           accessList,
				YParity:              tx.YParity,

				Events: events,
			})
		}

		blocksBatchJson.Blocks = append(blocksBatchJson.Blocks, seer_common.BlockJson{
			Difficulty:       fmt.Sprintf("%d", b.Difficulty),
			ExtraData:        b.ExtraData,
			GasLimit:         fmt.Sprintf("%d", b.GasLimit),
			GasUsed:          fmt.Sprintf("%d", b.GasUsed),
			Hash:             b.Hash,
			LogsBloom:        b.LogsBloom,
			Miner:            b.Miner,
			Nonce:            b.Nonce,
			BlockNumber:      fmt.Sprintf("%d", b.BlockNumber),
			ParentHash:       b.ParentHash,
			ReceiptsRoot:     b.ReceiptsRoot,
			Sha3Uncles:       b.Sha3Uncles,
			StateRoot:        b.StateRoot,
			Timestamp:        fmt.Sprintf("%d", b.Timestamp),
			TotalDifficulty:  b.TotalDifficulty,
			TransactionsRoot: b.TransactionsRoot,
			Size:             fmt.Sprintf("%d", b.Size),
			BaseFeePerGas:    b.BaseFeePerGas,
			IndexedAt:        fmt.Sprintf("%d", b.IndexedAt),

			
			
			
			

			Transactions: txs,
		})
	}

	return &blocksBatchJson
}

func ToProtoSingleBlock(obj *seer_common.BlockJson) *ImxZkevmBlock {
	return &ImxZkevmBlock{
		BlockNumber:      fromHex(obj.BlockNumber).Uint64(),
		Difficulty:       fromHex(obj.Difficulty).Uint64(),
		ExtraData:        obj.ExtraData,
		GasLimit:         fromHex(obj.GasLimit).Uint64(),
		GasUsed:          fromHex(obj.GasUsed).Uint64(),
		BaseFeePerGas:    obj.BaseFeePerGas,
		Hash:             obj.Hash,
		LogsBloom:        obj.LogsBloom,
		Miner:            obj.Miner,
		Nonce:            obj.Nonce,
		ParentHash:       obj.ParentHash,
		ReceiptsRoot:     obj.ReceiptsRoot,
		Sha3Uncles:       obj.Sha3Uncles,
		Size:             fromHex(obj.Size).Uint64(),
		StateRoot:        obj.StateRoot,
		Timestamp:        fromHex(obj.Timestamp).Uint64(),
		TotalDifficulty:  obj.TotalDifficulty,
		TransactionsRoot: obj.TransactionsRoot,
		IndexedAt:        fromHex(obj.IndexedAt).Uint64(),

		
		
		
		
	}
}

func ToProtoSingleTransaction(obj *seer_common.TransactionJson) *ImxZkevmTransaction {
	var accessList []*ImxZkevmTransactionAccessList
	for _, al := range obj.AccessList {
		accessList = append(accessList, &ImxZkevmTransactionAccessList{
			Address:     al.Address,
			StorageKeys: al.StorageKeys,
		})
	}

	return &ImxZkevmTransaction{
		Hash:                 obj.Hash,
		BlockNumber:          fromHex(obj.BlockNumber).Uint64(),
		BlockHash:            obj.BlockHash,
		FromAddress:          obj.FromAddress,
		ToAddress:            obj.ToAddress,
		Gas:                  obj.Gas,
		GasPrice:             obj.GasPrice,
		MaxFeePerGas:         obj.MaxFeePerGas,
		MaxPriorityFeePerGas: obj.MaxPriorityFeePerGas,
		Input:                obj.Input,
		Nonce:                obj.Nonce,
		TransactionIndex:     fromHex(obj.TransactionIndex).Uint64(),
		TransactionType:      fromHex(obj.TransactionType).Uint64(),
		Value:                obj.Value,
		IndexedAt:            fromHex(obj.IndexedAt).Uint64(),
		BlockTimestamp:       fromHex(obj.BlockTimestamp).Uint64(),

		ChainId: obj.ChainId,
		V:       obj.V,
		R:       obj.R,
		S:       obj.S,

		AccessList: accessList,
		YParity:    obj.YParity,
	}
}

func ToEvenFromLogProto(obj *ImxZkevmEventLog) *seer_common.EventJson {
	return &seer_common.EventJson{
		Address:         obj.Address,
		Topics:          obj.Topics,
		Data:            obj.Data,
		BlockNumber:     fmt.Sprintf("%d", obj.BlockNumber),
		TransactionHash: obj.TransactionHash,
		LogIndex:        fmt.Sprintf("%d", obj.LogIndex),
		BlockHash:       obj.BlockHash,
		Removed:         obj.Removed,
	}
}

func ToProtoSingleEventLog(obj *seer_common.EventJson) *ImxZkevmEventLog {
	return &ImxZkevmEventLog{
		Address:         obj.Address,
		Topics:          obj.Topics,
		Data:            obj.Data,
		BlockNumber:     fromHex(obj.BlockNumber).Uint64(),
		TransactionHash: obj.TransactionHash,
		LogIndex:        fromHex(obj.LogIndex).Uint64(),
		BlockHash:       obj.BlockHash,
		Removed:         obj.Removed,
	}
}

func (c *Client) DecodeProtoEventLogs(data []string) ([]*ImxZkevmEventLog, error) {
	var events []*ImxZkevmEventLog
	for _, d := range data {
		var event ImxZkevmEventLog
		base64Decoded, err := base64.StdEncoding.DecodeString(d)
		if err != nil {
			return nil, err
		}
		if err := proto.Unmarshal(base64Decoded, &event); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (c *Client) DecodeProtoTransactions(data []string) ([]*ImxZkevmTransaction, error) {
	var transactions []*ImxZkevmTransaction
	for _, d := range data {
		var transaction ImxZkevmTransaction
		base64Decoded, err := base64.StdEncoding.DecodeString(d)
		if err != nil {
			return nil, err
		}
		if err := proto.Unmarshal(base64Decoded, &transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (c *Client) DecodeProtoBlocks(data []string) ([]*ImxZkevmBlock, error) {
	var blocks []*ImxZkevmBlock
	for _, d := range data {
		var block ImxZkevmBlock
		base64Decoded, err := base64.StdEncoding.DecodeString(d)
		if err != nil {
			return nil, err
		}
		if err := proto.Unmarshal(base64Decoded, &block); err != nil {
			return nil, err
		}
		blocks = append(blocks, &block)
	}
	return blocks, nil
}

func (c *Client) DecodeProtoEntireBlockToJson(rawData *bytes.Buffer) (*seer_common.BlocksBatchJson, error) {
	var protoBlocksBatch ImxZkevmBlocksBatch

    dataBytes := rawData.Bytes()

    err := proto.Unmarshal(dataBytes, &protoBlocksBatch)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal data: %v", err)
    }

	blocksBatchJson := ToEntireBlocksBatchFromLogProto(&protoBlocksBatch)

	return blocksBatchJson, nil
}

func (c *Client) DecodeProtoEntireBlockToLabels(rawData *bytes.Buffer, abiMap map[string]map[string]map[string]string) ([]indexer.EventLabel, []indexer.TransactionLabel, error) {
	var protoBlocksBatch ImxZkevmBlocksBatch

	dataBytes := rawData.Bytes()

    err := proto.Unmarshal(dataBytes, &protoBlocksBatch)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to unmarshal data: %v", err)
    }

	var labels []indexer.EventLabel
	var txLabels []indexer.TransactionLabel
	var decodeErr error

	for _, b := range protoBlocksBatch.Blocks {
		for _, tx := range b.Transactions {
			var decodedArgsTx map[string]interface{}

			label := indexer.SeerCrawlerLabel

			if len(tx.Input) < 10 { // If input is less than 3 characters then it direct transfer
				continue
			}

			// Process transaction labels
			selector := tx.Input[:10]

			if abiMap[tx.ToAddress] != nil && abiMap[tx.ToAddress][selector] != nil {
				txContractAbi, err := abi.JSON(strings.NewReader(abiMap[tx.ToAddress][selector]["abi"]))
				if err != nil {
					fmt.Println("Error initializing contract ABI transactions: ", err)
					return nil, nil, err
				}

				inputData, err := hex.DecodeString(tx.Input[2:])
				if err != nil {
					fmt.Println("Error decoding input data: ", err)
					return nil, nil, err
				}
				
				decodedArgsTx, decodeErr = seer_common.DecodeTransactionInputDataToInterface(&txContractAbi, inputData)
				if decodeErr != nil {
					fmt.Println("Error decoding transaction not decoded data: ", tx.Hash, decodeErr)
					decodedArgsTx = map[string]interface{}{
						"input_raw": tx,
						"abi": abiMap[tx.ToAddress][selector]["abi"],
						"selector": selector,
						"error": decodeErr,
					}
					label = indexer.SeerCrawlerRawLabel
				}

				receipt, err := c.TransactionReceipt(context.Background(), common.HexToHash(tx.Hash))

				if err != nil {
					fmt.Println("Error fetching transaction receipt: ", err)
					return nil, nil, err
				}

				// check if the transaction was successful
				if receipt.Status == 1 {
					decodedArgsTx["status"] = 1
				} else {
					decodedArgsTx["status"] = 0
				}

				txLabelDataBytes, err := json.Marshal(decodedArgsTx)
				if err != nil {
					fmt.Println("Error converting decodedArgsTx to JSON: ", err)
					return nil, nil, err
				}

				// Convert transaction to label
				transactionLabel := indexer.TransactionLabel{
					Address:         tx.ToAddress,
					BlockNumber:     tx.BlockNumber,
					BlockHash:       tx.BlockHash,
					CallerAddress:   tx.FromAddress,
					LabelName:       abiMap[tx.ToAddress][selector]["abi_name"],
					LabelType:       "tx_call",
					OriginAddress:   tx.FromAddress,
					Label:           label,
					TransactionHash: tx.Hash,
					LabelData:       string(txLabelDataBytes), // Convert JSON byte slice to string
					BlockTimestamp:  b.Timestamp,
				}

				txLabels = append(txLabels, transactionLabel)
			}

			// Process events
			for _, e := range tx.Logs {
				var decodedArgsLogs map[string]interface{}
				label = indexer.SeerCrawlerLabel

				var topicSelector string

				if len(e.Topics) > 0 {
					topicSelector = e.Topics[0]
				} else {
					// 0x0 is the default topic selector
					topicSelector = "0x0"
				}

				if abiMap[e.Address] == nil || abiMap[e.Address][topicSelector] == nil {
					continue
				}

				// Get the ABI string
				contractAbi, err := abi.JSON(strings.NewReader(abiMap[e.Address][topicSelector]["abi"]))
				if err != nil {
					fmt.Println("Error initializing contract ABI: ", err)
					return nil, nil, err
				}


				// Decode the event data
				decodedArgsLogs, decodeErr = seer_common.DecodeLogArgsToLabelData(&contractAbi, e.Topics, e.Data)
				if decodeErr != nil {
					fmt.Println("Error decoding event not decoded data: ", e.TransactionHash, decodeErr)
					decodedArgsLogs = map[string]interface{}{
						"input_raw": e,
						"abi": abiMap[e.Address][topicSelector]["abi"],
						"selector": topicSelector,
						"error": decodeErr,
					}
					label = indexer.SeerCrawlerRawLabel
				}

				// Convert decodedArgsLogs map to JSON
				labelDataBytes, err := json.Marshal(decodedArgsLogs)
				if err != nil {
					fmt.Println("Error converting decodedArgsLogs to JSON: ", err)
					return nil, nil, err
				}

				// Convert event to label
				eventLabel := indexer.EventLabel{
					Label:           label,
					LabelName:       abiMap[e.Address][topicSelector]["abi_name"],
					LabelType:       "event",
					BlockNumber:     e.BlockNumber,
					BlockHash:       e.BlockHash,
					Address:         e.Address,
					OriginAddress:   tx.FromAddress,
					TransactionHash: e.TransactionHash,
					LabelData:       string(labelDataBytes), // Convert JSON byte slice to string
					BlockTimestamp:  b.Timestamp,
					LogIndex:        e.LogIndex,
				}

				labels = append(labels, eventLabel)
			}
		}
	}

	return labels, txLabels, nil
}

func (c *Client) DecodeProtoTransactionsToLabels(transactions []string, blocksCache map[uint64]uint64, abiMap map[string]map[string]map[string]string) ([]indexer.TransactionLabel, error) {

	decodedTransactions, err := c.DecodeProtoTransactions(transactions)

	if err != nil {
		return nil, err
	}

	var labels []indexer.TransactionLabel
	var decodedArgs map[string]interface{}
	var decodeErr error


	for _, transaction := range decodedTransactions {

		label := indexer.SeerCrawlerLabel

		selector := transaction.Input[:10]

		contractAbi, err := abi.JSON(strings.NewReader(abiMap[transaction.ToAddress][selector]["abi"]))

		if err != nil {
			return nil, err
		}

		inputData, err := hex.DecodeString(transaction.Input[2:])
		if err != nil {
			fmt.Println("Error decoding input data: ", err)
			return nil, err
		}

		decodedArgs, decodeErr = seer_common.DecodeTransactionInputDataToInterface(&contractAbi, inputData)

		if decodeErr != nil {
			fmt.Println("Error decoding transaction not decoded data: ", transaction.Hash, decodeErr)
			decodedArgs = map[string]interface{}{
				"input_raw": transaction,
				"abi": abiMap[transaction.ToAddress][selector]["abi"],
				"selector": selector,
				"error": decodeErr,
			}
			label = indexer.SeerCrawlerRawLabel
		}

		labelDataBytes, err := json.Marshal(decodedArgs)
		if err != nil {
			fmt.Println("Error converting decodedArgs to JSON: ", err)
			return nil, err
		}

		// Convert JSON byte slice to string
		labelDataString := string(labelDataBytes)

		// Convert transaction to label
		transactionLabel := indexer.TransactionLabel{
			Address:         transaction.ToAddress,
			BlockNumber:     transaction.BlockNumber,
			BlockHash:       transaction.BlockHash,
			CallerAddress:   transaction.FromAddress,
			LabelName:       abiMap[transaction.ToAddress][selector]["abi_name"],
			LabelType:       "tx_call",
			OriginAddress:   transaction.FromAddress,
			Label:           label,
			TransactionHash: transaction.Hash,
			LabelData:       labelDataString,
			BlockTimestamp:  blocksCache[transaction.BlockNumber],
		}

		labels = append(labels, transactionLabel)

	}

	return labels, nil
}

func (c *Client) GetTransactionByHash(ctx context.Context, hash string) (*seer_common.TransactionJson, error) {
	var tx *seer_common.TransactionJson
	err := c.rpcClient.CallContext(ctx, &tx, "eth_getTransactionByHash", hash)
	return tx, err
}

func (c *Client) GetTransactionsLabels(startBlock uint64, endBlock uint64, abiMap map[string]map[string]map[string]string) ([]indexer.TransactionLabel, map[uint64]seer_common.BlockWithTransactions, error) {
	var transactionsLabels []indexer.TransactionLabel

	var blocksCache map[uint64]seer_common.BlockWithTransactions

	// Get blocks in range
	blocks, err := c.FetchBlocksInRangeAsync(big.NewInt(int64(startBlock)), big.NewInt(int64(endBlock)), false, 1)

	if err != nil {
		return nil, nil, err
	}

	// Get transactions in range

	for _, block := range blocks {

		blockNumber, err := strconv.ParseUint(block.BlockNumber, 0, 64)
		if err != nil {
			log.Fatalf("Failed to convert BlockNumber to uint64: %v", err)
		}

		blockTimestamp, err := strconv.ParseUint(block.Timestamp, 0, 64)

		if err != nil {
			log.Fatalf("Failed to convert BlockTimestamp to uint64: %v", err)
		}

		blocksCache[blockNumber] = seer_common.BlockWithTransactions{
			BlockNumber:    blockNumber,
			BlockHash:      block.Hash,
			BlockTimestamp: blockTimestamp,
			Transactions:   make(map[string]seer_common.TransactionJson),
		}

		for _, tx := range block.Transactions {

			label := indexer.SeerCrawlerLabel

			if len(tx.Input) < 10 { // If input is less than 3 characters then it direct transfer
				continue
			}
			// Fill blocks cache
			blocksCache[blockNumber].Transactions[tx.Hash] = tx

			// Process transaction labels

			selector := tx.Input[:10]

			if abiMap[tx.ToAddress] != nil && abiMap[tx.ToAddress][selector] != nil {
				txContractAbi, err := abi.JSON(strings.NewReader(abiMap[tx.ToAddress][selector]["abi"]))
				if err != nil {
					fmt.Println("Error initializing contract ABI transactions: ", err)
					return nil, nil, err
				}

				inputData, err := hex.DecodeString(tx.Input[2:])
				if err != nil {
					fmt.Println("Error decoding input data: ", err)
					return nil, nil, err
				}

				decodedArgsTx, decodeErr := seer_common.DecodeTransactionInputDataToInterface(&txContractAbi, inputData)
				if decodeErr != nil {
					fmt.Println("Error decoding transaction not decoded data: ", tx.Hash, decodeErr)
					decodedArgsTx = map[string]interface{}{
						"input_raw": tx,
						"abi":       abiMap[tx.ToAddress][selector]["abi"],
						"selector":  selector,
						"error":     decodeErr,
					}
					label = indexer.SeerCrawlerRawLabel
				}

				receipt, err := c.TransactionReceipt(context.Background(), common.HexToHash(tx.Hash))

				if err != nil {
					fmt.Println("Error fetching transaction receipt: ", err)
					return nil, nil, err
				}

				// check if the transaction was successful
				if receipt.Status == 1 {
					decodedArgsTx["status"] = 1
				} else {
					decodedArgsTx["status"] = 0
				}

				txLabelDataBytes, err := json.Marshal(decodedArgsTx)
				if err != nil {
					fmt.Println("Error converting decodedArgsTx to JSON: ", err)
					return nil, nil, err
				}

				// Convert transaction to label
				transactionLabel := indexer.TransactionLabel{
					Address:         tx.ToAddress,
					BlockNumber:     blockNumber,
					BlockHash:       tx.BlockHash,
					CallerAddress:   tx.FromAddress,
					LabelName:       abiMap[tx.ToAddress][selector]["abi_name"],
					LabelType:       "tx_call",
					OriginAddress:   tx.FromAddress,
					Label:           label,
					TransactionHash: tx.Hash,
					LabelData:       string(txLabelDataBytes), // Convert JSON byte slice to string
					BlockTimestamp:  blockTimestamp,
				}

				transactionsLabels = append(transactionsLabels, transactionLabel)
			}

		}

	}

	return transactionsLabels, blocksCache, nil

}

func (c *Client) GetEventsLabels(startBlock uint64, endBlock uint64, abiMap map[string]map[string]map[string]string, blocksCache map[uint64]seer_common.BlockWithTransactions) ([]indexer.EventLabel, error) {
	var eventsLabels []indexer.EventLabel

	if blocksCache == nil {
		blocksCache = make(map[uint64]seer_common.BlockWithTransactions)
	}

	// Get events in range

	var addresses []common.Address
	var topics []common.Hash

	for address, selectorMap := range abiMap {
		for selector, _ := range selectorMap {
			topics = append(topics, common.HexToHash(selector))
		}

		addresses = append(addresses, common.HexToAddress(address))
	}

	// query filter from abiMap
	filter := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startBlock)),
		ToBlock:   big.NewInt(int64(endBlock)),
		Addresses: addresses,
		Topics:    [][]common.Hash{topics},
	}
	logs, err := c.ClientFilterLogs(context.Background(), filter, false)

	if err != nil {
		return nil, err
	}

	for _, log := range logs {
		var decodedArgsLogs map[string]interface{}
		label := indexer.SeerCrawlerLabel

		var topicSelector string

		if len(log.Topics) > 0 {
			topicSelector = log.Topics[0]
		} else {
			// 0x0 is the default topic selector
			topicSelector = "0x0"
		}

		if abiMap[log.Address] == nil || abiMap[log.Address][topicSelector] == nil {
			continue
		}

		// Get the ABI string
		contractAbi, err := abi.JSON(strings.NewReader(abiMap[log.Address][topicSelector]["abi"]))
		if err != nil {
			fmt.Println("Error initializing contract ABI: ", err)
			return nil, err
		}

		// Decode the event data
		decodedArgsLogs, decodeErr := seer_common.DecodeLogArgsToLabelData(&contractAbi, log.Topics, log.Data)
		if decodeErr != nil {
			fmt.Println("Error decoding event not decoded data: ", log.TransactionHash, decodeErr)
			decodedArgsLogs = map[string]interface{}{
				"input_raw": log,
				"abi":       abiMap[log.Address][topicSelector]["abi"],
				"selector":  topicSelector,
				"error":     decodeErr,
			}
			label = indexer.SeerCrawlerRawLabel
		}

		// Convert decodedArgsLogs map to JSON
		labelDataBytes, err := json.Marshal(decodedArgsLogs)
		if err != nil {
			fmt.Println("Error converting decodedArgsLogs to JSON: ", err)
			return nil, err
		}

		blockNumber, err := strconv.ParseUint(log.BlockNumber, 0, 64)
		if err != nil {
			return nil, err
		}

		if _, ok := blocksCache[blockNumber]; !ok {

			// get block from rpc
			block, err := c.GetBlockByNumber(context.Background(), big.NewInt(int64(blockNumber)), true)
			if err != nil {
				return nil, err
			}

			blockTimestamp, err := strconv.ParseUint(block.Timestamp, 0, 64)
			if err != nil {
				return nil, err
			}

			blocksCache[blockNumber] = seer_common.BlockWithTransactions{
				BlockNumber:    blockNumber,
				BlockHash:      block.Hash,
				BlockTimestamp: blockTimestamp,
				Transactions:   make(map[string]seer_common.TransactionJson),
			}

			for _, tx := range block.Transactions {
				blocksCache[blockNumber].Transactions[tx.Hash] = tx
			}

		}

		transaction := blocksCache[blockNumber].Transactions[log.TransactionHash]

		logIndex, err := strconv.ParseUint(log.LogIndex, 0, 64)
		if err != nil {
			return nil, err
		}

		// Convert event to label
		eventLabel := indexer.EventLabel{
			Label:           label,
			LabelName:       abiMap[log.Address][topicSelector]["abi_name"],
			LabelType:       "event",
			BlockNumber:     blockNumber,
			BlockHash:       log.BlockHash,
			Address:         log.Address,
			OriginAddress:   transaction.FromAddress,
			TransactionHash: log.TransactionHash,
			LabelData:       string(labelDataBytes), // Convert JSON byte slice to string
			BlockTimestamp:  blocksCache[blockNumber].BlockTimestamp,
			LogIndex:        logIndex,
		}

		eventsLabels = append(eventsLabels, eventLabel)

	}

	return eventsLabels, nil

}
