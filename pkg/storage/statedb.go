package storage

import (
	"encoding/json"
	"fmt"

	"github.com/apex/pkg/types"
)

// StateDB manages blockchain state
type StateDB struct {
	db *Database
}

// NewStateDB creates a new state database
func NewStateDB(db *Database) *StateDB {
	return &StateDB{db: db}
}

// GetAccount retrieves an account by address
func (s *StateDB) GetAccount(addr types.Address) (*types.Account, error) {
	key := accountKey(addr)
	data, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}
	
	var account types.Account
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, err
	}
	
	return &account, nil
}

// SetAccount stores an account
func (s *StateDB) SetAccount(account *types.Account) error {
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}
	
	key := accountKey(account.Address)
	return s.db.Put(key, data)
}

// DeleteAccount deletes an account
func (s *StateDB) DeleteAccount(addr types.Address) error {
	key := accountKey(addr)
	return s.db.Delete(key)
}

// GetValidator retrieves a validator by address
func (s *StateDB) GetValidator(addr types.Address) (*types.Validator, error) {
	key := validatorKey(addr)
	data, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}
	
	var validator types.Validator
	if err := json.Unmarshal(data, &validator); err != nil {
		return nil, err
	}
	
	return &validator, nil
}

// SetValidator stores a validator
func (s *StateDB) SetValidator(validator *types.Validator) error {
	data, err := json.Marshal(validator)
	if err != nil {
		return err
	}
	
	key := validatorKey(validator.Address)
	return s.db.Put(key, data)
}

// DeleteValidator deletes a validator
func (s *StateDB) DeleteValidator(addr types.Address) error {
	key := validatorKey(addr)
	return s.db.Delete(key)
}

// GetAllValidators retrieves all validators
func (s *StateDB) GetAllValidators() ([]*types.Validator, error) {
	validators := make([]*types.Validator, 0)
	
	prefix := []byte("validator:")
	iter := s.db.Iterator(prefix)
	defer iter.Close()
	
	for iter.Rewind(); iter.Valid(); iter.Next() {
		data, err := iter.Value()
		if err != nil {
			continue
		}
		
		var validator types.Validator
		if err := json.Unmarshal(data, &validator); err != nil {
			continue
		}
		
		validators = append(validators, &validator)
	}
	
	return validators, nil
}

// GetDelegation retrieves a delegation
func (s *StateDB) GetDelegation(delegator, validator types.Address) (*types.Delegation, error) {
	key := delegationKey(delegator, validator)
	data, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}
	
	var delegation types.Delegation
	if err := json.Unmarshal(data, &delegation); err != nil {
		return nil, err
	}
	
	return &delegation, nil
}

// SetDelegation stores a delegation
func (s *StateDB) SetDelegation(delegation *types.Delegation) error {
	data, err := json.Marshal(delegation)
	if err != nil {
		return err
	}
	
	key := delegationKey(delegation.Delegator, delegation.Validator)
	return s.db.Put(key, data)
}

// DeleteDelegation deletes a delegation
func (s *StateDB) DeleteDelegation(delegator, validator types.Address) error {
	key := delegationKey(delegator, validator)
	return s.db.Delete(key)
}

// GetStateRoot computes the state root hash
func (s *StateDB) GetStateRoot() (types.Hash, error) {
	// Simplified implementation - in production, use Merkle Patricia Tree
	var root types.Hash
	// Compute root from all state...
	return root, nil
}

// Commit commits state changes
func (s *StateDB) Commit() error {
	// In a full implementation, this would batch and commit all pending changes
	return nil
}

// Snapshot creates a state snapshot
func (s *StateDB) Snapshot() int {
	// Returns snapshot ID - not fully implemented
	return 0
}

// RevertToSnapshot reverts to a snapshot
func (s *StateDB) RevertToSnapshot(snapshot int) {
	// Revert logic - not fully implemented
}

// Key generation helpers
func accountKey(addr types.Address) []byte {
	return []byte(fmt.Sprintf("account:%s", addr.Hex()))
}

func validatorKey(addr types.Address) []byte {
	return []byte(fmt.Sprintf("validator:%s", addr.Hex()))
}

func delegationKey(delegator, validator types.Address) []byte {
	return []byte(fmt.Sprintf("delegation:%s:%s", delegator.Hex(), validator.Hex()))
}