// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.6.1
// source: blockchain/xai/xai_index_types.proto

package xai

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

type XaiTransactionAccessList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address     string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	StorageKeys []string `protobuf:"bytes,2,rep,name=storage_keys,json=storageKeys,proto3" json:"storage_keys,omitempty"`
}

func (x *XaiTransactionAccessList) Reset() {
	*x = XaiTransactionAccessList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XaiTransactionAccessList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XaiTransactionAccessList) ProtoMessage() {}

func (x *XaiTransactionAccessList) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XaiTransactionAccessList.ProtoReflect.Descriptor instead.
func (*XaiTransactionAccessList) Descriptor() ([]byte, []int) {
	return file_blockchain_xai_xai_index_types_proto_rawDescGZIP(), []int{0}
}

func (x *XaiTransactionAccessList) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *XaiTransactionAccessList) GetStorageKeys() []string {
	if x != nil {
		return x.StorageKeys
	}
	return nil
}

// Represents a single transaction within a block
type XaiTransaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash                 string                      `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`                                                                   // The hash of the transaction
	BlockNumber          uint64                      `protobuf:"varint,2,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`                                 // The block number the transaction is in
	FromAddress          string                      `protobuf:"bytes,3,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`                                  // The address the transaction is sent from
	ToAddress            string                      `protobuf:"bytes,4,opt,name=to_address,json=toAddress,proto3" json:"to_address,omitempty"`                                        // The address the transaction is sent to
	Gas                  string                      `protobuf:"bytes,5,opt,name=gas,proto3" json:"gas,omitempty"`                                                                     // The gas limit of the transaction
	GasPrice             string                      `protobuf:"bytes,6,opt,name=gas_price,json=gasPrice,proto3" json:"gas_price,omitempty"`                                           // The gas price of the transaction
	MaxFeePerGas         string                      `protobuf:"bytes,7,opt,name=max_fee_per_gas,json=maxFeePerGas,proto3" json:"max_fee_per_gas,omitempty"`                           // Used as a field to match potential EIP-1559 transaction types
	MaxPriorityFeePerGas string                      `protobuf:"bytes,8,opt,name=max_priority_fee_per_gas,json=maxPriorityFeePerGas,proto3" json:"max_priority_fee_per_gas,omitempty"` // Used as a field to match potential EIP-1559 transaction types
	Input                string                      `protobuf:"bytes,9,opt,name=input,proto3" json:"input,omitempty"`                                                                 // The input data of the transaction
	Nonce                string                      `protobuf:"bytes,10,opt,name=nonce,proto3" json:"nonce,omitempty"`                                                                // The nonce of the transaction
	TransactionIndex     uint64                      `protobuf:"varint,11,opt,name=transaction_index,json=transactionIndex,proto3" json:"transaction_index,omitempty"`                 // The index of the transaction in the block
	TransactionType      uint64                      `protobuf:"varint,12,opt,name=transaction_type,json=transactionType,proto3" json:"transaction_type,omitempty"`                    // Field to match potential EIP-1559 transaction types
	Value                string                      `protobuf:"bytes,13,opt,name=value,proto3" json:"value,omitempty"`                                                                // The value of the transaction
	IndexedAt            uint64                      `protobuf:"varint,14,opt,name=indexed_at,json=indexedAt,proto3" json:"indexed_at,omitempty"`                                      // When the transaction was indexed by crawler
	BlockTimestamp       uint64                      `protobuf:"varint,15,opt,name=block_timestamp,json=blockTimestamp,proto3" json:"block_timestamp,omitempty"`                       // The timestamp of this block
	BlockHash            string                      `protobuf:"bytes,16,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`                                       // The hash of the block the transaction is in
	ChainId              string                      `protobuf:"bytes,17,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`                                             // Used as a field to match potential EIP-1559 transaction types
	V                    string                      `protobuf:"bytes,18,opt,name=v,proto3" json:"v,omitempty"`                                                                        // Used as a field to match potential EIP-1559 transaction types
	R                    string                      `protobuf:"bytes,19,opt,name=r,proto3" json:"r,omitempty"`                                                                        // Used as a field to match potential EIP-1559 transaction types
	S                    string                      `protobuf:"bytes,20,opt,name=s,proto3" json:"s,omitempty"`                                                                        // Used as a field to match potential EIP-1559 transaction types
	AccessList           []*XaiTransactionAccessList `protobuf:"bytes,21,rep,name=access_list,json=accessList,proto3" json:"access_list,omitempty"`
	YParity              string                      `protobuf:"bytes,22,opt,name=y_parity,json=yParity,proto3" json:"y_parity,omitempty"` // Used as a field to match potential EIP-1559 transaction types
	Logs                 []*XaiEventLog              `protobuf:"bytes,23,rep,name=logs,proto3" json:"logs,omitempty"`                      // The logs generated by this transaction
}

