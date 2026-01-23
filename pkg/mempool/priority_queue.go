package mempool

import (
	"container/heap"
	"math/big"

	"github.com/apex/pkg/core"
)

// PriorityQueue implements a priority queue for transactions
type PriorityQueue struct {
	items txHeap
}

// NewPriorityQueue creates a new priority queue
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		items: make(txHeap, 0),
	}
	heap.Init(&pq.items)
	return pq
}

// Push adds a transaction to the queue
func (pq *PriorityQueue) Push(tx *core.Transaction) {
	heap.Push(&pq.items, tx)
}

// Pop removes and returns the highest priority transaction
func (pq *PriorityQueue) Pop() *core.Transaction {
	if len(pq.items) == 0 {
		return nil
	}
	return heap.Pop(&pq.items).(*core.Transaction)
}

// Top returns top N transactions without removing them
func (pq *PriorityQueue) Top(n int) []*core.Transaction {
	if n > len(pq.items) {
		n = len(pq.items)
	}
	
	result := make([]*core.Transaction, n)
	for i := 0; i < n; i++ {
		result[i] = pq.items[i]
	}
	
	return result
}

// Remove removes a specific transaction
func (pq *PriorityQueue) Remove(tx *core.Transaction) {
	for i, item := range pq.items {
		if item.Hash == tx.Hash {
			heap.Remove(&pq.items, i)
			return
		}
	}
}

// Len returns queue length
func (pq *PriorityQueue) Len() int {
	return len(pq.items)
}

// txHeap implements heap.Interface for transactions
type txHeap []*core.Transaction

func (h txHeap) Len() int {
	return len(h)
}

func (h txHeap) Less(i, j int) bool {
	// Higher gas price = higher priority
	return h[i].GasPrice.Cmp(h[j].GasPrice) > 0
}

func (h txHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *txHeap) Push(x interface{}) {
	*h = append(*h, x.(*core.Transaction))
}

func (h *txHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// CalculatePriority calculates transaction priority
func CalculatePriority(tx *core.Transaction) *big.Int {
	// Priority = gas price * gas limit
	priority := new(big.Int).Mul(tx.GasPrice, big.NewInt(int64(tx.GasLimit)))
	return priority
}