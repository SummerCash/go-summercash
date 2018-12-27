package accounts

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/space55/summertech-blockchain/accounts"
	"github.com/space55/summertech-blockchain/common"
	accountsProto "github.com/space55/summertech-blockchain/internal/rpc/proto/accounts"
)

// Server - RPC server
type Server struct{}

// NewAccount - accounts.NewAccount RPC handler
func (server *Server) NewAccount(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	account, err := accounts.NewAccount() // Create new account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", account.Address.String())}, nil // No error occurred, return response
}

// AccountFromKey - accounts.AccountFromKey RPC handler
func (server *Server) AccountFromKey(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	data, err := common.DecodeString(req.PrivateKey) // Decode private key

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	blockPrivate, _ := pem.Decode(data) // Decode

	x509EncodedPrivateKey := blockPrivate.Bytes // Get x509 byte val

	privateKey, err := x509.ParseECPrivateKey(x509EncodedPrivateKey) // Parse private key

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	account, err := accounts.AccountFromKey(privateKey) // Get account

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", account.Address.String())}, nil // No error occurred, return response
}

// GetAllAccounts - accounts.GetAllAccounts RPC handler
func (server *Server) GetAllAccounts(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	addresses, err := accounts.GetAllAccounts() // Walk

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", strings.Join(addresses, ", "))}, nil // No error occurred, return response
}

// MakeEncodingSafe - accounts.MakeEncodingSafe RPC handler
func (server *Server) MakeEncodingSafe(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	err = account.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\nmade account with address %s encoding safe", account.Address.String())}, nil // No error occurred, return response
}

// RecoverSafeEncoding - accounts.RecoverSafeEncoding RPC handler
func (server *Server) RecoverSafeEncoding(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	err = account.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\nrecovered account with address %s from safe encoding", account.Address.String())}, nil // No error occurred, return response
}

// String - accounts.String RPC handler
func (server *Server) String(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", account.String())}, nil // No error occurred, return response
}

// Bytes - accounts.Bytes RPC handler
func (server *Server) Bytes(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	encoded, err := common.EncodeString(account.Bytes()) // Encode byte val

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", encoded)}, nil // No error occurred, return response
}

// ReadAccountFromMemory - account.ReadAccountFromMemory RPC handler
func (server *Server) ReadAccountFromMemory(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", account.String())}, nil // No error occurred, return response
}
