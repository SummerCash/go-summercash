package accounts

import (
	"context"
	"crypto/x509"
	"encoding/pem"

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

	return &accountsProto.GeneralResponse{Message: account.Address.String()}, nil // No error occurred, return response
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

	return &accountsProto.GeneralResponse{Message: account.Address.String()}, nil // No error occurred, return response
}
