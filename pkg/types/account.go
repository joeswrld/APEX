package types

import (
	"math/big"
)

// Account represents a blockchain account
type Account struct {
	Address Address  `json:"address"`
	Balance *big.Int `json:"balance"`
	Nonce   uint64   `json:"nonce"`
	Staked  *big.Int `json:"staked"`  // Amount staked
	Locked  *big.Int `json:"locked"`  // Locked/unbonding tokens
}

// NewAccount creates a new account
func NewAccount(addr Address) *Account {
	return &Account{
		Address: addr,
		Balance: big.NewInt(0),
		Nonce:   0,
		Staked:  big.NewInt(0),
		Locked:  big.NewInt(0),
	}
}

// Copy creates a deep copy of the account
func (a *Account) Copy() *Account {
	return &Account{
		Address: a.Address,
		Balance: new(big.Int).Set(a.Balance),
		Nonce:   a.Nonce,
		Staked:  new(big.Int).Set(a.Staked),
		Locked:  new(big.Int).Set(a.Locked),
	}
}

// AddBalance adds to account balance
func (a *Account) AddBalance(amount *big.Int) {
	a.Balance = new(big.Int).Add(a.Balance, amount)
}

// SubBalance subtracts from account balance
func (a *Account) SubBalance(amount *big.Int) bool {
	if a.Balance.Cmp(amount) < 0 {
		return false
	}
	a.Balance = new(big.Int).Sub(a.Balance, amount)
	return true
}

// AddStake adds to staked amount
func (a *Account) AddStake(amount *big.Int) {
	a.Staked = new(big.Int).Add(a.Staked, amount)
}

// SubStake subtracts from staked amount
func (a *Account) SubStake(amount *big.Int) bool {
	if a.Staked.Cmp(amount) < 0 {
		return false
	}
	a.Staked = new(big.Int).Sub(a.Staked, amount)
	return true
}

// TotalBalance returns total balance including staked and locked
func (a *Account) TotalBalance() *big.Int {
	total := new(big.Int).Set(a.Balance)
	total.Add(total, a.Staked)
	total.Add(total, a.Locked)
	return total
}