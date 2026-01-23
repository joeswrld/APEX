package jsonrpc

import (
	"errors"
	"math/big"

	"github.com/apex/pkg/core"
	"github.com/apex/pkg/types"
)

// Handler handles RPC methods
type Handler struct {
	blockchain *core.Blockchain
}

// NewHandler creates a new handler
func NewHandler(blockchain *core.Blockchain) *Handler {
	return &Handler{
		blockchain: blockchain,
	}
}

// Handle routes RPC method to appropriate handler
func (h *Handler) Handle(req *RPCRequest) (interface{}, error) {
	switch req.Method {
	case "apex_blockNumber":
		return h.handleBlockNumber(req)
	case "apex_getBalance":
		return h.handleGetBalance(req)
	case "apex_getBlockByNumber":
		return h.handleGetBlockByNumber(req)
	case "apex_getBlockByHash":
		return h.handleGetBlockByHash(req)
	case "apex_getTransaction":
		return h.handleGetTransaction(req)
	case "apex_sendTransaction":
		return h.handleSendTransaction(req)
	case "apex_getValidators":
		return h.handleGetValidators(req)
	case "apex_stake":
		return h.handleStake(req)
	case "apex_unstake":
		return h.handleUnstake(req)
	case "apex_getStakingInfo":
		return h.handleGetStakingInfo(req)
	default:
		return nil, errors.New("method not found")
	}
}

// handleBlockNumber returns current block number
func (h *Handler) handleBlockNumber(req *RPCRequest) (interface{}, error) {
	height := h.blockchain.GetHeight()
	return map[string]interface{}{
		"blockNumber": height,
	}, nil
}

// handleGetBalance returns account balance
func (h *Handler) handleGetBalance(req *RPCRequest) (interface{}, error) {
	if len(req.Params) < 1 {
		return nil, errors.New("missing address parameter")
	}
	
	addrStr, ok := req.Params[0].(string)
	if !ok {
		return nil, errors.New("invalid address parameter")
	}
	
	addr := types.HexToAddress(addrStr)
	account, err := h.blockchain.GetStateDB().GetAccount(addr)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"address": addr.Hex(),
		"balance": types.FromWei(account.Balance),
		"staked":  types.FromWei(account.Staked),
		"locked":  types.FromWei(account.Locked),
		"nonce":   account.Nonce,
	}, nil
}

// handleGetBlockByNumber returns block by number
func (h *Handler) handleGetBlockByNumber(req *RPCRequest) (interface{}, error) {
	if len(req.Params) < 1 {
		return nil, errors.New("missing block number parameter")
	}
	
	blockNum, ok := req.Params[0].(float64)
	if !ok {
		return nil, errors.New("invalid block number parameter")
	}
	
	block, err := h.blockchain.GetBlockByNumber(uint64(blockNum))
	if err != nil {
		return nil, err
	}
	
	return h.formatBlock(block), nil
}

// handleGetBlockByHash returns block by hash
func (h *Handler) handleGetBlockByHash(req *RPCRequest) (interface{}, error) {
	if len(req.Params) < 1 {
		return nil, errors.New("missing block hash parameter")
	}
	
	hashStr, ok := req.Params[0].(string)
	if !ok {
		return nil, errors.New("invalid block hash parameter")
	}
	
	hash := types.HexToHash(hashStr)
	block, err := h.blockchain.GetBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	
	return h.formatBlock(block), nil
}

// handleGetValidators returns active validators
func (h *Handler) handleGetValidators(req *RPCRequest) (interface{}, error) {
	validators, err := h.blockchain.GetStateDB().GetAllValidators()
	if err != nil {
		return nil, err
	}
	
	result := make([]map[string]interface{}, len(validators))
	for i, val := range validators {
		result[i] = map[string]interface{}{
			"address":      val.Address.Hex(),
			"voting_power": types.FromWei(val.VotingPower),
			"commission":   float64(val.Commission) / 100,
			"status":       val.Status,
			"jailed":       val.Jailed,
		}
	}
	
	return result, nil
}

// formatBlock formats block for RPC response
func (h *Handler) formatBlock(block *core.Block) map[string]interface{} {
	txs := make([]string, len(block.Transactions))
	for i, tx := range block.Transactions {
		txs[i] = tx.Hash.Hex()
	}
	
	return map[string]interface{}{
		"number":           block.Header.Number,
		"hash":             block.Hash.Hex(),
		"previousHash":     block.Header.PreviousHash.Hex(),
		"timestamp":        block.Header.Timestamp,
		"validator":        block.Header.Validator.Hex(),
		"transactionRoot":  block.Header.TransactionRoot.Hex(),
		"stateRoot":        block.Header.StateRoot.Hex(),
		"gasUsed":          block.Header.GasUsed,
		"gasLimit":         block.Header.GasLimit,
		"transactions":     txs,
		"transactionCount": len(txs),
	}
}

// Stub handlers for transaction operations
func (h *Handler) handleGetTransaction(req *RPCRequest) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) handleSendTransaction(req *RPCRequest) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) handleStake(req *RPCRequest) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) handleUnstake(req *RPCRequest) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (h *Handler) handleGetStakingInfo(req *RPCRequest) (interface{}, error) {
	return nil, errors.New("not implemented")
}