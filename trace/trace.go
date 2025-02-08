package trace

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/triedb"

	"github.com/G7DAO/seer/blockchain/common"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

// Define OpcodeTracer as a new type whose underlying type is tracing.Hooks.
type OpcodeTracer tracing.Hooks

// NewOpcodeTracer returns a pointer to a new OpcodeTracer.
func NewOpcodeTracer() *OpcodeTracer {
	return &OpcodeTracer{}
}

// getStackValues uses reflection to extract the underlying []*big.Int from the stack.
func getStackValues(stack *vm.Stack) []string {
	v := reflect.ValueOf(stack).Elem()
	dataField := v.FieldByName("data")
	if !dataField.IsValid() {
		return nil
	}
	values := make([]string, dataField.Len())
	for i := 0; i < dataField.Len(); i++ {
		if val, ok := dataField.Index(i).Interface().(*big.Int); ok {
			values[i] = val.Text(16)
		} else {
			values[i] = "?"
		}
	}
	return values
}

// CaptureStart is called at the beginning of EVM execution.
func (ot *OpcodeTracer) CaptureStart(env *vm.EVM) {
	// No action necessary at start.
}

// CaptureState is called for each executed opcode.
func (ot *OpcodeTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, evm *vm.EVM, stack *vm.Stack, memory *vm.Memory, depth int, err error) {
	// Print the PC, opcode, and current stack values.
	stackValues := getStackValues(stack)
	fmt.Printf("PC: %d, Op: %s, Stack: %v\n", pc, op.String(), stackValues)
}

// CaptureFault is called if an opcode causes a fault.
func (ot *OpcodeTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, evm *vm.EVM, stack *vm.Stack, memory *vm.Memory, depth int, err error) {
	fmt.Printf("Fault at PC: %d, Op: %s, error: %v\n", pc, op.String(), err)
}

// CaptureEnd is called at the end of execution.
func (ot *OpcodeTracer) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	fmt.Printf("Execution ended, Gas Used: %d, Duration: %v, Error: %v\n", gasUsed, t, err)
}

// TraceTransaction performs execution trace of a transaction.
// It accepts a block, its parent block, and the transaction to trace.
//
// TODO: Set up state retrieval from parentBlock's state root,
// create an appropriate EVM context, and then execute the transaction.
func TraceTransaction(block *common.BlockJson, parentBlock *common.BlockJson, tx *common.TransactionJson) (*ExecutionResult, error) {
	// Convert string values to appropriate types.
	gasLimit, err := strconv.ParseUint(tx.Gas, 0, 64)
	if err != nil {
		return nil, err
	}
	value := new(big.Int)
	value.SetString(tx.Value, 0)

	// Create EVM configuration using the custom OpcodeTracer.
	cfg := vm.Config{
		// Convert our *OpcodeTracer to *tracing.Hooks since underlying types are identical.
		Tracer:                  (*tracing.Hooks)(NewOpcodeTracer()),
		NoBaseFee:               false,
		EnablePreimageRecording: true,
		ExtraEips:               nil,
	}

	// Build chain configuration.
	chainConfig := params.MainnetChainConfig // Adjust this as appropriate for the active chain.

	var toAddress gethCommon.Address
	toAddress = gethCommon.HexToAddress(tx.ToAddress)

	// Create message from transaction data.
	msg := core.Message{
		From:             gethCommon.HexToAddress(tx.FromAddress),
		To:               &toAddress,
		Value:            value,
		GasLimit:         gasLimit,
		GasPrice:         new(big.Int).SetBytes(hexutil.MustDecode(tx.GasPrice)),
		Data:             hexutil.MustDecode(tx.Input),
		AccessList:       nil, // TODO: Convert tx.AccessList if applicable.
		SkipNonceChecks:  false,
		SkipFromEOACheck: false,
	}

	// Retrieve the stateDB using parent's state root.
	db := rawdb.NewMemoryDatabase()
	// Wrap the mem DB into a triedb.Database using the defaults (this mirrors what is done in stateless.go):
	trieDb := triedb.NewDatabase(db, triedb.HashDefaults)
	// Create the stateDB from the trie database.
	stateDB, err := state.New(gethCommon.HexToHash(parentBlock.StateRoot), state.NewDatabase(trieDb, nil))
	if err != nil {
		return nil, err
	}

	// Ensure that the "to" account exists in state.
	// If this is a contract call, missing account will cause missing trie node errors.
	// Here we create the account with dummy code (empty byte slice) if it doesn't already exist.
	toAccount := gethCommon.HexToAddress(tx.ToAddress)
	if !stateDB.Exist(toAccount) {
		stateDB.CreateAccount(toAccount)
		// Optionally, if you know the contract's code, set it here.
		stateDB.SetCode(toAccount, []byte{})
	}

	// Convert block number from string to *big.Int.
	blockNumber := new(big.Int)
	if _, ok := blockNumber.SetString(block.BlockNumber, 10); !ok {
		return nil, errors.New("invalid block number")
	}

	// Since BlockJson has no Time field, default blockTime to 0.
	var blockTime uint64 = 0

	// Parse block difficulty (assuming block.Difficulty is a string representation of a big int).
	difficulty := new(big.Int)
	if _, ok := difficulty.SetString(block.Difficulty, 10); !ok {
		difficulty = big.NewInt(0)
	}

	// Since BlockJson has no Coinbase field, use a zero address.
	coinbase := gethCommon.Address{}

	// Create the EVM contexts.
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash: func(n uint64) gethCommon.Hash {
			// For simplicity, return zero hash.
			// In a full implementation, this should lookup the corresponding ancestor block's hash.
			return gethCommon.Hash{}
		},
		Coinbase:    coinbase,
		BlockNumber: blockNumber,
		Time:        blockTime,
		Difficulty:  difficulty,
		GasLimit:    gasLimit,
	}
	txContext := vm.TxContext{
		Origin:   gethCommon.HexToAddress(tx.FromAddress),
		GasPrice: new(big.Int).SetBytes(hexutil.MustDecode(tx.GasPrice)),
	}

	// Instantiate the EVM with the new contexts.
	evm := vm.NewEVM(blockContext, txContext, stateDB, chainConfig, cfg)

	// Create a new GasPool and add the available gas.
	gasPool := new(core.GasPool)
	gasPool.AddGas(gasLimit)

	// Execute the transaction using the EVM.
	result, err := core.ApplyMessage(evm, &msg, gasPool)
	if err != nil {
		return nil, err
	}

	// Map the result to our ExecutionResult.
	// Call the functions result.Failed() and result.Return() as they are function fields.
	return &ExecutionResult{
		Gas:         gasLimit - result.UsedGas,
		Failed:      result.Failed(),
		ReturnValue: result.Return(),
		StructLogs:  nil, // If available, assign structured logs here.
	}, nil
}

// ExecutionResult contains the structured execution trace results.
type ExecutionResult struct {
	Gas         uint64
	Failed      bool
	ReturnValue []byte
	StructLogs  []StructLog
}

// StructLog contains a structured log entry from the EVM execution.
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
