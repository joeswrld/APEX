package core

import (
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"time"

	"github.com/apex/pkg/types"
)

// TxType represents transaction type
type TxType uint8

const (
	TxTypeTransfer TxType = iota
	TxTypeStake
	TxTypeUnstake
	TxTypeDelegate
	TxTypeUndelegate
	TxTypeVote
	TxTypeCreateValidator
	TxTypeEditValidator
)

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash      types.Hash    `json:"hash"`
	Type      TxType        `json:"type"`
	From      types.Address `json:"from"`
	To        types.Address `json:"to"`
	Value     *big.Int      `json:"value"`
	Data      []byte        `json:"data"`
	Nonce     uint64        `json:"nonce"`
	GasLimit  uint64        `json:"gas_limit"`
	GasPrice  *big.Int      `json:"gas_price"`
	Signature types.Signature `json:"signature"`
	Timestamp time.Time     `json:"timestamp"`
}

// TxReceipt represents a transaction receipt
type TxReceipt struct {
	TxHash          types.Hash    `json:"tx_hash"`
	BlockHash       types.Hash    `json:"block_hash"`
	BlockNumber     uint64        `json:"block_number"`
	From            types.Address `json:"from"`
	To              types.Address `json:"to"`
	GasUsed         uint64        `json:"gas_used"`
	Status          uint8         `json:"status"` // 1 = success, 0 = failure
	Logs            []Log         `json:"logs"`
	ContractAddress types.Address `json:"contract_address,omitempty"`
}

// Log represents an event log
type Log struct {
	Address types.Address `json:"address"`
	Topics  []types.Hash  `json:"topics"`
	Data    []byte        `json:"data"`
}

// NewTransaction creates a new transaction
func NewTransaction(txType TxType, from, to types.Address, value *big.Int, data []byte, nonce uint64) *Transaction {
	return &Transaction{
		Type:      txType,
		From:      from,
		To:        to,
		Value:     value,
		Data:      data,
		Nonce:     nonce,
		GasLimit:  21000, // Base gas
		GasPrice:  big.NewInt(1000000000), // 1 Gwei default
		Timestamp: time.Now(),
	}
}

// ComputeHash computes transaction hash
func (tx *Transaction) ComputeHash() types.Hash {
	data, _ := json.Marshal(struct {
		Type     TxType
		From     types.Address
		To       types.Address
		Value    *big.Int
		Data     []byte
		Nonce    uint64
		GasLimit uint64
		GasPrice *big.Int
	}{
		Type:     tx.Type,
		From:     tx.From,
		To:       tx.To,
		Value:    tx.Value,
		Data:     tx.Data,
		Nonce:    tx.Nonce,
		GasLimit: tx.GasLimit,
		GasPrice: tx.GasPrice,
	})
	
	hash := sha256.Sum256(data)
	var txHash types.Hash
	copy(txHash[:], hash[:])
	return txHash
}

// Sign signs the transaction
func (tx *Transaction) Sign(signature types.Signature) {
	tx.Signature = signature
	tx.Hash = tx.ComputeHash()
}

// GetCost returns total transaction cost (value + gas)
func (tx *Transaction) GetCost() *big.Int {
	cost := new(big.Int).Set(tx.Value)
	gasCost := new(big.Int).Mul(big.NewInt(int64(tx.GasLimit)), tx.GasPrice)
	cost.Add(cost, gasCost)
	return cost
}

// Validate performs basic transaction validation
func (tx *Transaction) Validate() error {
	if tx.Value.Sign() < 0 {
		return ErrInvalidTxValue
	}
	if tx.GasLimit == 0 {
		return ErrInvalidGasLimit
	}
	if tx.GasPrice.Sign() <= 0 {
		return ErrInvalidGasPrice
	}
	if len(tx.Signature) == 0 {
		return ErrMissingSignature
	}
	return nil
}

// StakeData represents stake transaction data
type StakeData struct {
	Validator types.Address `json:"validator"`
	Amount    *big.Int      `json:"amount"`
}

// VoteData represents vote transaction data
type VoteData struct {
	Validator types.Address `json:"validator"`
	Weight    *big.Int      `json:"weight"`
}

// CreateValidatorData represents create validator transaction data
type CreateValidatorData struct {
	PublicKey  []byte   `json:"public_key"`
	Commission uint64   `json:"commission"`
	SelfStake  *big.Int `json:"self_stake"`
	Moniker    string   `json:"moniker"`
	Website    string   `json:"website"`
	Details    string   `json:"details"`
}