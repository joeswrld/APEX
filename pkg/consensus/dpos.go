package consensus

import (
	"errors"
	"math/big"
	"sort"
	"sync"
	"time"

	
	"github.com/apex/pkg/types"
)

// DPoS implements Delegated Proof of Stake consensus
type DPoS struct {
	validators       map[types.Address]*types.Validator
	delegations      map[types.Address]map[types.Address]*types.Delegation // delegator -> validator -> delegation
	activeValidators []*types.Validator
	currentEpoch     uint64
	mu               sync.RWMutex
}

// NewDPoS creates a new DPoS consensus engine
func NewDPoS() *DPoS {
	return &DPoS{
		validators:       make(map[types.Address]*types.Validator),
		delegations:      make(map[types.Address]map[types.Address]*types.Delegation),
		activeValidators: make([]*types.Validator, 0),
		currentEpoch:     0,
	}
}

// RegisterValidator registers a new validator
func (d *DPoS) RegisterValidator(validator *types.Validator) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	if _, exists := d.validators[validator.Address]; exists {
		return errors.New("validator already registered")
	}
	
	// Check minimum stake requirement
	minStake := types.ToWei(float64(types.MinStakeAmount))
	if validator.SelfStake.Cmp(minStake) < 0 {
		return errors.New("insufficient self-stake")
	}
	
	d.validators[validator.Address] = validator
	return nil
}

// Delegate allows a user to delegate stake to a validator
func (d *DPoS) Delegate(delegator, validator types.Address, amount *big.Int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Check validator exists
	val, exists := d.validators[validator]
	if !exists {
		return errors.New("validator not found")
	}
	
	// Create or update delegation
	if d.delegations[delegator] == nil {
		d.delegations[delegator] = make(map[types.Address]*types.Delegation)
	}
	
	delegation, exists := d.delegations[delegator][validator]
	if !exists {
		delegation = types.NewDelegation(delegator, validator, amount)
		d.delegations[delegator][validator] = delegation
	} else {
		delegation.Amount.Add(delegation.Amount, amount)
	}
	
	// Update validator voting power
	val.AddVotingPower(amount)
	
	return nil
}

// Undelegate initiates undelegation
func (d *DPoS) Undelegate(delegator, validator types.Address, amount *big.Int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Check delegation exists
	delegation, exists := d.delegations[delegator][validator]
	if !exists {
		return errors.New("delegation not found")
	}
	
	if delegation.Amount.Cmp(amount) < 0 {
		return errors.New("insufficient delegated amount")
	}
	
	// Update delegation
	delegation.Amount.Sub(delegation.Amount, amount)
	
	// Update validator voting power
	val := d.validators[validator]
	val.SubVotingPower(amount)
	
	// Remove delegation if amount is zero
	if delegation.Amount.Sign() == 0 {
		delete(d.delegations[delegator], validator)
	}
	
	return nil
}

// SelectValidators selects active validators for next epoch
func (d *DPoS) SelectValidators() []*types.Validator {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Get all validators
	validators := make([]*types.Validator, 0, len(d.validators))
	for _, val := range d.validators {
		if val.CanProduceBlocks() {
			validators = append(validators, val)
		}
	}
	
	// Sort by voting power (descending)
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].VotingPower.Cmp(validators[j].VotingPower) > 0
	})
	
	// Select top MaxValidators
	maxValidators := types.MaxValidators
	if len(validators) < maxValidators {
		maxValidators = len(validators)
	}
	
	d.activeValidators = validators[:maxValidators]
	return d.activeValidators
}

// GetBlockProducer returns the validator that should produce the next block
func (d *DPoS) GetBlockProducer(blockNumber uint64) (*types.Validator, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if len(d.activeValidators) == 0 {
		return nil, errors.New("no active validators")
	}
	
	// Round-robin selection based on block number
	index := blockNumber % uint64(len(d.activeValidators))
	return d.activeValidators[index], nil
}

// ValidateBlock validates a block according to DPoS rules
func (d *DPoS) ValidateBlock(block *core.Block) error {
	// Get expected validator
	expectedValidator, err := d.GetBlockProducer(block.Header.Number)
	if err != nil {
		return err
	}
	
	// Check if block producer is correct
	if block.Header.Validator != expectedValidator.Address {
		return errors.New("invalid block producer")
	}
	
	// Check if validator is active
	if !expectedValidator.IsActive() {
		return errors.New("validator is not active")
	}
	
	// Validate block timestamp
	if block.Header.Timestamp.Before(time.Now().Add(-time.Minute)) {
		return errors.New("block timestamp too old")
	}
	
	return nil
}

// UpdateEpoch updates the current epoch and rotates validators
func (d *DPoS) UpdateEpoch(blockNumber uint64) {
	if blockNumber%types.EpochLength == 0 {
		d.mu.Lock()
		d.currentEpoch++
		d.mu.Unlock()
		
		// Select new validator set
		d.SelectValidators()
	}
}

// GetValidator returns a validator by address
func (d *DPoS) GetValidator(addr types.Address) (*types.Validator, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	val, exists := d.validators[addr]
	if !exists {
		return nil, errors.New("validator not found")
	}
	
	return val, nil
}

// GetActiveValidators returns current active validators
func (d *DPoS) GetActiveValidators() []*types.Validator {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	result := make([]*types.Validator, len(d.activeValidators))
	copy(result, d.activeValidators)
	return result
}

// GetDelegation returns delegation information
func (d *DPoS) GetDelegation(delegator, validator types.Address) (*types.Delegation, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if d.delegations[delegator] == nil {
		return nil, errors.New("no delegations found")
	}
	
	delegation, exists := d.delegations[delegator][validator]
	if !exists {
		return nil, errors.New("delegation not found")
	}
	
	return delegation, nil
}

// GetTotalVotingPower returns total voting power of all validators
func (d *DPoS) GetTotalVotingPower() *big.Int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	total := big.NewInt(0)
	for _, val := range d.activeValidators {
		total.Add(total, val.VotingPower)
	}
	
	return total
}

// GetCurrentEpoch returns current epoch number
func (d *DPoS) GetCurrentEpoch() uint64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.currentEpoch
}