func (x *XaiTransaction) Reset() {
	*x = XaiTransaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XaiTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XaiTransaction) ProtoMessage() {}

func (x *XaiTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XaiTransaction.ProtoReflect.Descriptor instead.
func (*XaiTransaction) Descriptor() ([]byte, []int) {
	return file_blockchain_xai_xai_index_types_proto_rawDescGZIP(), []int{1}
}

func (x *XaiTransaction) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *XaiTransaction) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *XaiTransaction) GetFromAddress() string {
	if x != nil {
		return x.FromAddress
	}
	return ""
}

func (x *XaiTransaction) GetToAddress() string {
	if x != nil {
		return x.ToAddress
	}
	return ""
}

func (x *XaiTransaction) GetGas() string {
	if x != nil {
		return x.Gas
	}
	return ""
}

func (x *XaiTransaction) GetGasPrice() string {
	if x != nil {
		return x.GasPrice
	}
	return ""
}

func (x *XaiTransaction) GetMaxFeePerGas() string {
	if x != nil {
		return x.MaxFeePerGas
	}
	return ""
}

func (x *XaiTransaction) GetMaxPriorityFeePerGas() string {
	if x != nil {
		return x.MaxPriorityFeePerGas
	}
	return ""
}

func (x *XaiTransaction) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

func (x *XaiTransaction) GetNonce() string {
	if x != nil {
		return x.Nonce
	}
	return ""
}

func (x *XaiTransaction) GetTransactionIndex() uint64 {
	if x != nil {
		return x.TransactionIndex
	}
	return 0
}

func (x *XaiTransaction) GetTransactionType() uint64 {
	if x != nil {
		return x.TransactionType
	}
	return 0
}

func (x *XaiTransaction) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *XaiTransaction) GetIndexedAt() uint64 {
	if x != nil {
		return x.IndexedAt
	}
	return 0
}

func (x *XaiTransaction) GetBlockTimestamp() uint64 {
	if x != nil {
		return x.BlockTimestamp
	}
	return 0
}

func (x *XaiTransaction) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *XaiTransaction) GetChainId() string {
	if x != nil {
		return x.ChainId
	}
	return ""
}

func (x *XaiTransaction) GetV() string {
	if x != nil {
		return x.V
	}
	return ""
}

func (x *XaiTransaction) GetR() string {
	if x != nil {
		return x.R
	}
	return ""
}

func (x *XaiTransaction) GetS() string {
	if x != nil {
		return x.S
	}
	return ""
}

func (x *XaiTransaction) GetAccessList() []*XaiTransactionAccessList {
	if x != nil {
		return x.AccessList
	}
	return nil
}

func (x *XaiTransaction) GetYParity() string {
	if x != nil {
		return x.YParity
	}
	return ""
}

func (x *XaiTransaction) GetLogs() []*XaiEventLog {
	if x != nil {
		return x.Logs
	}
	return nil
}

