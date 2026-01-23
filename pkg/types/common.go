package types

import (
	"encoding/hex"
	"math/big"
)

// Hash represents a 32-byte hash
type Hash [32]byte

// Address represents a 20-byte account address
type Address [20]byte

// Signature represents a cryptographic signature
type Signature []byte

// Constants
const (
	// TotalSupply is the total supply of APX tokens (500M)
	TotalSupply = 500_000_000
	
	// Decimals for APX token (18 decimals like ETH)
	Decimals = 18
	
	// BlockTime target in seconds
	BlockTime = 3
	
	// MaxValidators is the maximum number of active validators
	MaxValidators = 21
	
	// MinStakeAmount is minimum stake to become a validator (100K APX)
	MinStakeAmount = 100_000
	
	// EpochLength is blocks per epoch for validator rotation
	EpochLength = 1200 // ~1 hour at 3s blocks
	
	// UnbondingPeriod in blocks (~7 days)
	UnbondingPeriod = 201_600
)

// ToWei converts APX amount to wei (smallest unit)
func ToWei(apx float64) *big.Int {
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(Decimals), nil)
	amount := new(big.Float).Mul(big.NewFloat(apx), new(big.Float).SetInt(multiplier))
	result, _ := amount.Int(nil)
	return result
}

// FromWei converts wei to APX
func FromWei(wei *big.Int) float64 {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(Decimals), nil)
	result := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(divisor))
	f, _ := result.Float64()
	return f
}

// HexToHash converts hex string to Hash
func HexToHash(s string) Hash {
	var h Hash
	b, _ := hex.DecodeString(s)
	copy(h[:], b)
	return h
}

// HexToAddress converts hex string to Address
func HexToAddress(s string) Address {
	var a Address
	b, _ := hex.DecodeString(s)
	copy(a[:], b)
	return a
}

// Hex returns hex encoding of hash
func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

// Hex returns hex encoding of address
func (a Address) Hex() string {
	return hex.EncodeToString(a[:])
}

// Bytes returns byte representation
func (h Hash) Bytes() []byte {
	return h[:]
}

// Bytes returns byte representation
func (a Address) Bytes() []byte {
	return a[:]
}