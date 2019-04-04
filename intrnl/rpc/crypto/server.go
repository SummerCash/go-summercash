package crypto

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/SummerCash/go-summercash/crypto"
	cryptoProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/crypto"
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

// Sha3N - crypto.Sha3n RPC handler
func (server *Server) Sha3N(ctx context.Context, req *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	result := crypto.Sha3n(req.Input, uint(req.N)) // Hash input

	return &cryptoProto.GeneralResponse{Message: fmt.Sprintf("\n%s", hex.EncodeToString(result))}, nil // Return response
}

// Sha3NString - crypto.Sha3nString RPC handler
func (server *Server) Sha3NString(ctx context.Context, req *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	result := crypto.Sha3nString(req.Input, uint(req.N)) // Hash input

	return &cryptoProto.GeneralResponse{Message: fmt.Sprintf("\n%s", result)}, nil // Return response
}

// Sha3D - crypto.Sha3d RPC handler
func (server *Server) Sha3D(ctx context.Context, req *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	result := crypto.Sha3d(req.Input) // Hash input

	return &cryptoProto.GeneralResponse{Message: fmt.Sprintf("\n%s", hex.EncodeToString(result))}, nil // Return response
}

// Sha3DString - crypto.Sha3dString RPC handler
func (server *Server) Sha3DString(ctx context.Context, req *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	result := crypto.Sha3dString(req.Input) // Hash input

	return &cryptoProto.GeneralResponse{Message: fmt.Sprintf("\n%s", result)}, nil // Return response
}
