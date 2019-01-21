package transaction

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SummerCash/go-summercash/accounts"
	"github.com/SummerCash/go-summercash/common"
	transactionProto "github.com/SummerCash/go-summercash/internal/rpc/proto/transaction"
	"github.com/SummerCash/go-summercash/types"
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

	transaction := types.Transaction{} // Init buffer

	accountChain, err := types.ReadChainFromMemory(sender) // Read account chain from persistent memory

	if err != nil { // Check for errors
		newTransaction, err := types.NewTransaction(0, nil, &sender, &recipient, req.Amount, req.Payload) // Init transaction

		if err != nil { // Check for errors
			return &transactionProto.GeneralResponse{}, err // Return found error
		}

		transaction = *newTransaction // Write tx to buffer
	} else {
		nonce := uint64(0)                      // Init nonce
		lastTransaction := &types.Transaction{} // Init buffer

		if len(accountChain.Transactions) > 0 { // Check is not 0 index
			for _, transaction := range accountChain.Transactions { // Iterate through transactions
				if *transaction.Recipient == sender || *transaction.Sender == sender { // Check match
					if transaction.AccountNonce == uint64(len(accountChain.Transactions)) { // Check is last transaction
						lastTransaction = transaction // Set last transaction
					}

					nonce++ // Increment
				}
			}
		}

		newTransaction, err := types.NewTransaction(nonce, lastTransaction, &sender, &recipient, req.Amount, req.Payload) // Init transaction

		if err != nil { // Check for errors
			return &transactionProto.GeneralResponse{}, err // Return found error
		}

		transaction = *newTransaction // Write tx to buffer
	}

	err = transaction.WriteToMemory() // Write transaction to persistent memory

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\n%s\n\nhash: %s", transaction.String(), transaction.Hash.String())}, nil // Return response
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

	if transaction.Payload != nil && strings.Contains(string(transaction.Payload), "(") { // Check is contract call
		return handleContractCall(transaction) // Handle contract call
	}

	err = transaction.Publish() // Publish transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\npublished transaction %s", transaction.Hash)}, nil // Return response
}

func handleContractCall(transaction *types.Transaction) (*transactionProto.GeneralResponse, error) {
	err := transaction.Publish() // Publish transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	chain, err := types.ReadChainFromMemory(*transaction.Recipient) // Read recipient chain

	if err != nil { // Check for errors
		coordinationChain, err := types.ReadCoordinationChainFromMemory() // Read coordination chain

		if err != nil { // Check for errors
			return &transactionProto.GeneralResponse{}, err // Return found error
		}

		chain, err = coordinationChain.GetChain(*transaction.Recipient) // Get contract chain

		if err != nil { // Check for errors
			return &transactionProto.GeneralResponse{}, err // Return found error
		}

		err = chain.WriteToMemory() // Write chain to persistent memory

		if err != nil { // Check for errors
			return &transactionProto.GeneralResponse{}, err // Return found error
		}
	}

	startTime := time.Now() // Get start time

	txHash := *transaction.Hash         // Get hash
	recipient := *transaction.Recipient // Get recipient
	transaction = nil                   // Init tx buffer

	for time.Now().Sub(startTime) < 5*time.Second { // Async read tx
		chain, err = types.ReadChainFromMemory(recipient) // Read recipient chain

		if err != nil { // Check for errors
			return &transactionProto.GeneralResponse{}, err // Return found error
		}

		transaction, _ = chain.QueryTransaction(txHash) // Query TX

		if transaction != nil { // Check not nil
			break // Break
		}
	}

	if transaction == nil { // Check tx nil
		return &transactionProto.GeneralResponse{}, types.ErrNilTransaction // Return error
	}

	if transaction.Logs == nil || len(transaction.Logs) == 0 { // Check no logs
		return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\npublished transaction %s", transaction.Hash)}, nil // Return response
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("\ncontract call response: %s", types.StringLogs(transaction.Logs))}, nil // Return response
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
