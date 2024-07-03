// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.6.1
// source: blockchain/ethereum/ethereum_index_types.proto

package ethereum

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EthereumTransactionAccessList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address     string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	StorageKeys []string `protobuf:"bytes,2,rep,name=storage_keys,json=storageKeys,proto3" json:"storage_keys,omitempty"`
}

func (x *EthereumTransactionAccessList) Reset() {
	*x = EthereumTransactionAccessList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EthereumTransactionAccessList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EthereumTransactionAccessList) ProtoMessage() {}

func (x *EthereumTransactionAccessList) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EthereumTransactionAccessList.ProtoReflect.Descriptor instead.
func (*EthereumTransactionAccessList) Descriptor() ([]byte, []int) {
	return file_blockchain_ethereum_ethereum_index_types_proto_rawDescGZIP(), []int{0}
}

func (x *EthereumTransactionAccessList) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *EthereumTransactionAccessList) GetStorageKeys() []string {
	if x != nil {
		return x.StorageKeys
	}
	return nil
}

// Represents a single transaction within a block
type EthereumTransaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash                 string                           `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	BlockNumber          uint64                           `protobuf:"varint,2,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	FromAddress          string                           `protobuf:"bytes,3,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
	ToAddress            string                           `protobuf:"bytes,4,opt,name=to_address,json=toAddress,proto3" json:"to_address,omitempty"`
	Gas                  string                           `protobuf:"bytes,5,opt,name=gas,proto3" json:"gas,omitempty"` // using string to handle big numeric values
	GasPrice             string                           `protobuf:"bytes,6,opt,name=gas_price,json=gasPrice,proto3" json:"gas_price,omitempty"`
	MaxFeePerGas         string                           `protobuf:"bytes,7,opt,name=max_fee_per_gas,json=maxFeePerGas,proto3" json:"max_fee_per_gas,omitempty"`
	MaxPriorityFeePerGas string                           `protobuf:"bytes,8,opt,name=max_priority_fee_per_gas,json=maxPriorityFeePerGas,proto3" json:"max_priority_fee_per_gas,omitempty"`
	Input                string                           `protobuf:"bytes,9,opt,name=input,proto3" json:"input,omitempty"` // could be a long text
	Nonce                string                           `protobuf:"bytes,10,opt,name=nonce,proto3" json:"nonce,omitempty"`
	TransactionIndex     uint64                           `protobuf:"varint,11,opt,name=transaction_index,json=transactionIndex,proto3" json:"transaction_index,omitempty"`
	TransactionType      uint64                           `protobuf:"varint,12,opt,name=transaction_type,json=transactionType,proto3" json:"transaction_type,omitempty"`
	Value                string                           `protobuf:"bytes,13,opt,name=value,proto3" json:"value,omitempty"`                                          // using string to handle big numeric values
	IndexedAt            uint64                           `protobuf:"varint,14,opt,name=indexed_at,json=indexedAt,proto3" json:"indexed_at,omitempty"`                // using uint64 to represent timestamp
	BlockTimestamp       uint64                           `protobuf:"varint,15,opt,name=block_timestamp,json=blockTimestamp,proto3" json:"block_timestamp,omitempty"` // using uint64 to represent timestam
	BlockHash            string                           `protobuf:"bytes,16,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`                 // Added field for block hash
	ChainId              string                           `protobuf:"bytes,17,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`                       // Used as a field to match potential EIP-1559 transaction types
	V                    string                           `protobuf:"bytes,18,opt,name=v,proto3" json:"v,omitempty"`                                                  // Used as a field to match potential EIP-1559 transaction types
	R                    string                           `protobuf:"bytes,19,opt,name=r,proto3" json:"r,omitempty"`                                                  // Used as a field to match potential EIP-1559 transaction types
	S                    string                           `protobuf:"bytes,20,opt,name=s,proto3" json:"s,omitempty"`                                                  // Used as a field to match potential EIP-1559 transaction types
	AccessList           []*EthereumTransactionAccessList `protobuf:"bytes,21,rep,name=access_list,json=accessList,proto3" json:"access_list,omitempty"`
	YParity              string                           `protobuf:"bytes,22,opt,name=y_parity,json=yParity,proto3" json:"y_parity,omitempty"` // Used as a field to match potential EIP-1559 transaction types
	Logs                 []*EthereumEventLog              `protobuf:"bytes,23,rep,name=logs,proto3" json:"logs,omitempty"`                      // The logs generated by this transaction
}

