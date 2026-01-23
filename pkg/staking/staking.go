package staking

import (
	"errors"
	"math/big"
	"time"

	"github.com/apex/pkg/consensus"
	"github.com/apex/pkg/types"
)

// StakingManager manages staking operations
type StakingManager struct {
	dpos              *consensus.DPoS
	unbondingQueue    map[types.Address][]*types.UnbondingDelegation
	minStakeAmount    *big.Int
	unbondingPeriod   uint64
}

// NewStakingManager creates a new staking manager
func NewStakingManager(dpos *consensus.DPoS) *StakingManager {
	return &StakingManager{
		dpos:            dpos,
		unbondingQueue:  make(map[types.Address][]*types.UnbondingDelegation),
		minStakeAmount:  types.ToWei(10.0), // Minimum 10 APX to stake
		unbondingPeriod: types.UnbondingPeriod,
	}
}

// Stake allows user to stake tokens to a validator
func (sm *StakingManager) Stake(
	delegator, validator types.Address,
	amount *big.Int,
) error {
	// Validate amount
	if amount.Cmp(sm.minStakeAmount) < 0 {
		return errors.New("stake amount below minimum")
	}
	
	// Delegate to validator
	return sm.dpos.Delegate(delegator, validator, amount)
}

// Unstake initiates unstaking process
func (sm *StakingManager) Unstake(
	delegator, validator types.Address,
	amount *big.Int,
	currentBlock uint64,
) error {
	// Initiate undelegation
	err := sm.dpos.Undelegate(delegator, validator, amount)
	if err != nil {
		return err
	}
	
	// Create unbonding delegation
	unbonding := &types.UnbondingDelegation{
		Delegator:       delegator,
		Validator:       validator,
		Amount:          new(big.Int).Set(amount),
		CompletionBlock: currentBlock + sm.unbondingPeriod,
		CreatedAt:       time.Now(),
	}
	
	// Add to unbonding queue
	sm.unbondingQueue[delegator] = append(sm.unbondingQueue[delegator], unbonding)
	
	return nil
}

// ProcessUnbonding processes completed unbonding delegations
func (sm *StakingManager) ProcessUnbonding(
	currentBlock uint64,
) ([]*types.UnbondingDelegation, error) {
	completed := make([]*types.UnbondingDelegation, 0)
	
	for delegator, unbondings := range sm.unbondingQueue {
		remaining := make([]*types.UnbondingDelegation, 0)
		
		for _, unbonding := range unbondings {
			if currentBlock >= unbonding.CompletionBlock {
				// Unbonding complete
				completed = append(completed, unbonding)
			} else {
				// Still unbonding
				remaining = append(remaining, unbonding)
			}
		}
		
		if len(remaining) > 0 {
			sm.unbondingQueue[delegator] = remaining
		} else {
			delete(sm.unbondingQueue, delegator)
		}
	}
	
	return completed, nil
}

// GetUnbondingDelegations returns unbonding delegations for a delegator
func (sm *StakingManager) GetUnbondingDelegations(
	delegator types.Address,
) []*types.UnbondingDelegation {
	return sm.unbondingQueue[delegator]
}

// Redelegate moves stake from one validator to another
func (sm *StakingManager) Redelegate(
	delegator, srcValidator, dstValidator types.Address,
	amount *big.Int,
) error {
	// Undelegate from source
	err := sm.dpos.Undelegate(delegator, srcValidator, amount)
	if err != nil {
		return err
	}
	
	// Delegate to destination
	return sm.dpos.Delegate(delegator, dstValidator, amount)
}

// GetStakingInfo returns staking information for a delegator
func (sm *StakingManager) GetStakingInfo(
	delegator, validator types.Address,
) (map[string]interface{}, error) {
	delegation, err := sm.dpos.GetDelegation(delegator, validator)
	if err != nil {
		return nil, err
	}
	
	unbondings := sm.GetUnbondingDelegations(delegator)
	
	return map[string]interface{}{
		"delegator":         delegator.Hex(),
		"validator":         validator.Hex(),
		"staked_amount":     types.FromWei(delegation.Amount),
		"rewards":           types.FromWei(delegation.Rewards),
		"unbonding_count":   len(unbondings),
		"created_at":        delegation.CreatedAt,
	}, nil
}

// CalculateStakingReturns calculates expected staking returns
func (sm *StakingManager) CalculateStakingReturns(
	amount *big.Int,
	days uint64,
	apy float64,
) *big.Int {
	// Daily return = (amount * APY) / 365
	dailyReturn := new(big.Float).Mul(
		new(big.Float).SetInt(amount),
		big.NewFloat(apy/365/100),
	)
	
	// Total return = daily return * days
	totalReturn := new(big.Float).Mul(dailyReturn, big.NewFloat(float64(days)))
	
	result, _ := totalReturn.Int(nil)
	return result
}