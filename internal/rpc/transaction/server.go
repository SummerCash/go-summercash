package transaction

import (
	"context"
	"fmt"

	"github.com/space55/summertech-blockchain/common"
	transactionProto "github.com/space55/summertech-blockchain/internal/rpc/proto/transaction"
	"github.com/space55/summertech-blockchain/types"
)

// Server - RPC server
type Server struct{}

// NewTransaction - transaction.NewTransaction RPC handler
func (server *Server) NewTransaction(ctx context.Context, req *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	sender, err := common.StringToAddress(req.Address) // Convert address param to address literal

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	recipient, err := common.StringToAddress(req.Address2) // Check for errors

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	accountChain, err := types.ReadChainFromMemory(sender) // Read account chain from persistent memory

	if err != nil { // CHeck for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	nonce := uint64(0)                      // Init nonce
	lastTransaction := &types.Transaction{} // Init buffer

	for _, transaction := range accountChain.Transactions { // Iterate through transactions
		if *transaction.Sender == sender { // Check match
			if transaction.AccountNonce == uint64(len(accountChain.Transactions)) { // Check is last transaction
				lastTransaction = transaction // Set last transaction
			}

			nonce++ // Increment
		}
	}

	transaction, err := types.NewTransaction(nonce, lastTransaction, &sender, &recipient, req.Amount, req.Payload) // Init transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\n%s", transaction.String())}, nil // Return response
}

// TransactionFromBytes - transaction.TransactionFromBytes RPC handler
func (server *Server) TransactionFromBytes(ctx context.Context, req *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	transaction, err := types.TransactionFromBytes(req.Payload) // Get tx literal

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\n%s", transaction.String())}, nil // Return response
}
