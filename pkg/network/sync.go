package network

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/apex/pkg/core"
	"go.uber.org/zap"
)

// Syncer handles blockchain synchronization
type Syncer struct {
	blockchain *core.Blockchain
	network    *P2PNetwork
	protocol   *Protocol
	logger     *zap.Logger
	syncing    bool
	mu         sync.RWMutex
}

// NewSyncer creates a new syncer
func NewSyncer(blockchain *core.Blockchain, network *P2PNetwork, protocol *Protocol, logger *zap.Logger) *Syncer {
	return &Syncer{
		blockchain: blockchain,
		network:    network,
		protocol:   protocol,
		logger:     logger,
		syncing:    false,
	}
}

// Start starts the synchronization process
func (s *Syncer) Start() {
	go s.syncLoop()
}

// syncLoop periodically checks if sync is needed
func (s *Syncer) syncLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		if s.IsSyncing() {
			continue
		}
		
		if s.needsSync() {
			s.sync()
		}
	}
}

// needsSync checks if blockchain needs synchronization
func (s *Syncer) needsSync() bool {
	// Get current height
	currentHeight := s.blockchain.GetHeight()
	
	// Get peer heights
	peerHeight := s.getPeerHeight()
	
	// Need sync if peers are ahead by more than 10 blocks
	return peerHeight > currentHeight+10
}

// sync synchronizes the blockchain
func (s *Syncer) sync() {
	s.mu.Lock()
	if s.syncing {
		s.mu.Unlock()
		return
	}
	s.syncing = true
	s.mu.Unlock()
	
	defer func() {
		s.mu.Lock()
		s.syncing = false
		s.mu.Unlock()
	}()
	
	s.logger.Info("Starting blockchain synchronization")
	
	currentHeight := s.blockchain.GetHeight()
	targetHeight := s.getPeerHeight()
	
	// Sync in batches
	batchSize := uint64(100)
	for start := currentHeight + 1; start <= targetHeight; start += batchSize {
		end := start + batchSize - 1
		if end > targetHeight {
			end = targetHeight
		}
		
		s.logger.Info("Syncing blocks",
			zap.Uint64("from", start),
			zap.Uint64("to", end),
		)
		
		if err := s.syncRange(start, end); err != nil {
			s.logger.Error("Failed to sync range", zap.Error(err))
			return
		}
	}
	
	s.logger.Info("Synchronization completed",
		zap.Uint64("height", s.blockchain.GetHeight()),
	)
}

// syncRange syncs a range of blocks
func (s *Syncer) syncRange(start, end uint64) error {
	// Request blocks from peers
	req := struct {
		Start uint64 `json:"start"`
		End   uint64 `json:"end"`
	}{
		Start: start,
		End:   end,
	}
	
	reqData, _ := json.Marshal(req)
	msg := Message{
		Type: MsgTypeGetBlocks,
		Data: reqData,
	}
	
	msgData, _ := json.Marshal(msg)
	
	// Broadcast request
	if err := s.network.Broadcast("sync", msgData); err != nil {
		return err
	}
	
	// Wait for blocks (simplified - in production use proper request/response)
	time.Sleep(2 * time.Second)
	
	return nil
}

// getPeerHeight returns the highest block height from peers
func (s *Syncer) getPeerHeight() uint64 {
	// In a full implementation, this would query peers
	// For now, return current height
	return s.blockchain.GetHeight()
}

// IsSyncing returns true if currently syncing
func (s *Syncer) IsSyncing() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.syncing
}

// GetSyncProgress returns sync progress percentage
func (s *Syncer) GetSyncProgress() float64 {
	if !s.IsSyncing() {
		return 100.0
	}
	
	currentHeight := s.blockchain.GetHeight()
	targetHeight := s.getPeerHeight()
	
	if targetHeight == 0 {
		return 100.0
	}
	
	return float64(currentHeight) / float64(targetHeight) * 100.0
}

// FastSync performs fast synchronization using state snapshots
func (s *Syncer) FastSync() error {
	s.logger.Info("Fast sync not implemented yet")
	return nil
}

// SyncStatus returns current sync status
func (s *Syncer) SyncStatus() map[string]interface{} {
	return map[string]interface{}{
		"syncing":      s.IsSyncing(),
		"progress":     s.GetSyncProgress(),
		"current_height": s.blockchain.GetHeight(),
		"target_height":  s.getPeerHeight(),
		"peer_count":   s.network.GetPeerCount(),
	}
}