// Represents a block in the blockchain
type XaiBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockNumber      uint64            `protobuf:"varint,1,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`          // The block number
	Difficulty       uint64            `protobuf:"varint,2,opt,name=difficulty,proto3" json:"difficulty,omitempty"`                               // The difficulty of this block
	ExtraData        string            `protobuf:"bytes,3,opt,name=extra_data,json=extraData,proto3" json:"extra_data,omitempty"`                 // Extra data included in the block
	GasLimit         uint64            `protobuf:"varint,4,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`                   // The gas limit for this block
	GasUsed          uint64            `protobuf:"varint,5,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`                      // The total gas used by all transactions in this block
	BaseFeePerGas    string            `protobuf:"bytes,6,opt,name=base_fee_per_gas,json=baseFeePerGas,proto3" json:"base_fee_per_gas,omitempty"` // The base fee per gas for this block
	Hash             string            `protobuf:"bytes,7,opt,name=hash,proto3" json:"hash,omitempty"`                                            // The hash of this block
	LogsBloom        string            `protobuf:"bytes,8,opt,name=logs_bloom,json=logsBloom,proto3" json:"logs_bloom,omitempty"`                 // The logs bloom filter for this block
	Miner            string            `protobuf:"bytes,9,opt,name=miner,proto3" json:"miner,omitempty"`                                          // The address of the miner who mined this block
	Nonce            string            `protobuf:"bytes,10,opt,name=nonce,proto3" json:"nonce,omitempty"`                                         // The nonce of this block
	ParentHash       string            `protobuf:"bytes,11,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`             // The hash of the parent block
	ReceiptsRoot     string            `protobuf:"bytes,12,opt,name=receipts_root,json=receiptsRoot,proto3" json:"receipts_root,omitempty"`       // The root hash of the receipts trie
	Sha3Uncles       string            `protobuf:"bytes,13,opt,name=sha3_uncles,json=sha3Uncles,proto3" json:"sha3_uncles,omitempty"`             // The SHA3 hash of the uncles data in this block
	Size             uint64            `protobuf:"varint,14,opt,name=size,proto3" json:"size,omitempty"`                                          // The size of this block
	StateRoot        string            `protobuf:"bytes,15,opt,name=state_root,json=stateRoot,proto3" json:"state_root,omitempty"`                // The root hash of the state trie
	Timestamp        uint64            `protobuf:"varint,16,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	TotalDifficulty  string            `protobuf:"bytes,17,opt,name=total_difficulty,json=totalDifficulty,proto3" json:"total_difficulty,omitempty"`    // The total difficulty of the chain until this block
	TransactionsRoot string            `protobuf:"bytes,18,opt,name=transactions_root,json=transactionsRoot,proto3" json:"transactions_root,omitempty"` // The root hash of the transactions trie
	IndexedAt        uint64            `protobuf:"varint,19,opt,name=indexed_at,json=indexedAt,proto3" json:"indexed_at,omitempty"`                     // When the block was indexed by crawler
	Transactions     []*XaiTransaction `protobuf:"bytes,20,rep,name=transactions,proto3" json:"transactions,omitempty"`                                 // The transactions included in this block
	MixHash          string            `protobuf:"bytes,21,opt,name=mix_hash,json=mixHash,proto3" json:"mix_hash,omitempty"`                            // The timestamp of this block
	SendCount        string            `protobuf:"bytes,22,opt,name=send_count,json=sendCount,proto3" json:"send_count,omitempty"`                      // The number of sends in this block
	SendRoot         string            `protobuf:"bytes,23,opt,name=send_root,json=sendRoot,proto3" json:"send_root,omitempty"`                         // The root hash of the sends trie
	L1BlockNumber    uint64            `protobuf:"varint,24,opt,name=l1_block_number,json=l1BlockNumber,proto3" json:"l1_block_number,omitempty"`       // The block number of the corresponding L1 block
}

func (x *XaiBlock) Reset() {
	*x = XaiBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XaiBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XaiBlock) ProtoMessage() {}

func (x *XaiBlock) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XaiBlock.ProtoReflect.Descriptor instead.
func (*XaiBlock) Descriptor() ([]byte, []int) {
	return file_blockchain_xai_xai_index_types_proto_rawDescGZIP(), []int{2}
}

func (x *XaiBlock) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *XaiBlock) GetDifficulty() uint64 {
	if x != nil {
		return x.Difficulty
	}
	return 0
}

func (x *XaiBlock) GetExtraData() string {
	if x != nil {
		return x.ExtraData
	}
	return ""
}

func (x *XaiBlock) GetGasLimit() uint64 {
	if x != nil {
		return x.GasLimit
	}
	return 0
}

func (x *XaiBlock) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *XaiBlock) GetBaseFeePerGas() string {
	if x != nil {
		return x.BaseFeePerGas
	}
	return ""
}

func (x *XaiBlock) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *XaiBlock) GetLogsBloom() string {
	if x != nil {
		return x.LogsBloom
	}
	return ""
}

func (x *XaiBlock) GetMiner() string {
	if x != nil {
		return x.Miner
	}
	return ""
}

func (x *XaiBlock) GetNonce() string {
	if x != nil {
		return x.Nonce
	}
	return ""
}

func (x *XaiBlock) GetParentHash() string {
	if x != nil {
		return x.ParentHash
	}
	return ""
}

func (x *XaiBlock) GetReceiptsRoot() string {
	if x != nil {
		return x.ReceiptsRoot
	}
	return ""
}

func (x *XaiBlock) GetSha3Uncles() string {
	if x != nil {
		return x.Sha3Uncles
	}
	return ""
}

func (x *XaiBlock) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *XaiBlock) GetStateRoot() string {
	if x != nil {
		return x.StateRoot
	}
	return ""
}

func (x *XaiBlock) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *XaiBlock) GetTotalDifficulty() string {
	if x != nil {
		return x.TotalDifficulty
	}
	return ""
}

func (x *XaiBlock) GetTransactionsRoot() string {
	if x != nil {
		return x.TransactionsRoot
	}
	return ""
}

func (x *XaiBlock) GetIndexedAt() uint64 {
	if x != nil {
		return x.IndexedAt
	}
	return 0
}

func (x *XaiBlock) GetTransactions() []*XaiTransaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *XaiBlock) GetMixHash() string {
	if x != nil {
		return x.MixHash
	}
	return ""
}

func (x *XaiBlock) GetSendCount() string {
	if x != nil {
		return x.SendCount
	}
	return ""
}

func (x *XaiBlock) GetSendRoot() string {
	if x != nil {
		return x.SendRoot
	}
	return ""
}

func (x *XaiBlock) GetL1BlockNumber() uint64 {
	if x != nil {
		return x.L1BlockNumber
	}
	return 0
}

type XaiEventLog struct {
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

func (x *XaiEventLog) Reset() {
	*x = XaiEventLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XaiEventLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XaiEventLog) ProtoMessage() {}

func (x *XaiEventLog) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XaiEventLog.ProtoReflect.Descriptor instead.
func (*XaiEventLog) Descriptor() ([]byte, []int) {
	return file_blockchain_xai_xai_index_types_proto_rawDescGZIP(), []int{3}
}

func (x *XaiEventLog) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *XaiEventLog) GetTopics() []string {
	if x != nil {
		return x.Topics
	}
	return nil
}

func (x *XaiEventLog) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *XaiEventLog) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *XaiEventLog) GetTransactionHash() string {
	if x != nil {
		return x.TransactionHash
	}
	return ""
}

func (x *XaiEventLog) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *XaiEventLog) GetRemoved() bool {
	if x != nil {
		return x.Removed
	}
	return false
}

func (x *XaiEventLog) GetLogIndex() uint64 {
	if x != nil {
		return x.LogIndex
	}
	return 0
}

func (x *XaiEventLog) GetTransactionIndex() uint64 {
	if x != nil {
		return x.TransactionIndex
	}
	return 0
}

type XaiBlocksBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Blocks      []*XaiBlock `protobuf:"bytes,1,rep,name=blocks,proto3" json:"blocks,omitempty"`
	SeerVersion string      `protobuf:"bytes,2,opt,name=seer_version,json=seerVersion,proto3" json:"seer_version,omitempty"`
}

func (x *XaiBlocksBatch) Reset() {
	*x = XaiBlocksBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *XaiBlocksBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*XaiBlocksBatch) ProtoMessage() {}

func (x *XaiBlocksBatch) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_xai_xai_index_types_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use XaiBlocksBatch.ProtoReflect.Descriptor instead.
func (*XaiBlocksBatch) Descriptor() ([]byte, []int) {
	return file_blockchain_xai_xai_index_types_proto_rawDescGZIP(), []int{4}
}

func (x *XaiBlocksBatch) GetBlocks() []*XaiBlock {
	if x != nil {
		return x.Blocks
	}
	return nil
}

func (x *XaiBlocksBatch) GetSeerVersion() string {
	if x != nil {
		return x.SeerVersion
	}
	return ""
}

var File_blockchain_xai_xai_index_types_proto protoreflect.FileDescriptor

var file_blockchain_xai_xai_index_types_proto_rawDesc = []byte{
	0x0a, 0x24, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f, 0x78, 0x61, 0x69,
	0x2f, 0x78, 0x61, 0x69, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x57, 0x0a, 0x18, 0x58, 0x61, 0x69, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x21, 0x0a, 0x0c,
	0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4b, 0x65, 0x79, 0x73, 0x22,
	0xd6, 0x05, 0x0a, 0x0e, 0x58, 0x61, 0x69, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f,
	0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x72, 0x6f,
	0x6d, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x66, 0x72, 0x6f, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1d, 0x0a, 0x0a,
	0x74, 0x6f, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x74, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x67,
	0x61, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x67, 0x61, 0x73, 0x12, 0x1b, 0x0a,
	0x09, 0x67, 0x61, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x67, 0x61, 0x73, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x25, 0x0a, 0x0f, 0x6d, 0x61,
	0x78, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x67, 0x61, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x46, 0x65, 0x65, 0x50, 0x65, 0x72, 0x47, 0x61,
	0x73, 0x12, 0x36, 0x0a, 0x18, 0x6d, 0x61, 0x78, 0x5f, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74,
	0x79, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x67, 0x61, 0x73, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x14, 0x6d, 0x61, 0x78, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79,
	0x46, 0x65, 0x65, 0x50, 0x65, 0x72, 0x47, 0x61, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x76, 0x18, 0x12, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x01, 0x76, 0x12, 0x0c, 0x0a, 0x01, 0x72, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01,
	0x72, 0x12, 0x0c, 0x0a, 0x01, 0x73, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x73, 0x12,
	0x3a, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x15,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x58, 0x61, 0x69, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x79,
	0x5f, 0x70, 0x61, 0x72, 0x69, 0x74, 0x79, 0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x79,
	0x50, 0x61, 0x72, 0x69, 0x74, 0x79, 0x12, 0x20, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x17,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x58, 0x61, 0x69, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c,
	0x6f, 0x67, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x22, 0x8f, 0x06, 0x0a, 0x08, 0x58, 0x61, 0x69,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x69, 0x66, 0x66,
	0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x64, 0x69,
	0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x74, 0x72,
	0x61, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x78,
	0x74, 0x72, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61, 0x73, 0x5f, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x67, 0x61, 0x73, 0x4c,
	0x69, 0x6d, 0x69, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x61, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x67, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12,
	0x27, 0x0a, 0x10, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x5f,
	0x67, 0x61, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x62, 0x61, 0x73, 0x65, 0x46,
	0x65, 0x65, 0x50, 0x65, 0x72, 0x47, 0x61, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a,
	0x6c, 0x6f, 0x67, 0x73, 0x5f, 0x62, 0x6c, 0x6f, 0x6f, 0x6d, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6c, 0x6f, 0x67, 0x73, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x6d,
	0x69, 0x6e, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x69, 0x6e, 0x65,
	0x72, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65,
	0x69, 0x70, 0x74, 0x73, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x73, 0x68, 0x61, 0x33, 0x5f, 0x75, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x73, 0x68, 0x61, 0x33, 0x55, 0x6e, 0x63, 0x6c, 0x65, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x6f, 0x6f, 0x74,
	0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x10,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x29, 0x0a, 0x10, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75,
	0x6c, 0x74, 0x79, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x44, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x12, 0x2b, 0x0a, 0x11, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x18,
	0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x13, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x65, 0x64, 0x41, 0x74, 0x12, 0x33, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x14, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x58,
	0x61, 0x69, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x6d,
	0x69, 0x78, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x69, 0x78, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x6e, 0x64,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x6f,
	0x6f, 0x74, 0x18, 0x17, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x52, 0x6f,
	0x6f, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x6c, 0x31, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x18, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x6c, 0x31, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0xa4, 0x02, 0x0a, 0x0b, 0x58,
	0x61, 0x69, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d,
	0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x18, 0x0a,
	0x07, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x5f, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6c, 0x6f, 0x67, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x12, 0x2b, 0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x22, 0x56, 0x0a, 0x0e, 0x58, 0x61, 0x69, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x12, 0x21, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x58, 0x61, 0x69, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x06,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x65, 0x72, 0x5f, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65,
	0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x47, 0x37, 0x44, 0x41, 0x4f, 0x2f, 0x73, 0x65,
	0x65, 0x72, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f, 0x78, 0x61,
	0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_blockchain_xai_xai_index_types_proto_rawDescOnce sync.Once
	file_blockchain_xai_xai_index_types_proto_rawDescData = file_blockchain_xai_xai_index_types_proto_rawDesc
)

func file_blockchain_xai_xai_index_types_proto_rawDescGZIP() []byte {
	file_blockchain_xai_xai_index_types_proto_rawDescOnce.Do(func() {
		file_blockchain_xai_xai_index_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_blockchain_xai_xai_index_types_proto_rawDescData)
	})
	return file_blockchain_xai_xai_index_types_proto_rawDescData
}

var file_blockchain_xai_xai_index_types_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_blockchain_xai_xai_index_types_proto_goTypes = []any{
	(*XaiTransactionAccessList)(nil), // 0: XaiTransactionAccessList
	(*XaiTransaction)(nil),           // 1: XaiTransaction
	(*XaiBlock)(nil),                 // 2: XaiBlock
	(*XaiEventLog)(nil),              // 3: XaiEventLog
	(*XaiBlocksBatch)(nil),           // 4: XaiBlocksBatch
}
var file_blockchain_xai_xai_index_types_proto_depIdxs = []int32{
	0, // 0: XaiTransaction.access_list:type_name -> XaiTransactionAccessList
	3, // 1: XaiTransaction.logs:type_name -> XaiEventLog
	1, // 2: XaiBlock.transactions:type_name -> XaiTransaction
	2, // 3: XaiBlocksBatch.blocks:type_name -> XaiBlock
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_blockchain_xai_xai_index_types_proto_init() }
func file_blockchain_xai_xai_index_types_proto_init() {
	if File_blockchain_xai_xai_index_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_blockchain_xai_xai_index_types_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*XaiTransactionAccessList); i {
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
		file_blockchain_xai_xai_index_types_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*XaiTransaction); i {
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
		file_blockchain_xai_xai_index_types_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*XaiBlock); i {
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
		file_blockchain_xai_xai_index_types_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*XaiEventLog); i {
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
		file_blockchain_xai_xai_index_types_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*XaiBlocksBatch); i {
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
			RawDescriptor: file_blockchain_xai_xai_index_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_blockchain_xai_xai_index_types_proto_goTypes,
		DependencyIndexes: file_blockchain_xai_xai_index_types_proto_depIdxs,
		MessageInfos:      file_blockchain_xai_xai_index_types_proto_msgTypes,
	}.Build()
	File_blockchain_xai_xai_index_types_proto = out.File
	file_blockchain_xai_xai_index_types_proto_rawDesc = nil
	file_blockchain_xai_xai_index_types_proto_goTypes = nil
	file_blockchain_xai_xai_index_types_proto_depIdxs = nil
}