func (x *EthereumTransaction) Reset() {
	*x = EthereumTransaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EthereumTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EthereumTransaction) ProtoMessage() {}

func (x *EthereumTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EthereumTransaction.ProtoReflect.Descriptor instead.
func (*EthereumTransaction) Descriptor() ([]byte, []int) {
	return file_blockchain_ethereum_ethereum_index_types_proto_rawDescGZIP(), []int{1}
}

func (x *EthereumTransaction) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *EthereumTransaction) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *EthereumTransaction) GetFromAddress() string {
	if x != nil {
		return x.FromAddress
	}
	return ""
}

func (x *EthereumTransaction) GetToAddress() string {
	if x != nil {
		return x.ToAddress
	}
	return ""
}

func (x *EthereumTransaction) GetGas() string {
	if x != nil {
		return x.Gas
	}
	return ""
}

func (x *EthereumTransaction) GetGasPrice() string {
	if x != nil {
		return x.GasPrice
	}
	return ""
}

func (x *EthereumTransaction) GetMaxFeePerGas() string {
	if x != nil {
		return x.MaxFeePerGas
	}
	return ""
}

func (x *EthereumTransaction) GetMaxPriorityFeePerGas() string {
	if x != nil {
		return x.MaxPriorityFeePerGas
	}
	return ""
}

func (x *EthereumTransaction) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

func (x *EthereumTransaction) GetNonce() string {
	if x != nil {
		return x.Nonce
	}
	return ""
}

func (x *EthereumTransaction) GetTransactionIndex() uint64 {
	if x != nil {
		return x.TransactionIndex
	}
	return 0
}

func (x *EthereumTransaction) GetTransactionType() uint64 {
	if x != nil {
		return x.TransactionType
	}
	return 0
}

func (x *EthereumTransaction) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *EthereumTransaction) GetIndexedAt() uint64 {
	if x != nil {
		return x.IndexedAt
	}
	return 0
}

func (x *EthereumTransaction) GetBlockTimestamp() uint64 {
	if x != nil {
		return x.BlockTimestamp
	}
	return 0
}

func (x *EthereumTransaction) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *EthereumTransaction) GetChainId() string {
	if x != nil {
		return x.ChainId
	}
	return ""
}

func (x *EthereumTransaction) GetV() string {
	if x != nil {
		return x.V
	}
	return ""
}

func (x *EthereumTransaction) GetR() string {
	if x != nil {
		return x.R
	}
	return ""
}

func (x *EthereumTransaction) GetS() string {
	if x != nil {
		return x.S
	}
	return ""
}

func (x *EthereumTransaction) GetAccessList() []*EthereumTransactionAccessList {
	if x != nil {
		return x.AccessList
	}
	return nil
}

func (x *EthereumTransaction) GetYParity() string {
	if x != nil {
		return x.YParity
	}
	return ""
}

func (x *EthereumTransaction) GetLogs() []*EthereumEventLog {
	if x != nil {
		return x.Logs
	}
	return nil
}

