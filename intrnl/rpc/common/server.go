package common

import (
	"context"
	"fmt"

	"github.com/SummerCash/go-summercash/common"
	commonProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/common"
)

// Server - RPC server
type Server struct{}

// Encode - common.Encode RPC handler
func (server *Server) Encode(ctx context.Context, req *commonProto.GeneralRequest) (*commonProto.GeneralResponse, error) {
	encoded, err := common.Encode(req.Input) // Encode

	if err != nil { // Check for errors
		return &commonProto.GeneralResponse{}, err // Return found error
	}

	return &commonProto.GeneralResponse{Message: fmt.Sprintf("\n%s", encoded)}, nil // Return response
}

// EncodeString - common.EncodeString RPC handler
func (server *Server) EncodeString(ctx context.Context, req *commonProto.GeneralRequest) (*commonProto.GeneralResponse, error) {
	encoded, err := common.EncodeString(req.Input) // Encode

	if err != nil { // Check for errors
		return &commonProto.GeneralResponse{}, err // Return found error
	}

	return &commonProto.GeneralResponse{Message: fmt.Sprintf("\n%s", encoded)}, nil // Return response
}

// Decode - common.Decode RPC handler
func (server *Server) Decode(ctx context.Context, req *commonProto.GeneralRequest) (*commonProto.GeneralResponse, error) {
	decoded, err := common.Decode(req.Input) // Decode

	if err != nil { // Check for errors
		return &commonProto.GeneralResponse{}, err // Return found error
	}

	return &commonProto.GeneralResponse{Message: fmt.Sprintf("\n%s", decoded)}, nil // Return response
}

// DecodeString - common.DecodeString RPC handler
func (server *Server) DecodeString(ctx context.Context, req *commonProto.GeneralRequest) (*commonProto.GeneralResponse, error) {
	decoded, err := common.DecodeString(req.S) // Decode

	if err != nil { // Check for errors
		return &commonProto.GeneralResponse{}, err // Return found error
	}

	return &commonProto.GeneralResponse{Message: fmt.Sprintf("\n%s", decoded)}, nil // Return response
}
