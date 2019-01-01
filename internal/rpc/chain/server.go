package chain

import (
	"context"
	"fmt"

	"github.com/SummerCash/go-summercash/common"
	chainProto "github.com/SummerCash/go-summercash/internal/rpc/proto/chain"
	"github.com/SummerCash/go-summercash/types"
)

// Server - RPC server
type Server struct{}

// GetBalance - chain.GetBalance RPC handler
func (server *Server) GetBalance(ctx context.Context, req *chainProto.GeneralRequest) (*chainProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address primitive value

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	chain, err := types.ReadChainFromMemory(address) // Read chain from persistent memory

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	balance := chain.CalculateBalance() // Get balance

	return &chainProto.GeneralResponse{Message: fmt.Sprintf("\nbalance: %f", balance)}, nil // Return response
}