// Represents a single blockchain block
type EthereumBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockNumber      uint64                 `protobuf:"varint,1,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	Difficulty       uint64                 `protobuf:"varint,2,opt,name=difficulty,proto3" json:"difficulty,omitempty"`
	ExtraData        string                 `protobuf:"bytes,3,opt,name=extra_data,json=extraData,proto3" json:"extra_data,omitempty"`
	GasLimit         uint64                 `protobuf:"varint,4,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`
	GasUsed          uint64                 `protobuf:"varint,5,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	BaseFeePerGas    string                 `protobuf:"bytes,6,opt,name=base_fee_per_gas,json=baseFeePerGas,proto3" json:"base_fee_per_gas,omitempty"` // using string to handle big numeric values
	Hash             string                 `protobuf:"bytes,7,opt,name=hash,proto3" json:"hash,omitempty"`
	LogsBloom        string                 `protobuf:"bytes,8,opt,name=logs_bloom,json=logsBloom,proto3" json:"logs_bloom,omitempty"`
	Miner            string                 `protobuf:"bytes,9,opt,name=miner,proto3" json:"miner,omitempty"`
	Nonce            string                 `protobuf:"bytes,10,opt,name=nonce,proto3" json:"nonce,omitempty"`
	ParentHash       string                 `protobuf:"bytes,11,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`
	ReceiptsRoot     string                 `protobuf:"bytes,12,opt,name=receipts_root,json=receiptsRoot,proto3" json:"receipts_root,omitempty"`
	Sha3Uncles       string                 `protobuf:"bytes,13,opt,name=sha3_uncles,json=sha3Uncles,proto3" json:"sha3_uncles,omitempty"`
	Size             uint64                 `protobuf:"varint,14,opt,name=size,proto3" json:"size,omitempty"`
	StateRoot        string                 `protobuf:"bytes,15,opt,name=state_root,json=stateRoot,proto3" json:"state_root,omitempty"`
	Timestamp        uint64                 `protobuf:"varint,16,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	TotalDifficulty  string                 `protobuf:"bytes,17,opt,name=total_difficulty,json=totalDifficulty,proto3" json:"total_difficulty,omitempty"`
	TransactionsRoot string                 `protobuf:"bytes,18,opt,name=transactions_root,json=transactionsRoot,proto3" json:"transactions_root,omitempty"`
	IndexedAt        uint64                 `protobuf:"varint,19,opt,name=indexed_at,json=indexedAt,proto3" json:"indexed_at,omitempty"` // using uint64 to represent timestamp
	Transactions     []*EthereumTransaction `protobuf:"bytes,20,rep,name=transactions,proto3" json:"transactions,omitempty"`
}

func (x *EthereumBlock) Reset() {
	*x = EthereumBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EthereumBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EthereumBlock) ProtoMessage() {}

func (x *EthereumBlock) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EthereumBlock.ProtoReflect.Descriptor instead.
func (*EthereumBlock) Descriptor() ([]byte, []int) {
	return file_blockchain_ethereum_ethereum_index_types_proto_rawDescGZIP(), []int{2}
}

func (x *EthereumBlock) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *EthereumBlock) GetDifficulty() uint64 {
	if x != nil {
		return x.Difficulty
	}
	return 0
}

func (x *EthereumBlock) GetExtraData() string {
	if x != nil {
		return x.ExtraData
	}
	return ""
}

func (x *EthereumBlock) GetGasLimit() uint64 {
	if x != nil {
		return x.GasLimit
	}
	return 0
}

func (x *EthereumBlock) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *EthereumBlock) GetBaseFeePerGas() string {
	if x != nil {
		return x.BaseFeePerGas
	}
	return ""
}

func (x *EthereumBlock) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *EthereumBlock) GetLogsBloom() string {
	if x != nil {
		return x.LogsBloom
	}
	return ""
}

func (x *EthereumBlock) GetMiner() string {
	if x != nil {
		return x.Miner
	}
	return ""
}

func (x *EthereumBlock) GetNonce() string {
	if x != nil {
		return x.Nonce
	}
	return ""
}

func (x *EthereumBlock) GetParentHash() string {
	if x != nil {
		return x.ParentHash
	}
	return ""
}

func (x *EthereumBlock) GetReceiptsRoot() string {
	if x != nil {
		return x.ReceiptsRoot
	}
	return ""
}

func (x *EthereumBlock) GetSha3Uncles() string {
	if x != nil {
		return x.Sha3Uncles
	}
	return ""
}

func (x *EthereumBlock) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *EthereumBlock) GetStateRoot() string {
	if x != nil {
		return x.StateRoot
	}
	return ""
}

func (x *EthereumBlock) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *EthereumBlock) GetTotalDifficulty() string {
	if x != nil {
		return x.TotalDifficulty
	}
	return ""
}

func (x *EthereumBlock) GetTransactionsRoot() string {
	if x != nil {
		return x.TransactionsRoot
	}
	return ""
}

func (x *EthereumBlock) GetIndexedAt() uint64 {
	if x != nil {
		return x.IndexedAt
	}
	return 0
}

func (x *EthereumBlock) GetTransactions() []*EthereumTransaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

type EthereumEventLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address          string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`                                            // The address of the contract that generated the log
	Topics           []string `protobuf:"bytes,2,rep,name=topics,proto3" json:"topics,omitempty"`                                              // Topics are indexed parameters during log generation
	Data             string   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`                                                  // The data field from the log
	BlockNumber      uint64   `protobuf:"varint,4,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`                // The block number where this log was in
	TransactionHash  string   `protobuf:"bytes,5,opt,name=transaction_hash,json=transactionHash,proto3" json:"transaction_hash,omitempty"`     // The hash of the transaction that generated this log
	BlockHash        string   `protobuf:"bytes,6,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`                       // The hash of the block where this log was in
	Removed          bool     `protobuf:"varint,7,opt,name=removed,proto3" json:"removed,omitempty"`                                           // True if the log was reverted due to a chain reorganization
	LogIndex         uint64   `protobuf:"varint,8,opt,name=log_index,json=logIndex,proto3" json:"log_index,omitempty"`                         // The index of the log in the block
	TransactionIndex uint64   `protobuf:"varint,9,opt,name=transaction_index,json=transactionIndex,proto3" json:"transaction_index,omitempty"` // The index of the transaction in the block
}

