package consensus

import (
	"errors"
	"math/big"
	"time"

	"github.com/apex-blockchain/apex/pkg/types"
)

// ValidatorManager manages validator lifecycle
type ValidatorManager struct {
	dpos *DPoS
}

// NewValidatorManager creates a new validator manager
func NewValidatorManager(dpos *DPoS) *ValidatorManager {
	return &ValidatorManager{
		dpos: dpos,
	}
}

// CreateValidator creates a new validator
func (vm *ValidatorManager) CreateValidator(
	address types.Address,
	pubKey []byte,
	selfStake *big.Int,
	commission uint64,
) error {
	// Validate commission rate (0-100%)
	if commission > 10000 {
		return errors.New("commission rate must be <= 100%")
	}
	
	// Validate minimum stake
	minStake := types.ToWei(float64(types.MinStakeAmount))
	if selfStake.Cmp(minStake) < 0 {
		return errors.New("insufficient self-stake for validator")
	}
	
	// Create validator
	validator := types.NewValidator(address, pubKey, selfStake, commission)
	
	// Register with DPoS
	return vm.dpos.RegisterValidator(validator)
}

// EditValidator updates validator metadata
func (vm *ValidatorManager) EditValidator(
	address types.Address,
	commission *uint64,
) error {
	validator, err := vm.dpos.GetValidator(address)
	if err != nil {
		return err
	}
	
	// Update commission if provided
	if commission != nil {
		if *commission > 10000 {
			return errors.New("commission rate must be <= 100%")
		}
		validator.Commission = *commission
	}
	
	return nil
}

// JailValidator jails a validator for misbehavior
func (vm *ValidatorManager) JailValidator(address types.Address, reason string) error {
	validator, err := vm.dpos.GetValidator(address)
	if err != nil {
		return err
	}
	
	if validator.Jailed {
		return errors.New("validator already jailed")
	}
	
	validator.Jailed = true
	validator.JailTime = time.Now()
	validator.Status = types.ValidatorStatusJailed
	
	return nil
}

// UnjailValidator unjails a validator
func (vm *ValidatorManager) UnjailValidator(address types.Address) error {
	validator, err := vm.dpos.GetValidator(address)
	if err != nil {
		return err
	}
	
	if !validator.Jailed {
		return errors.New("validator is not jailed")
	}
	
	// Check if jail period has passed (e.g., 24 hours)
	jailDuration := 24 * time.Hour
	if time.Since(validator.JailTime) < jailDuration {
		return errors.New("jail period not completed")
	}
	
	validator.Jailed = false
	validator.Status = types.ValidatorStatusActive
	validator.MissedBlocks = 0
	
	return nil
}

// IncrementMissedBlocks increments missed block count
func (vm *ValidatorManager) IncrementMissedBlocks(address types.Address) error {
	validator, err := vm.dpos.GetValidator(address)
	if err != nil {
		return err
	}
	
	validator.MissedBlocks++
	
	// Auto-jail after threshold (e.g., 100 missed blocks)
	const missedBlockThreshold = 100
	if validator.MissedBlocks >= missedBlockThreshold && !validator.Jailed {
		return vm.JailValidator(address, "excessive missed blocks")
	}
	
	return nil
}

// IncrementProducedBlocks increments produced block count
func (vm *ValidatorManager) IncrementProducedBlocks(address types.Address) error {
	validator, err := vm.dpos.GetValidator(address)
	if err != nil {
		return err
	}
	
	validator.ProducedBlocks++
	validator.LastActiveEpoch = vm.dpos.GetCurrentEpoch()
	
	return nil
}

// CalculateUptime calculates validator uptime percentage
func (vm *ValidatorManager) CalculateUptime(address types.Address) (float64, error) {
	validator, err := vm.dpos.GetValidator(address)
	if err != nil {
		return 0, err
	}
	
	total := validator.ProducedBlocks + validator.MissedBlocks
	if total == 0 {
		return 100.0, nil
	}
	
	uptime := float64(validator.ProducedBlocks) / float64(total) * 100
	return uptime, nil
}

// GetValidatorStats returns validator statistics
func (vm *ValidatorManager) GetValidatorStats(address types.Address) (map[string]interface{}, error) {
	validator, err := vm.dpos.GetValidator(address)
	if err != nil {
		return nil, err
	}
	
	uptime, _ := vm.CalculateUptime(address)
	
	return map[string]interface{}{
		"address":         validator.Address.Hex(),
		"voting_power":    types.FromWei(validator.VotingPower),
		"self_stake":      types.FromWei(validator.SelfStake),
		"commission":      float64(validator.Commission) / 100,
		"status":          validator.Status,
		"jailed":          validator.Jailed,
		"missed_blocks":   validator.MissedBlocks,
		"produced_blocks": validator.ProducedBlocks,
		"uptime":          uptime,
		"total_rewards":   types.FromWei(validator.TotalRewards),
	}, nil
}

// IsActiveValidator checks if validator is in active set
func (vm *ValidatorManager) IsActiveValidator(address types.Address) bool {
	activeValidators := vm.dpos.GetActiveValidators()
	for _, val := range activeValidators {
		if val.Address == address {
			return true
		}
	}
	return false
}