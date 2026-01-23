package mempool

import (
	"errors"
	"sync"

	"github.com/apex/pkg/core"
	"github.com/apex/pkg/types"
)

// Mempool manages pending transactions
type Mempool struct {
	transactions map[types.Hash]*core.Transaction
	queue        *PriorityQueue
	maxSize      int
	mu           sync.RWMutex
}

// NewMempool creates a new mempool
func NewMempool(maxSize int) *Mempool {
	return &Mempool{
		transactions: make(map[types.Hash]*core.Transaction),
		queue:        NewPriorityQueue(),
		maxSize:      maxSize,
	}
}

// AddTransaction adds a transaction to mempool
func (m *Mempool) AddTransaction(tx *core.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if already exists
	if _, exists := m.transactions[tx.Hash]; exists {
		return errors.New("transaction already in mempool")
	}
	
	// Validate transaction
	if err := tx.Validate(); err != nil {
		return err
	}
	
	// Check mempool size
	if len(m.transactions) >= m.maxSize {
		// Remove lowest priority transaction
		if lowest := m.queue.Pop(); lowest != nil {
			delete(m.transactions, lowest.Hash)
		}
	}
	
	// Add to mempool
	m.transactions[tx.Hash] = tx
	m.queue.Push(tx)
	
	return nil
}

// RemoveTransaction removes a transaction from mempool
func (m *Mempool) RemoveTransaction(hash types.Hash) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if tx, exists := m.transactions[hash]; exists {
		delete(m.transactions, hash)
		m.queue.Remove(tx)
	}
}

// GetTransaction retrieves a transaction by hash
func (m *Mempool) GetTransaction(hash types.Hash) (*core.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	tx, exists := m.transactions[hash]
	if !exists {
		return nil, errors.New("transaction not found in mempool")
	}
	
	return tx, nil
}

// GetTransactions returns top N transactions by priority
func (m *Mempool) GetTransactions(limit int) []*core.Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return m.queue.Top(limit)
}

// Size returns current mempool size
func (m *Mempool) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return len(m.transactions)
}

// Clear clears the mempool
func (m *Mempool) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.transactions = make(map[types.Hash]*core.Transaction)
	m.queue = NewPriorityQueue()
}

// Has checks if transaction exists in mempool
func (m *Mempool) Has(hash types.Hash) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	_, exists := m.transactions[hash]
	return exists
}

// RemoveTransactions removes multiple transactions
func (m *Mempool) RemoveTransactions(hashes []types.Hash) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, hash := range hashes {
		if tx, exists := m.transactions[hash]; exists {
			delete(m.transactions, hash)
			m.queue.Remove(tx)
		}
	}
}

// GetPendingCount returns count of pending transactions
func (m *Mempool) GetPendingCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return len(m.transactions)
}