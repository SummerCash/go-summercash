package chain

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

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

// Bytes - chain.Bytes RPC handler
func (server *Server) Bytes(ctx context.Context, req *chainProto.GeneralRequest) (*chainProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address primitive value

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	chain, err := types.ReadChainFromMemory(address) // Read chain from persistent memory

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	hex, err := common.EncodeString(chain.Bytes()) // Encode chain byte value

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	return &chainProto.GeneralResponse{Message: fmt.Sprintf("\n%s", hex)}, nil // Return response
}

// String - chain.String RPC handler
func (server *Server) String(ctx context.Context, req *chainProto.GeneralRequest) (*chainProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address primitive value

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	chain, err := types.ReadChainFromMemory(address) // Read chain from persistent memory

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	return &chainProto.GeneralResponse{Message: fmt.Sprintf("\n%s", chain.String())}, nil // Return response
}

// ReadChainFromMemory - chain.ReadChainFromMemory RPC handler
func (server *Server) ReadChainFromMemory(ctx context.Context, req *chainProto.GeneralRequest) (*chainProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address primitive value

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	chain, err := types.ReadChainFromMemory(address) // Read chain from persistent memory

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	return &chainProto.GeneralResponse{Message: fmt.Sprintf("\n%s", chain.String())}, nil // Return response
}

// QueryTransaction - chain.QueryTransaction RPC handler
func (server *Server) QueryTransaction(ctx context.Context, req *chainProto.GeneralRequest) (*chainProto.GeneralResponse, error) {
	files, err := ioutil.ReadDir(filepath.FromSlash(fmt.Sprintf("%s/db/chain", common.DataDir))) // Walk chain dir

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	hash, err := common.StringToHash(req.Address) // Get hash value

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	for _, file := range files { // Iterate through files
		address, err := common.StringToAddress(strings.Split(strings.Split(file.Name(), "chain_")[1], ".json")[0]) // Get address value

		if err == nil { // Check for success
			chain, err := types.ReadChainFromMemory(address) // Read chain

			if err == nil { // Check successfully read
				transaction, err := chain.QueryTransaction(hash) // Query for transaction

				if err != nil { // Check for errors
					continue
				}

				return &chainProto.GeneralResponse{Message: fmt.Sprintf("\n%s", transaction.String())}, nil // Return response
			}
		}
	}

	return &chainProto.GeneralResponse{}, types.ErrNilTransaction // Return error
}

// GetNumTransactions - get total number of transactions in given account chain
func (server *Server) GetNumTransactions(ctx context.Context, req *chainProto.GeneralRequest) (*chainProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address primitive value

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	chain, err := types.ReadChainFromMemory(address) // Read chain from persistent memory

	if err != nil { // Check for errors
		return &chainProto.GeneralResponse{}, err // Return found error
	}

	numTx := 0 // Init counter

	for range chain.Transactions { // Iterate through transactions
		numTx++ // Increment
	}

	return &chainProto.GeneralResponse{Message: fmt.Sprintf("\n%d", numTx)}, nil // Return response
}
