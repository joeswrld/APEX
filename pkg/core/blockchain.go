package core

import (
	"crypto/ecdsa"
	"errors"
	"sync"
	"time"

	"github.com/apex/pkg/consensus"
	"github.com/apex/pkg/crypto"
	"github.com/apex/pkg/storage"
	"github.com/apex/pkg/types"
)

// Blockchain represents the main blockchain
type Blockchain struct {
	blocks       []*Block
	blocksByHash map[types.Hash]*Block
	stateDB      *storage.StateDB
	blockStore   *storage.BlockStore
	dpos         *consensus.DPoS
	rewardCalc   *consensus.RewardCalculator
	slasher      *consensus.Slasher
	executor     *Executor
	mu           sync.RWMutex
}

// NewBlockchain creates a new blockchain
func NewBlockchain(
	stateDB *storage.StateDB,
	blockStore *storage.BlockStore,
	dpos *consensus.DPoS,
) *Blockchain {
	bc := &Blockchain{
		blocks:       make([]*Block, 0),
		blocksByHash: make(map[types.Hash]*Block),
		stateDB:      stateDB,
		blockStore:   blockStore,
		dpos:         dpos,
		rewardCalc:   consensus.NewRewardCalculator(dpos),
		slasher:      consensus.NewSlasher(dpos),
	}
	
	bc.executor = NewExecutor(bc, stateDB)
	
	return bc
}

// InitGenesis initializes blockchain with genesis block
func (bc *Blockchain) InitGenesis(genesisValidators []*types.Validator, genesisAccounts []*types.Account) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	// Create genesis block
	genesis := NewBlock(0, types.Hash{}, types.Address{})
	genesis.Header.Timestamp = time.Now()
	
	// Initialize validators
	for _, validator := range genesisValidators {
		if err := bc.dpos.RegisterValidator(validator); err != nil {
			return err
		}
		if err := bc.stateDB.SetValidator(validator); err != nil {
			return err
		}
	}
	
	// Initialize accounts
	for _, account := range genesisAccounts {
		if err := bc.stateDB.SetAccount(account); err != nil {
			return err
		}
	}
	
	// Select initial validator set
	bc.dpos.SelectValidators()
	
	// Finalize genesis block
	stateRoot, _ := bc.stateDB.GetStateRoot()
	genesis.Finalize(stateRoot)
	
	// Store genesis block
	bc.blocks = append(bc.blocks, genesis)
	bc.blocksByHash[genesis.Hash] = genesis
	
	if err := bc.blockStore.PutBlock(genesis); err != nil {
		return err
	}
	
	return nil
}

// ProduceBlock produces a new block (called by validator)
func (bc *Blockchain) ProduceBlock(
	validatorKey *ecdsa.PrivateKey,
	transactions []*Transaction,
) (*Block, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	// Get validator address
	validatorAddr := crypto.PublicKeyToAddress(&validatorKey.PublicKey)
	
	// Check if validator should produce this block
	currentHeight := bc.GetHeight()
	expectedValidator, err := bc.dpos.GetBlockProducer(currentHeight + 1)
	if err != nil {
		return nil, err
	}
	
	if expectedValidator.Address != validatorAddr {
		return nil, errors.New("not your turn to produce block")
	}
	
	// Get previous block
	previousBlock := bc.GetLatestBlock()
	
	// Create new block
	block := NewBlock(currentHeight+1, previousBlock.Hash, validatorAddr)
	
	// Add transactions
	totalFees := types.ToWei(0)
	for _, tx := range transactions {
		if block.AddTransaction(tx) {
			// Calculate fees
			fee := tx.GasPrice
			totalFees.Add(totalFees, fee)
		}
	}
	
	// Execute transactions and update state
	if err := bc.executor.ExecuteBlock(block); err != nil {
		return nil, err
	}
	
	// Distribute block rewards
	if err := bc.rewardCalc.DistributeBlockReward(validatorAddr, block.Header.Number, totalFees); err != nil {
		return nil, err
	}
	
	// Get state root
	stateRoot, _ := bc.stateDB.GetStateRoot()
	
	// Finalize block
	block.Finalize(stateRoot)
	
	// Sign block
	signature, err := crypto.SignHash(block.Hash, validatorKey)
	if err != nil {
		return nil, err
	}
	block.Sign(signature)
	
	return block, nil
}

// AddBlock adds a validated block to the chain
func (bc *Blockchain) AddBlock(block *Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	// Validate block
	if err := bc.ValidateBlock(block); err != nil {
		return err
	}
	
	// Execute block
	if err := bc.executor.ExecuteBlock(block); err != nil {
		return err
	}
	
	// Store block
	bc.blocks = append(bc.blocks, block)
	bc.blocksByHash[block.Hash] = block
	
	if err := bc.blockStore.PutBlock(block); err != nil {
		return err
	}
	
	// Update epoch if needed
	bc.dpos.UpdateEpoch(block.Header.Number)
	
	return nil
}

// ValidateBlock validates a block
func (bc *Blockchain) ValidateBlock(block *Block) error {
	// Basic validation
	if err := block.Validate(); err != nil {
		return err
	}
	
	// Check block number
	expectedHeight := bc.GetHeight() + 1
	if block.Header.Number != expectedHeight {
		return errors.New("invalid block number")
	}
	
	// Check previous hash
	previousBlock := bc.GetLatestBlock()
	if block.Header.PreviousHash != previousBlock.Hash {
		return errors.New("invalid previous hash")
	}
	
	// Validate with DPoS
	if err := bc.dpos.ValidateBlock(block); err != nil {
		return err
	}
	
	return nil
}

// GetLatestBlock returns the latest block
func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	
	if len(bc.blocks) == 0 {
		return nil
	}
	
	return bc.blocks[len(bc.blocks)-1]
}

// GetHeight returns current blockchain height
func (bc *Blockchain) GetHeight() uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	
	if len(bc.blocks) == 0 {
		return 0
	}
	
	return bc.blocks[len(bc.blocks)-1].Header.Number
}

// GetBlockByHash retrieves a block by hash
func (bc *Blockchain) GetBlockByHash(hash types.Hash) (*Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	
	block, exists := bc.blocksByHash[hash]
	if !exists {
		return bc.blockStore.GetBlock(hash)
	}
	
	return block, nil
}

// GetBlockByNumber retrieves a block by number
func (bc *Blockchain) GetBlockByNumber(number uint64) (*Block, error) {
	return bc.blockStore.GetBlockByNumber(number)
}

// GetStateDB returns the state database
func (bc *Blockchain) GetStateDB() *storage.StateDB {
	return bc.stateDB
}