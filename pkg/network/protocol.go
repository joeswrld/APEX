package network

import (
	"encoding/json"
	"io"

	"github.com/apex/pkg/core"
	"github.com/libp2p/go-libp2p/core/network"
	"go.uber.org/zap"
)

// MessageType represents different message types
type MessageType uint8

const (
	MsgTypeBlock MessageType = iota
	MsgTypeTransaction
	MsgTypeGetBlocks
	MsgTypeGetBlockHeaders
	MsgTypeBlockHeaders
	MsgTypeGetState
	MsgTypeState
)

// Message represents a network message
type Message struct {
	Type MessageType `json:"type"`
	Data []byte      `json:"data"`
}

// Protocol handles network protocol messages
type Protocol struct {
	blockchain *core.Blockchain
	network    *P2PNetwork
	logger     *zap.Logger
}

// NewProtocol creates a new protocol handler
func NewProtocol(blockchain *core.Blockchain, network *P2PNetwork, logger *zap.Logger) *Protocol {
	return &Protocol{
		blockchain: blockchain,
		network:    network,
		logger:     logger,
	}
}

// HandleStream handles incoming protocol streams
func (p *Protocol) HandleStream(stream network.Stream) {
	defer stream.Close()
	
	// Read message
	var msg Message
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&msg); err != nil {
		if err != io.EOF {
			p.logger.Error("Failed to decode message", zap.Error(err))
		}
		return
	}
	
	// Route message to handler
	switch msg.Type {
	case MsgTypeBlock:
		p.handleBlock(msg.Data, stream)
	case MsgTypeTransaction:
		p.handleTransaction(msg.Data, stream)
	case MsgTypeGetBlocks:
		p.handleGetBlocks(msg.Data, stream)
	case MsgTypeGetBlockHeaders:
		p.handleGetBlockHeaders(msg.Data, stream)
	default:
		p.logger.Warn("Unknown message type", zap.Uint8("type", uint8(msg.Type)))
	}
}

// handleBlock handles incoming block messages
func (p *Protocol) handleBlock(data []byte, stream network.Stream) {
	var block core.Block
	if err := json.Unmarshal(data, &block); err != nil {
		p.logger.Error("Failed to unmarshal block", zap.Error(err))
		return
	}
	
	p.logger.Info("Received block", zap.Uint64("number", block.Header.Number))
	
	// Add block to blockchain
	if err := p.blockchain.AddBlock(&block); err != nil {
		p.logger.Error("Failed to add block", zap.Error(err))
		return
	}
	
	// Broadcast to other peers
	p.broadcastBlock(&block)
}

// handleTransaction handles incoming transaction messages
func (p *Protocol) handleTransaction(data []byte, stream network.Stream) {
	var tx core.Transaction
	if err := json.Unmarshal(data, &tx); err != nil {
		p.logger.Error("Failed to unmarshal transaction", zap.Error(err))
		return
	}
	
	p.logger.Debug("Received transaction", zap.String("hash", tx.Hash.Hex()))
	
	// Process transaction (add to mempool, etc.)
	// This would integrate with mempool in a full implementation
}

// handleGetBlocks handles block requests
func (p *Protocol) handleGetBlocks(data []byte, stream network.Stream) {
	var req struct {
		Start uint64 `json:"start"`
		End   uint64 `json:"end"`
	}
	
	if err := json.Unmarshal(data, &req); err != nil {
		p.logger.Error("Failed to unmarshal get blocks request", zap.Error(err))
		return
	}
	
	p.logger.Debug("Received get blocks request",
		zap.Uint64("start", req.Start),
		zap.Uint64("end", req.End),
	)
	
	// Get blocks from blockchain
	blocks := make([]*core.Block, 0)
	for i := req.Start; i <= req.End && i <= req.Start+100; i++ {
		block, err := p.blockchain.GetBlockByNumber(i)
		if err != nil {
			break
		}
		blocks = append(blocks, block)
	}
	
	// Send blocks
	encoder := json.NewEncoder(stream)
	for _, block := range blocks {
		msg := Message{
			Type: MsgTypeBlock,
			Data: mustMarshal(block),
		}
		if err := encoder.Encode(msg); err != nil {
			p.logger.Error("Failed to send block", zap.Error(err))
			break
		}
	}
}

// handleGetBlockHeaders handles block header requests
func (p *Protocol) handleGetBlockHeaders(data []byte, stream network.Stream) {
	// Similar to handleGetBlocks but only sends headers
	p.logger.Debug("Received get block headers request")
}

// broadcastBlock broadcasts a block to all peers
func (p *Protocol) broadcastBlock(block *core.Block) {
	msg := Message{
		Type: MsgTypeBlock,
		Data: mustMarshal(block),
	}
	
	msgData := mustMarshal(msg)
	if err := p.network.Broadcast("blocks", msgData); err != nil {
		p.logger.Error("Failed to broadcast block", zap.Error(err))
	}
}

// BroadcastTransaction broadcasts a transaction to all peers
func (p *Protocol) BroadcastTransaction(tx *core.Transaction) {
	msg := Message{
		Type: MsgTypeTransaction,
		Data: mustMarshal(tx),
	}
	
	msgData := mustMarshal(msg)
	if err := p.network.Broadcast("transactions", msgData); err != nil {
		p.logger.Error("Failed to broadcast transaction", zap.Error(err))
	}
}

// mustMarshal marshals data or panics
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
