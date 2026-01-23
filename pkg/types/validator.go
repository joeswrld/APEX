package types

import (
	"math/big"
	"time"
)

// ValidatorStatus represents validator state
type ValidatorStatus int

const (
	ValidatorStatusActive ValidatorStatus = iota
	ValidatorStatusInactive
	ValidatorStatusJailed
	ValidatorStatusUnbonding
)

// Validator represents a network validator
type Validator struct {
	Address           Address         `json:"address"`
	PublicKey         []byte          `json:"public_key"`
	VotingPower       *big.Int        `json:"voting_power"`      // Total delegated stake
	SelfStake         *big.Int        `json:"self_stake"`        // Validator's own stake
	Commission        uint64          `json:"commission"`        // Commission rate (basis points, 10000 = 100%)
	Status            ValidatorStatus `json:"status"`
	Jailed            bool            `json:"jailed"`
	JailTime          time.Time       `json:"jail_time"`
	MissedBlocks      uint64          `json:"missed_blocks"`
	ProducedBlocks    uint64          `json:"produced_blocks"`
	LastActiveEpoch   uint64          `json:"last_active_epoch"`
	TotalRewards      *big.Int        `json:"total_rewards"`
	CreatedAt         time.Time       `json:"created_at"`
}

// Delegation represents a stake delegation
type Delegation struct {
	Delegator   Address  `json:"delegator"`
	Validator   Address  `json:"validator"`
	Amount      *big.Int `json:"amount"`
	Rewards     *big.Int `json:"rewards"`
	CreatedAt   time.Time `json:"created_at"`
}

// UnbondingDelegation represents an unbonding delegation
type UnbondingDelegation struct {
	Delegator       Address  `json:"delegator"`
	Validator       Address  `json:"validator"`
	Amount          *big.Int `json:"amount"`
	CompletionBlock uint64   `json:"completion_block"`
	CreatedAt       time.Time `json:"created_at"`
}

// NewValidator creates a new validator
func NewValidator(addr Address, pubKey []byte, selfStake *big.Int, commission uint64) *Validator {
	return &Validator{
		Address:         addr,
		PublicKey:       pubKey,
		VotingPower:     new(big.Int).Set(selfStake),
		SelfStake:       new(big.Int).Set(selfStake),
		Commission:      commission,
		Status:          ValidatorStatusActive,
		Jailed:          false,
		MissedBlocks:    0,
		ProducedBlocks:  0,
		TotalRewards:    big.NewInt(0),
		CreatedAt:       time.Now(),
	}
}

// IsActive returns true if validator is active
func (v *Validator) IsActive() bool {
	return v.Status == ValidatorStatusActive && !v.Jailed
}

// CanProduceBlocks returns true if validator can produce blocks
func (v *Validator) CanProduceBlocks() bool {
	return v.IsActive() && v.VotingPower.Cmp(ToWei(float64(MinStakeAmount))) >= 0
}

// AddVotingPower adds to validator's voting power
func (v *Validator) AddVotingPower(amount *big.Int) {
	v.VotingPower = new(big.Int).Add(v.VotingPower, amount)
}

// SubVotingPower subtracts from validator's voting power
func (v *Validator) SubVotingPower(amount *big.Int) {
	v.VotingPower = new(big.Int).Sub(v.VotingPower, amount)
	if v.VotingPower.Sign() < 0 {
		v.VotingPower = big.NewInt(0)
	}
}

// CalculateCommission calculates commission from rewards
func (v *Validator) CalculateCommission(rewards *big.Int) *big.Int {
	commission := new(big.Int).Mul(rewards, big.NewInt(int64(v.Commission)))
	commission.Div(commission, big.NewInt(10000))
	return commission
}

// NewDelegation creates a new delegation
func NewDelegation(delegator, validator Address, amount *big.Int) *Delegation {
	return &Delegation{
		Delegator: delegator,
		Validator: validator,
		Amount:    new(big.Int).Set(amount),
		Rewards:   big.NewInt(0),
		CreatedAt: time.Now(),
	}
}

// AddRewards adds rewards to delegation
func (d *Delegation) AddRewards(amount *big.Int) {
	d.Rewards = new(big.Int).Add(d.Rewards, amount)
}