func (x *EthereumEventLog) Reset() {
	*x = EthereumEventLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EthereumEventLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EthereumEventLog) ProtoMessage() {}

func (x *EthereumEventLog) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EthereumEventLog.ProtoReflect.Descriptor instead.
func (*EthereumEventLog) Descriptor() ([]byte, []int) {
	return file_blockchain_ethereum_ethereum_index_types_proto_rawDescGZIP(), []int{3}
}

func (x *EthereumEventLog) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *EthereumEventLog) GetTopics() []string {
	if x != nil {
		return x.Topics
	}
	return nil
}

func (x *EthereumEventLog) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *EthereumEventLog) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *EthereumEventLog) GetTransactionHash() string {
	if x != nil {
		return x.TransactionHash
	}
	return ""
}

func (x *EthereumEventLog) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *EthereumEventLog) GetRemoved() bool {
	if x != nil {
		return x.Removed
	}
	return false
}

func (x *EthereumEventLog) GetLogIndex() uint64 {
	if x != nil {
		return x.LogIndex
	}
	return 0
}

func (x *EthereumEventLog) GetTransactionIndex() uint64 {
	if x != nil {
		return x.TransactionIndex
	}
	return 0
}

type EthereumBlocksBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Blocks []*EthereumBlock `protobuf:"bytes,1,rep,name=blocks,proto3" json:"blocks,omitempty"`
}

func (x *EthereumBlocksBatch) Reset() {
	*x = EthereumBlocksBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EthereumBlocksBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EthereumBlocksBatch) ProtoMessage() {}

func (x *EthereumBlocksBatch) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EthereumBlocksBatch.ProtoReflect.Descriptor instead.
func (*EthereumBlocksBatch) Descriptor() ([]byte, []int) {
	return file_blockchain_ethereum_ethereum_index_types_proto_rawDescGZIP(), []int{4}
}

func (x *EthereumBlocksBatch) GetBlocks() []*EthereumBlock {
	if x != nil {
		return x.Blocks
	}
	return nil
}

var File_blockchain_ethereum_ethereum_index_types_proto protoreflect.FileDescriptor

var file_blockchain_ethereum_ethereum_index_types_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f, 0x65, 0x74, 0x68,
	0x65, 0x72, 0x65, 0x75, 0x6d, 0x2f, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x5f, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x5c, 0x0a, 0x1d, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4b, 0x65, 0x79, 0x73, 0x22, 0xe5,
	0x05, 0x0a, 0x13, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x21, 0x0a,
	0x0c, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x66, 0x72, 0x6f, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x10, 0x0a, 0x03, 0x67, 0x61, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x67, 0x61,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x61, 0x73, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x25,
	0x0a, 0x0f, 0x6d, 0x61, 0x78, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x67, 0x61,
	0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x46, 0x65, 0x65, 0x50,
	0x65, 0x72, 0x47, 0x61, 0x73, 0x12, 0x36, 0x0a, 0x18, 0x6d, 0x61, 0x78, 0x5f, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x79, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x67, 0x61,
	0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x6d, 0x61, 0x78, 0x50, 0x72, 0x69, 0x6f,
	0x72, 0x69, 0x74, 0x79, 0x46, 0x65, 0x65, 0x50, 0x65, 0x72, 0x47, 0x61, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x65, 0x64, 0x41, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x10, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x19,
	0x0a, 0x08, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x76, 0x18, 0x12,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x76, 0x12, 0x0c, 0x0a, 0x01, 0x72, 0x18, 0x13, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x01, 0x72, 0x12, 0x0c, 0x0a, 0x01, 0x73, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x01, 0x73, 0x12, 0x3f, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6c, 0x69,
	0x73, 0x74, 0x18, 0x15, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x45, 0x74, 0x68, 0x65, 0x72,
	0x65, 0x75, 0x6d, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x79, 0x5f, 0x70, 0x61, 0x72, 0x69, 0x74, 0x79,
	0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x79, 0x50, 0x61, 0x72, 0x69, 0x74, 0x79, 0x12,
	0x25, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x17, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67,
	0x52, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x22, 0x9a, 0x05, 0x0a, 0x0d, 0x45, 0x74, 0x68, 0x65, 0x72,
	0x65, 0x75, 0x6d, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x64,
	0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0a, 0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x65,
	0x78, 0x74, 0x72, 0x61, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x65, 0x78, 0x74, 0x72, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61,
	0x73, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x67,
	0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x61, 0x73, 0x5f, 0x75,
	0x73, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73,
	0x65, 0x64, 0x12, 0x27, 0x0a, 0x10, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x70,
	0x65, 0x72, 0x5f, 0x67, 0x61, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x62, 0x61,
	0x73, 0x65, 0x46, 0x65, 0x65, 0x50, 0x65, 0x72, 0x47, 0x61, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12,
	0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x73, 0x5f, 0x62, 0x6c, 0x6f, 0x6f, 0x6d, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x73, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x12, 0x14,
	0x0a, 0x05, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d,
	0x69, 0x6e, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x52, 0x6f, 0x6f, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x61, 0x33, 0x5f, 0x75, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x68, 0x61, 0x33, 0x55, 0x6e, 0x63, 0x6c, 0x65,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x72,
	0x6f, 0x6f, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x6f, 0x6f, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x10, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x64, 0x69, 0x66, 0x66,
	0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x44, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x12, 0x2b, 0x0a,
	0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x72, 0x6f,
	0x6f, 0x74, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x13, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x64, 0x41, 0x74, 0x12, 0x38, 0x0a, 0x0c, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x14, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x22, 0xa9, 0x02, 0x0a, 0x10, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x21,
	0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x72,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6c, 0x6f, 0x67, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x2b, 0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x22,
	0x3d, 0x0a, 0x13, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x73, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x26, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75,
	0x6d, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x42, 0x33,
	0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x6f,
	0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2d, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x65, 0x72, 0x2f,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f, 0x65, 0x74, 0x68, 0x65, 0x72,
	0x65, 0x75, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_blockchain_ethereum_ethereum_index_types_proto_rawDescOnce sync.Once
	file_blockchain_ethereum_ethereum_index_types_proto_rawDescData = file_blockchain_ethereum_ethereum_index_types_proto_rawDesc
)

