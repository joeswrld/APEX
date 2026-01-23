package storage

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/apex/pkg/core"
	"github.com/apex/pkg/types"
)

// BlockStore manages block storage
type BlockStore struct {
	db *Database
}

// NewBlockStore creates a new block store
func NewBlockStore(db *Database) *BlockStore {
	return &BlockStore{db: db}
}

// PutBlock stores a block
func (bs *BlockStore) PutBlock(block *core.Block) error {
	// Serialize block
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}
	
	// Store by hash
	hashKey := blockHashKey(block.Hash)
	if err := bs.db.Put(hashKey, data); err != nil {
		return err
	}
	
	// Store by number (for quick lookups)
	numberKey := blockNumberKey(block.Header.Number)
	if err := bs.db.Put(numberKey, block.Hash[:]); err != nil {
		return err
	}
	
	// Update latest block number
	return bs.updateLatestBlockNumber(block.Header.Number)
}

// GetBlock retrieves a block by hash
func (bs *BlockStore) GetBlock(hash types.Hash) (*core.Block, error) {
	key := blockHashKey(hash)
	data, err := bs.db.Get(key)
	if err != nil {
		return nil, err
	}
	
	var block core.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, err
	}
	
	return &block, nil
}

// GetBlockByNumber retrieves a block by number
func (bs *BlockStore) GetBlockByNumber(number uint64) (*core.Block, error) {
	// Get hash from number index
	numberKey := blockNumberKey(number)
	hashData, err := bs.db.Get(numberKey)
	if err != nil {
		return nil, err
	}
	
	var hash types.Hash
	copy(hash[:], hashData)
	
	// Get block by hash
	return bs.GetBlock(hash)
}

// GetLatestBlock retrieves the latest block
func (bs *BlockStore) GetLatestBlock() (*core.Block, error) {
	number, err := bs.GetLatestBlockNumber()
	if err != nil {
		return nil, err
	}
	
	return bs.GetBlockByNumber(number)
}

// GetLatestBlockNumber retrieves the latest block number
func (bs *BlockStore) GetLatestBlockNumber() (uint64, error) {
	key := []byte("latest_block_number")
	data, err := bs.db.Get(key)
	if err != nil {
		return 0, err
	}
	
	return binary.BigEndian.Uint64(data), nil
}

// updateLatestBlockNumber updates the latest block number
func (bs *BlockStore) updateLatestBlockNumber(number uint64) error {
	key := []byte("latest_block_number")
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, number)
	return bs.db.Put(key, data)
}

// PutTransaction stores a transaction
func (bs *BlockStore) PutTransaction(tx *core.Transaction) error {
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	
	key := txKey(tx.Hash)
	return bs.db.Put(key, data)
}

// GetTransaction retrieves a transaction by hash
func (bs *BlockStore) GetTransaction(hash types.Hash) (*core.Transaction, error) {
	key := txKey(hash)
	data, err := bs.db.Get(key)
	if err != nil {
		return nil, err
	}
	
	var tx core.Transaction
	if err := json.Unmarshal(data, &tx); err != nil {
		return nil, err
	}
	
	return &tx, nil
}

// PutReceipt stores a transaction receipt
func (bs *BlockStore) PutReceipt(receipt *core.TxReceipt) error {
	data, err := json.Marshal(receipt)
	if err != nil {
		return err
	}
	
	key := receiptKey(receipt.TxHash)
	return bs.db.Put(key, data)
}

// GetReceipt retrieves a transaction receipt
func (bs *BlockStore) GetReceipt(txHash types.Hash) (*core.TxReceipt, error) {
	key := receiptKey(txHash)
	data, err := bs.db.Get(key)
	if err != nil {
		return nil, err
	}
	
	var receipt core.TxReceipt
	if err := json.Unmarshal(data, &receipt); err != nil {
		return nil, err
	}
	
	return &receipt, nil
}

// GetBlockRange retrieves blocks in a range
func (bs *BlockStore) GetBlockRange(start, end uint64) ([]*core.Block, error) {
	blocks := make([]*core.Block, 0, end-start+1)
	
	for i := start; i <= end; i++ {
		block, err := bs.GetBlockByNumber(i)
		if err != nil {
			continue
		}
		blocks = append(blocks, block)
	}
	
	return blocks, nil
}

// HasBlock checks if a block exists
func (bs *BlockStore) HasBlock(hash types.Hash) bool {
	key := blockHashKey(hash)
	has, _ := bs.db.Has(key)
	return has
}

// Key generation helpers
func blockHashKey(hash types.Hash) []byte {
	return []byte(fmt.Sprintf("block:hash:%s", hash.Hex()))
}

func blockNumberKey(number uint64) []byte {
	return []byte(fmt.Sprintf("block:number:%d", number))
}

func txKey(hash types.Hash) []byte {
	return []byte(fmt.Sprintf("tx:%s", hash.Hex()))
}

func receiptKey(hash types.Hash) []byte {
	return []byte(fmt.Sprintf("receipt:%s", hash.Hex()))
}