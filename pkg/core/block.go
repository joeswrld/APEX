package core

import (
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"time"

	"github.com/apex-blockchain/apex/pkg/types"
)

// BlockHeader represents a block header
type BlockHeader struct {
	Number          uint64        `json:"number"`
	PreviousHash    types.Hash    `json:"previous_hash"`
	Timestamp       time.Time     `json:"timestamp"`
	TransactionRoot types.Hash    `json:"transaction_root"`
	StateRoot       types.Hash    `json:"state_root"`
	Validator       types.Address `json:"validator"`
	Signature       types.Signature `json:"signature"`
	GasUsed         uint64        `json:"gas_used"`
	GasLimit        uint64        `json:"gas_limit"`
}

// Block represents a blockchain block
type Block struct {
	Header       *BlockHeader   `json:"header"`
	Transactions []*Transaction `json:"transactions"`
	Hash         types.Hash     `json:"hash"`
}

// NewBlock creates a new block
func NewBlock(number uint64, previousHash types.Hash, validator types.Address) *Block {
	return &Block{
		Header: &BlockHeader{
			Number:       number,
			PreviousHash: previousHash,
			Timestamp:    time.Now(),
			Validator:    validator,
			GasLimit:     10_000_000, // 10M gas per block
			GasUsed:      0,
		},
		Transactions: make([]*Transaction, 0),
	}
}

// AddTransaction adds a transaction to the block
func (b *Block) AddTransaction(tx *Transaction) bool {
	// Check gas limit
	txGas := tx.GasLimit
	if b.Header.GasUsed+txGas > b.Header.GasLimit {
		return false
	}
	
	b.Transactions = append(b.Transactions, tx)
	b.Header.GasUsed += txGas
	return true
}

// ComputeHash computes block hash
func (b *Block) ComputeHash() types.Hash {
	headerData, _ := json.Marshal(b.Header)
	hash := sha256.Sum256(headerData)
	var blockHash types.Hash
	copy(blockHash[:], hash[:])
	return blockHash
}

// ComputeTransactionRoot computes merkle root of transactions
func (b *Block) ComputeTransactionRoot() types.Hash {
	if len(b.Transactions) == 0 {
		return types.Hash{}
	}
	
	// Simple implementation - hash all transaction hashes together
	var allHashes []byte
	for _, tx := range b.Transactions {
		allHashes = append(allHashes, tx.Hash[:]...)
	}
	
	hash := sha256.Sum256(allHashes)
	var root types.Hash
	copy(root[:], hash[:])
	return root
}

// Finalize finalizes the block (compute roots and hash)
func (b *Block) Finalize(stateRoot types.Hash) {
	b.Header.TransactionRoot = b.ComputeTransactionRoot()
	b.Header.StateRoot = stateRoot
	b.Hash = b.ComputeHash()
}

// Sign signs the block
func (b *Block) Sign(signature types.Signature) {
	b.Header.Signature = signature
	b.Hash = b.ComputeHash()
}

// Validate performs basic block validation
func (b *Block) Validate() error {
	// Check timestamp
	if b.Header.Timestamp.After(time.Now().Add(time.Minute)) {
		return ErrFutureBlock
	}
	
	// Check gas usage
	if b.Header.GasUsed > b.Header.GasLimit {
		return ErrGasLimitExceeded
	}
	
	// Validate transactions
	for _, tx := range b.Transactions {
		if err := tx.Validate(); err != nil {
			return err
		}
	}
	
	// Check signature
	if len(b.Header.Signature) == 0 {
		return ErrMissingSignature
	}
	
	return nil
}

// GetTotalValue returns total value transferred in block
func (b *Block) GetTotalValue() *big.Int {
	total := big.NewInt(0)
	for _, tx := range b.Transactions {
		total.Add(total, tx.Value)
	}
	return total
}

// GetTransactionCount returns number of transactions
func (b *Block) GetTransactionCount() int {
	return len(b.Transactions)
}

// Errors
var (
	ErrFutureBlock       = &BlockError{msg: "block timestamp is in the future"}
	ErrGasLimitExceeded  = &BlockError{msg: "block gas limit exceeded"}
	ErrInvalidTxValue    = &BlockError{msg: "invalid transaction value"}
	ErrInvalidGasLimit   = &BlockError{msg: "invalid gas limit"}
	ErrInvalidGasPrice   = &BlockError{msg: "invalid gas price"}
	ErrMissingSignature  = &BlockError{msg: "missing signature"}
	ErrInvalidBlockHash  = &BlockError{msg: "invalid block hash"}
	ErrInvalidValidator  = &BlockError{msg: "invalid validator"}
)

type BlockError struct {
	msg string
}

func (e *BlockError) Error() string {
	return e.msg
}