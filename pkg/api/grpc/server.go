package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/apex/pkg/core"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server represents gRPC server
type Server struct {
	blockchain *core.Blockchain
	logger     *zap.Logger
	port       int
	server     *grpc.Server
}

// NewServer creates a new gRPC server
func NewServer(blockchain *core.Blockchain, port int, logger *zap.Logger) *Server {
	return &Server{
		blockchain: blockchain,
		logger:     logger,
		port:       port,
		server:     grpc.NewServer(),
	}
}

// Start starts the gRPC server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.logger.Info("Starting gRPC server", zap.String("address", addr))
	
	// Register services here
	// RegisterApexServiceServer(s.server, s)
	
	return s.server.Serve(listener)
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

// Example RPC method implementations
func (s *Server) GetBlockNumber(ctx context.Context, req interface{}) (interface{}, error) {
	height := s.blockchain.GetHeight()
	return map[string]interface{}{"height": height}, nil
}

func (s *Server) GetBalance(ctx context.Context, req interface{}) (interface{}, error) {
	// Implementation would go here
	return map[string]interface{}{"balance": "0"}, nil
}