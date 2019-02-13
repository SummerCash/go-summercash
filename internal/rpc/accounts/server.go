package accounts

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/SummerCash/go-summercash/accounts"
	"github.com/SummerCash/go-summercash/common"
	accountsProto "github.com/SummerCash/go-summercash/internal/rpc/proto/accounts"
)

// Server - RPC server
type Server struct{}

// NewAccount - accounts.NewAccount RPC handler
func (server *Server) NewAccount(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	account, err := accounts.NewAccount() // Create new account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	err = account.WriteToMemory() // Write to persistent memory

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	marshaledPrivateKey, err := x509.MarshalECPrivateKey(account.PrivateKey) // Marshal private key

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshaledPrivateKey}) // Encode to memory

	stringEncoded, err := common.EncodeString(pemEncoded) // Encode to string

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\nAddress: %s, PrivateKey: %s", account.Address.String(), stringEncoded)}, nil // No error occurred, return response
}

// NewContractAccount - accounts.NewContractAccount RPC handler
func (server *Server) NewContractAccount(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	filepath, err := filepath.Abs(filepath.FromSlash(req.Address)) // Get filepath

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	contractSource, err := ioutil.ReadFile(filepath) // Read contract source

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	contractInstance, err := accounts.NewContractAccount(contractSource) // Deploy from contract source

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	marshaledPrivateKey, err := x509.MarshalECPrivateKey(contractInstance.PrivateKey) // Marshal private key

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshaledPrivateKey}) // Encode to memory

	stringEncoded, err := common.EncodeString(pemEncoded) // Encode to string

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\nAddress: %s, PrivateKey: %s", contractInstance.Address.String(), stringEncoded)}, nil // Return response
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
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := accounts.AccountFromKey(privateKey) // Get account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	err = account.WriteToMemory() // Write to persistent memory

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", account.Address.String())}, nil // No error occurred, return response
}

// GetAllAccounts - accounts.GetAllAccounts RPC handler
func (server *Server) GetAllAccounts(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	addresses, err := accounts.GetAllAccounts() // Walk

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", strings.Join(addresses, ", "))}, nil // No error occurred, return response
}

// GetAllContracts - accounts.GetAllContracts RPC handler
func (server *Server) GetAllContracts(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	panic("test")
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	addresses, err := accounts.GetAllContracts(address) // Walk with deploying address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	fmt.Println(strings.Join(addresses, ", "))

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", strings.Join(addresses, ", "))}, nil // No error occurred, return response
}

// MakeEncodingSafe - accounts.MakeEncodingSafe RPC handler
func (server *Server) MakeEncodingSafe(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	err = account.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	marshaledVal, err := json.MarshalIndent(*account, "", "  ") // Marshal JSON

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: string(marshaledVal)}, nil // No error occurred, return response
}

// RecoverSafeEncoding - accounts.RecoverSafeEncoding RPC handler
func (server *Server) RecoverSafeEncoding(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	err = account.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	marshaledVal, err := json.MarshalIndent(*account, "", "  ") // Marshal JSON

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: string(marshaledVal)}, nil // No error occurred, return response
}

// String - accounts.String RPC handler
func (server *Server) String(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", account.String())}, nil // No error occurred, return response
}

// Bytes - accounts.Bytes RPC handler
func (server *Server) Bytes(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	encoded, err := common.EncodeString(account.Bytes()) // Encode byte val

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\n%s", encoded)}, nil // No error occurred, return response
}

// ReadAccountFromMemory - account.ReadAccountFromMemory RPC handler
func (server *Server) ReadAccountFromMemory(ctx context.Context, req *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	address, err := common.StringToAddress(req.Address) // Get address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(address) // Read account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	err = account.MakeEncodingSafe() // Make safe for encoding

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	encoded, err := common.EncodeString(account.SerializedPrivateKey) // Encode private key

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	err = account.RecoverSafeEncoding() // Reverse

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: fmt.Sprintf("\nAddress: %s, PrivateKey: %s", account.Address, encoded)}, nil // No error occurred, return response
}