func file_blockchain_ethereum_ethereum_index_types_proto_rawDescGZIP() []byte {
	file_blockchain_ethereum_ethereum_index_types_proto_rawDescOnce.Do(func() {
		file_blockchain_ethereum_ethereum_index_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_blockchain_ethereum_ethereum_index_types_proto_rawDescData)
	})
	return file_blockchain_ethereum_ethereum_index_types_proto_rawDescData
}

var file_blockchain_ethereum_ethereum_index_types_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_blockchain_ethereum_ethereum_index_types_proto_goTypes = []interface{}{
	(*EthereumTransactionAccessList)(nil), // 0: EthereumTransactionAccessList
	(*EthereumTransaction)(nil),           // 1: EthereumTransaction
	(*EthereumBlock)(nil),                 // 2: EthereumBlock
	(*EthereumEventLog)(nil),              // 3: EthereumEventLog
	(*EthereumBlocksBatch)(nil),           // 4: EthereumBlocksBatch
}
var file_blockchain_ethereum_ethereum_index_types_proto_depIdxs = []int32{
	0, // 0: EthereumTransaction.access_list:type_name -> EthereumTransactionAccessList
	3, // 1: EthereumTransaction.logs:type_name -> EthereumEventLog
	1, // 2: EthereumBlock.transactions:type_name -> EthereumTransaction
	2, // 3: EthereumBlocksBatch.blocks:type_name -> EthereumBlock
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_blockchain_ethereum_ethereum_index_types_proto_init() }
func file_blockchain_ethereum_ethereum_index_types_proto_init() {
	if File_blockchain_ethereum_ethereum_index_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EthereumTransactionAccessList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EthereumTransaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EthereumBlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EthereumEventLog); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_ethereum_ethereum_index_types_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EthereumBlocksBatch); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_blockchain_ethereum_ethereum_index_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_blockchain_ethereum_ethereum_index_types_proto_goTypes,
		DependencyIndexes: file_blockchain_ethereum_ethereum_index_types_proto_depIdxs,
		MessageInfos:      file_blockchain_ethereum_ethereum_index_types_proto_msgTypes,
	}.Build()
	File_blockchain_ethereum_ethereum_index_types_proto = out.File
	file_blockchain_ethereum_ethereum_index_types_proto_rawDesc = nil
	file_blockchain_ethereum_ethereum_index_types_proto_goTypes = nil
	file_blockchain_ethereum_ethereum_index_types_proto_depIdxs = nil
}
