package consensus

import (
	"errors"
	"math/big"

	"github.com/apex/pkg/types"
)

// SlashingReason represents reason for slashing
type SlashingReason int

const (
	SlashingReasonDoubleSign SlashingReason = iota
	SlashingReasonDowntime
	SlashingReasonInvalidBlock
)

// SlashingEvent represents a slashing event
type SlashingEvent struct {
	Validator types.Address
	Reason    SlashingReason
	Amount    *big.Int
	BlockNum  uint64
}

// Slasher handles validator slashing
type Slasher struct {
	dpos   *DPoS
	events []SlashingEvent
}

// NewSlasher creates a new slasher
func NewSlasher(dpos *DPoS) *Slasher {
	return &Slasher{
		dpos:   dpos,
		events: make([]SlashingEvent, 0),
	}
}

// SlashValidator slashes a validator's stake
func (s *Slasher) SlashValidator(
	address types.Address,
	reason SlashingReason,
	blockNum uint64,
) error {
	validator, err := s.dpos.GetValidator(address)
	if err != nil {
		return err
	}
	
	// Calculate slash amount based on reason
	var slashPercentage float64
	switch reason {
	case SlashingReasonDoubleSign:
		slashPercentage = 0.05 // 5% slash for double signing
	case SlashingReasonDowntime:
		slashPercentage = 0.01 // 1% slash for downtime
	case SlashingReasonInvalidBlock:
		slashPercentage = 0.03 // 3% slash for invalid block
	default:
		return errors.New("unknown slashing reason")
	}
	
	// Calculate slash amount from voting power
	slashAmount := new(big.Int).Mul(
		validator.VotingPower,
		big.NewInt(int64(slashPercentage*10000)),
	)
	slashAmount.Div(slashAmount, big.NewInt(10000))
	
	// Apply slash
	validator.SubVotingPower(slashAmount)
	
	// Also slash self-stake proportionally
	selfStakeSlash := new(big.Int).Mul(
		validator.SelfStake,
		big.NewInt(int64(slashPercentage*10000)),
	)
	selfStakeSlash.Div(selfStakeSlash, big.NewInt(10000))
	validator.SelfStake.Sub(validator.SelfStake, selfStakeSlash)
	
	// Record slashing event
	event := SlashingEvent{
		Validator: address,
		Reason:    reason,
		Amount:    slashAmount,
		BlockNum:  blockNum,
	}
	s.events = append(s.events, event)
	
	// Jail validator for serious offenses
	if reason == SlashingReasonDoubleSign || reason == SlashingReasonInvalidBlock {
		validator.Jailed = true
		validator.Status = types.ValidatorStatusJailed
	}
	
	return nil
}

// DetectDoubleSign detects if validator signed multiple blocks at same height
func (s *Slasher) DetectDoubleSign(
	address types.Address,
	block1Height, block2Height uint64,
) error {
	if block1Height == block2Height {
		return s.SlashValidator(address, SlashingReasonDoubleSign, block1Height)
	}
	return nil
}

// CheckDowntime checks validator downtime and applies penalties
func (s *Slasher) CheckDowntime(address types.Address, missedBlocks uint64) error {
	// Slash if missed too many blocks in a row
	const downtimeThreshold = 50
	if missedBlocks >= downtimeThreshold {
		validator, err := s.dpos.GetValidator(address)
		if err != nil {
			return err
		}
		
		return s.SlashValidator(address, SlashingReasonDowntime, validator.LastActiveEpoch)
	}
	return nil
}

// GetSlashingEvents returns all slashing events
func (s *Slasher) GetSlashingEvents() []SlashingEvent {
	return s.events
}

// GetValidatorSlashingHistory returns slashing history for a validator
func (s *Slasher) GetValidatorSlashingHistory(address types.Address) []SlashingEvent {
	var history []SlashingEvent
	for _, event := range s.events {
		if event.Validator == address {
			history = append(history, event)
		}
	}
	return history
}
