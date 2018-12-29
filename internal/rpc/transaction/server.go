package transaction

import (
	"context"
	"fmt"

	"github.com/space55/summertech-blockchain/accounts"
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

	err = transaction.WriteToMemory() // Write transaction to persistent memory

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

// Publish - transaction.Publish RPC handler
func (server *Server) Publish(ctx context.Context, req *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	hash, err := common.StringToHash(req.Address) // String to hash

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := types.ReadTransactionFromMemory(hash) // Read transaction from hash

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	err = transaction.Publish() // Publish transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\npublished transaction %s", transaction.Hash)}, nil // Return response
}

// Bytes - transaction.Bytes RPC handler
func (server *Server) Bytes(ctx context.Context, req *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	hash, err := common.StringToHash(req.Address) // String to hash

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := types.ReadTransactionFromMemory(hash) // Read transaction from hash

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	hex, err := common.EncodeString(transaction.Bytes()) // Encode byte value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\n%s", hex)}, nil // Return response
}

// String - transaction.String RPC handler
func (server *Server) String(ctx context.Context, req *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	hash, err := common.StringToHash(req.Address) // String to hash

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := types.ReadTransactionFromMemory(hash) // Read transaction from hash

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\n%s", transaction.String())}, nil // Return response
}

// SignTransaction - transaction.SignTransaction RPC handler
func (server *Server) SignTransaction(ctx context.Context, req *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	hash, err := common.StringToHash(req.Address) // String to hash

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := types.ReadTransactionFromMemory(hash) // Read transaction from hash

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(*transaction.Sender) // Read account

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	err = types.SignTransaction(transaction, account.PrivateKey) // Sign transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	err = transaction.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\n%s", transaction.Signature.String())}, nil // Return response
}

// VerifyTransactionSignature - transaction.VerifyTransactionSignature RPC handler
func (server *Server) VerifyTransactionSignature(ctx context.Context, req *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	hash, err := common.StringToHash(req.Address) // Get hash value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := types.ReadTransactionFromMemory(hash) // Read transaction from mempool

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	verified, err := types.VerifyTransactionSignature(transaction) // Verify signature

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\n%t", verified)}, nil // Return response
}
