package crypto

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/space55/summertech-blockchain/crypto"
	cryptoProto "github.com/space55/summertech-blockchain/internal/rpc/proto/crypto"
)

// Server - RPC server
type Server struct{}

// Sha3 - crypto.Sha3 RPC handler
func (server *Server) Sha3(ctx context.Context, req *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	result := crypto.Sha3(req.Input) // Hash input

	return &cryptoProto.GeneralResponse{Message: fmt.Sprintf("\n%s", hex.EncodeToString(result))}, nil // Return response
}

// Sha3String - crypto.Sha3String RPC handler
func (server *Server) Sha3String(ctx context.Context, req *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	result := crypto.Sha3String(req.Input) // Hash input

	return &cryptoProto.GeneralResponse{Message: fmt.Sprintf("\n%s", result)}, nil // Return response
}
