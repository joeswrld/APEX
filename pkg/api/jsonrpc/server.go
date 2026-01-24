package jsonrpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apex/pkg/core"
	"go.uber.org/zap"
)

// Server represents JSON-RPC server
type Server struct {
	blockchain *core.Blockchain
	handler    *Handler
	logger     *zap.Logger
	port       int
}

// NewServer creates a new JSON-RPC server
func NewServer(blockchain *core.Blockchain, port int, logger *zap.Logger) *Server {
	return &Server{
		blockchain: blockchain,
		handler:    NewHandler(blockchain),
		logger:     logger,
		port:       port,
	}
}

// Start starts the JSON-RPC server
func (s *Server) Start() error {
	http.HandleFunc("/", s.handleRequest)
	
	addr := fmt.Sprintf(":%d", s.port)
	s.logger.Info("Starting JSON-RPC server", zap.String("address", addr))
	
	return http.ListenAndServe(addr, nil)
}

// handleRequest handles incoming JSON-RPC requests
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req RPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, nil, -32700, "Parse error", err.Error())
		return
	}
	
	// Route to handler
	result, err := s.handler.Handle(&req)
	if err != nil {
		s.sendError(w, req.ID, -32603, "Internal error", err.Error())
		return
	}
	
	s.sendResponse(w, req.ID, result)
}

// sendResponse sends successful response
func (s *Server) sendResponse(w http.ResponseWriter, id interface{}, result interface{}) {
	resp := RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// sendError sends error response
func (s *Server) sendError(w http.ResponseWriter, id interface{}, code int, message, data string) {
	resp := RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RPCRequest represents JSON-RPC request
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

// RPCResponse represents JSON-RPC response
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}
