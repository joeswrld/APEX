package core

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/apex/pkg/storage"
	"github.com/apex/pkg/types"
)

// Executor executes transactions and updates state
type Executor struct {
	blockchain *Blockchain
	stateDB    *storage.StateDB
}

// NewExecutor creates a new executor
func NewExecutor(blockchain *Blockchain, stateDB *storage.StateDB) *Executor {
	return &Executor{
		blockchain: blockchain,
		stateDB:    stateDB,
	}
}

// ExecuteBlock executes all transactions in a block
func (e *Executor) ExecuteBlock(block *Block) error {
	for _, tx := range block.Transactions {
		if err := e.ExecuteTransaction(tx); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteTransaction executes a single transaction
func (e *Executor) ExecuteTransaction(tx *Transaction) error {
	switch tx.Type {
	case TxTypeTransfer:
		return e.executeTransfer(tx)
	case TxTypeStake:
		return e.executeStake(tx)
	case TxTypeUnstake:
		return e.executeUnstake(tx)
	case TxTypeDelegate:
		return e.executeDelegate(tx)
	case TxTypeUndelegate:
		return e.executeUndelegate(tx)
	case TxTypeCreateValidator:
		return e.executeCreateValidator(tx)
	default:
		return errors.New("unknown transaction type")
	}
}

// executeTransfer executes a transfer transaction
func (e *Executor) executeTransfer(tx *Transaction) error {
	// Get sender account
	sender, err := e.stateDB.GetAccount(tx.From)
	if err != nil {
		return err
	}
	
	// Check nonce
	if sender.Nonce != tx.Nonce {
		return errors.New("invalid nonce")
	}
	
	// Check sufficient balance (value + gas)
	cost := tx.GetCost()
	if sender.Balance.Cmp(cost) < 0 {
		return errors.New("insufficient balance")
	}
	
	// Deduct from sender
	sender.SubBalance(cost)
	sender.Nonce++
	
	// Get recipient account
	recipient, err := e.stateDB.GetAccount(tx.To)
	if err != nil {
		recipient = types.NewAccount(tx.To)
	}
	
	// Add to recipient
	recipient.AddBalance(tx.Value)
	
	// Save accounts
	if err := e.stateDB.SetAccount(sender); err != nil {
		return err
	}
	if err := e.stateDB.SetAccount(recipient); err != nil {
		return err
	}
	
	return nil
}

// executeStake executes a stake transaction
func (e *Executor) executeStake(tx *Transaction) error {
	var data StakeData
	if err := json.Unmarshal(tx.Data, &data); err != nil {
		return err
	}
	
	// Get account
	account, err := e.stateDB.GetAccount(tx.From)
	if err != nil {
		return err
	}
	
	// Check balance
	if account.Balance.Cmp(data.Amount) < 0 {
		return errors.New("insufficient balance for staking")
	}
	
	// Transfer to staked
	account.SubBalance(data.Amount)
	account.AddStake(data.Amount)
	account.Nonce++
	
	// Save account
	return e.stateDB.SetAccount(account)
}

// executeUnstake executes an unstake transaction
func (e *Executor) executeUnstake(tx *Transaction) error {
	var data StakeData
	if err := json.Unmarshal(tx.Data, &data); err != nil {
		return err
	}
	
	// Get account
	account, err := e.stateDB.GetAccount(tx.From)
	if err != nil {
		return err
	}
	
	// Check staked amount
	if account.Staked.Cmp(data.Amount) < 0 {
		return errors.New("insufficient staked amount")
	}
	
	// Move to locked (unbonding)
	account.SubStake(data.Amount)
	account.Locked.Add(account.Locked, data.Amount)
	account.Nonce++
	
	// Save account
	return e.stateDB.SetAccount(account)
}

// executeDelegate executes a delegation transaction
func (e *Executor) executeDelegate(tx *Transaction) error {
	var data StakeData
	if err := json.Unmarshal(tx.Data, &data); err != nil {
		return err
	}
	
	// Get delegator account
	account, err := e.stateDB.GetAccount(tx.From)
	if err != nil {
		return err
	}
	
	// Check balance
	if account.Balance.Cmp(data.Amount) < 0 {
		return errors.New("insufficient balance for delegation")
	}
	
	// Get validator
	validator, err := e.stateDB.GetValidator(data.Validator)
	if err != nil {
		return errors.New("validator not found")
	}
	
	// Create or update delegation
	delegation, err := e.stateDB.GetDelegation(tx.From, data.Validator)
	if err != nil {
		delegation = types.NewDelegation(tx.From, data.Validator, data.Amount)
	} else {
		delegation.Amount.Add(delegation.Amount, data.Amount)
	}
	
	// Update account
	account.SubBalance(data.Amount)
	account.AddStake(data.Amount)
	account.Nonce++
	
	// Update validator voting power
	validator.AddVotingPower(data.Amount)
	
	// Save state
	if err := e.stateDB.SetAccount(account); err != nil {
		return err
	}
	if err := e.stateDB.SetValidator(validator); err != nil {
		return err
	}
	if err := e.stateDB.SetDelegation(delegation); err != nil {
		return err
	}
	
	return nil
}

// executeUndelegate executes an undelegation transaction
func (e *Executor) executeUndelegate(tx *Transaction) error {
	var data StakeData
	if err := json.Unmarshal(tx.Data, &data); err != nil {
		return err
	}
	
	// Get delegation
	delegation, err := e.stateDB.GetDelegation(tx.From, data.Validator)
	if err != nil {
		return errors.New("delegation not found")
	}
	
	// Check amount
	if delegation.Amount.Cmp(data.Amount) < 0 {
		return errors.New("insufficient delegated amount")
	}
	
	// Get validator
	validator, err := e.stateDB.GetValidator(data.Validator)
	if err != nil {
		return err
	}
	
	// Get account
	account, err := e.stateDB.GetAccount(tx.From)
	if err != nil {
		return err
	}
	
	// Update delegation
	delegation.Amount.Sub(delegation.Amount, data.Amount)
	
	// Update validator voting power
	validator.SubVotingPower(data.Amount)
	
	// Move to unbonding
	account.SubStake(data.Amount)
	account.Locked.Add(account.Locked, data.Amount)
	account.Nonce++
	
	// Save state
	if err := e.stateDB.SetAccount(account); err != nil {
		return err
	}
	if err := e.stateDB.SetValidator(validator); err != nil {
		return err
	}
	
	if delegation.Amount.Sign() == 0 {
		return e.stateDB.DeleteDelegation(tx.From, data.Validator)
	}
	
	return e.stateDB.SetDelegation(delegation)
}

// executeCreateValidator executes a create validator transaction
func (e *Executor) executeCreateValidator(tx *Transaction) error {
	var data CreateValidatorData
	if err := json.Unmarshal(tx.Data, &data); err != nil {
		return err
	}
	
	// Get account
	account, err := e.stateDB.GetAccount(tx.From)
	if err != nil {
		return err
	}
	
	// Check balance
	if account.Balance.Cmp(data.SelfStake) < 0 {
		return errors.New("insufficient balance for self-stake")
	}
	
	// Create validator
	validator := types.NewValidator(tx.From, data.PublicKey, data.SelfStake, data.Commission)
	
	// Update account
	account.SubBalance(data.SelfStake)
	account.AddStake(data.SelfStake)
	account.Nonce++
	
	// Save state
	if err := e.stateDB.SetAccount(account); err != nil {
		return err
	}
	
	return e.stateDB.SetValidator(validator)
}