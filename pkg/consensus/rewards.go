package consensus

import (
	"errors"
	"math/big"

	"github.com/apex/pkg/types"
)

// RewardCalculator calculates and distributes rewards
type RewardCalculator struct {
	dpos *DPoS
}

// NewRewardCalculator creates a new reward calculator
func NewRewardCalculator(dpos *DPoS) *RewardCalculator {
	return &RewardCalculator{
		dpos: dpos,
	}
}

// CalculateBlockReward calculates block production reward
func (rc *RewardCalculator) CalculateBlockReward(blockNumber uint64) *big.Int {
	// Base reward: 2 APX per block
	baseReward := types.ToWei(2.0)
	
	// Halving every 4 years (approximate)
	blocksPerYear := uint64(365 * 24 * 60 * 60 / types.BlockTime)
	halvingPeriod := blocksPerYear * 4
	halvings := blockNumber / halvingPeriod
	
	// Apply halving
	reward := new(big.Int).Set(baseReward)
	for i := uint64(0); i < halvings; i++ {
		reward.Div(reward, big.NewInt(2))
	}
	
	// Minimum reward: 0.1 APX
	minReward := types.ToWei(0.1)
	if reward.Cmp(minReward) < 0 {
		reward = minReward
	}
	
	return reward
}

// DistributeBlockReward distributes block reward to validator and delegators
func (rc *RewardCalculator) DistributeBlockReward(
	validatorAddr types.Address,
	blockNumber uint64,
	txFees *big.Int,
) error {
	validator, err := rc.dpos.GetValidator(validatorAddr)
	if err != nil {
		return err
	}
	
	// Calculate total reward (block reward + tx fees)
	blockReward := rc.CalculateBlockReward(blockNumber)
	totalReward := new(big.Int).Add(blockReward, txFees)
	
	// Calculate validator commission
	commission := validator.CalculateCommission(totalReward)
	
	// Remaining reward for delegators
	delegatorReward := new(big.Int).Sub(totalReward, commission)
	
	// Add commission to validator
	validator.TotalRewards.Add(validator.TotalRewards, commission)
	
	// Distribute to delegators proportionally
	if validator.VotingPower.Sign() > 0 {
		rc.distributeToDelegate(validatorAddr, delegatorReward)
	}
	
	return nil
}

// distributeToDelegate distributes rewards to delegators
func (rc *RewardCalculator) distributeToDelegate(
	validatorAddr types.Address,
	totalReward *big.Int,
) {
	validator, _ := rc.dpos.GetValidator(validatorAddr)
	
	// Get all delegations for this validator
	rc.dpos.mu.RLock()
	defer rc.dpos.mu.RUnlock()
	
	for delegatorAddr, delegations := range rc.dpos.delegations {
		if delegation, exists := delegations[validatorAddr]; exists {
			// Calculate delegator's share
			share := new(big.Int).Mul(totalReward, delegation.Amount)
			share.Div(share, validator.VotingPower)
			
			// Add to delegation rewards
			delegation.AddRewards(share)
			
			// Update account (if integrated with state)
			_ = delegatorAddr // Use delegatorAddr to update account state
		}
	}
}

// CalculateAPY calculates Annual Percentage Yield for staking
func (rc *RewardCalculator) CalculateAPY() float64 {
	totalStaked := rc.dpos.GetTotalVotingPower()
	if totalStaked.Sign() == 0 {
		return 0
	}
	
	// Approximate blocks per year
	blocksPerYear := uint64(365 * 24 * 60 * 60 / types.BlockTime)
	
	// Average reward per block
	avgReward := rc.CalculateBlockReward(0) // Use block 0 as reference
	
	// Annual rewards
	annualRewards := new(big.Int).Mul(avgReward, big.NewInt(int64(blocksPerYear)))
	
	// APY = (annual rewards / total staked) * 100
	apy := new(big.Float).Quo(
		new(big.Float).SetInt(annualRewards),
		new(big.Float).SetInt(totalStaked),
	)
	
	apyFloat, _ := apy.Float64()
	return apyFloat * 100
}

// ClaimRewards allows delegator to claim accumulated rewards
func (rc *RewardCalculator) ClaimRewards(
	delegatorAddr, validatorAddr types.Address,
) (*big.Int, error) {
	delegation, err := rc.dpos.GetDelegation(delegatorAddr, validatorAddr)
	if err != nil {
		return nil, err
	}
	
	if delegation.Rewards.Sign() == 0 {
		return nil, errors.New("no rewards to claim")
	}
	
	// Get rewards amount
	rewards := new(big.Int).Set(delegation.Rewards)
	
	// Reset delegation rewards
	delegation.Rewards = big.NewInt(0)
	
	return rewards, nil
}

// GetDelegatorRewards returns accumulated rewards for a delegator
func (rc *RewardCalculator) GetDelegatorRewards(
	delegatorAddr, validatorAddr types.Address,
) (*big.Int, error) {
	delegation, err := rc.dpos.GetDelegation(delegatorAddr, validatorAddr)
	if err != nil {
		return nil, err
	}
	
	return delegation.Rewards, nil
}
