package trace

import (
	"fmt"

	"github.com/G7DAO/seer/blockchain/common"
)

// TraceTransaction performs execution trace of a transaction
func TraceTransaction(block *common.BlockJson, parentBlock *common.BlockJson, tx *common.TransactionJson) (*ExecutionResult, error) {
	fmt.Printf("Block: %v\n", block)
	fmt.Printf("Parent Block: %v\n", parentBlock)
	fmt.Printf("Transaction: %v\n", tx)

	return nil, nil
}

// ExecutionResult contains the structured execution trace results
type ExecutionResult struct {
	Gas         uint64
	Failed      bool
	ReturnValue []byte
	StructLogs  []StructLog
}

// StructLog contains a structured log entry from the EVM execution
type StructLog struct {
	Pc      uint64            `json:"pc"`
	Op      string            `json:"op"`
	Gas     uint64            `json:"gas"`
	GasCost uint64            `json:"gasCost"`
	Depth   int               `json:"depth"`
	Stack   []string          `json:"stack"`
	Memory  []string          `json:"memory"`
	Storage map[string]string `json:"storage"`
}
