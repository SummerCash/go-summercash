package config

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	configProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/config"
)

// Server - RPC server
type Server struct{}

// NewChainConfig - config.NewChainConfig RPC handler
func (server *Server) NewChainConfig(ctx context.Context, req *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	chainConfig, err := config.NewChainConfig(req.GenesisPath) // Init config
	if err != nil {                                            // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	oldDataDir := common.DataDir // Set old data dir

	if req.GenesisPath == "examples/genesis.json" { // Check is example
		common.DataDir, _ = filepath.Abs("./examples") // Temp data
	}

	err = chainConfig.WriteToMemory() // Write to persistent memory

	if err != nil { // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	common.DataDir = oldDataDir // Reset

	return &configProto.GeneralResponse{Message: fmt.Sprintf("\n%s", chainConfig.String())}, nil // Return response
}

// Bytes - config.Bytes RPC handler
func (server *Server) Bytes(ctx context.Context, req *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	chainConfig, err := config.ReadChainConfigFromMemory() // Read chain config from memory
	if err != nil {                                        // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	hex, err := common.EncodeString(chainConfig.Bytes()) // Encode bytes to string
	if err != nil {                                      // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	return &configProto.GeneralResponse{Message: fmt.Sprintf("\n%s", hex)}, nil // Return response
}

// String - config.String RPC handler
func (server *Server) String(ctx context.Context, req *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	chainConfig, err := config.ReadChainConfigFromMemory() // Read chain config from memory
	if err != nil {                                        // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	return &configProto.GeneralResponse{Message: fmt.Sprintf("\n%s", chainConfig.String())}, nil // Return response
}

// WriteToMemory - config.WriteToMemory RPC handler
func (server *Server) WriteToMemory(ctx context.Context, req *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	chainConfig, err := config.ReadChainConfigFromMemory() // Read chain config from memory
	if err != nil {                                        // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	err = chainConfig.WriteToMemory() // Write chain config to memory

	if err != nil { // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	return &configProto.GeneralResponse{Message: fmt.Sprintf("\nwrote config %s to memory at dir %s", chainConfig.ChainID.String(), fmt.Sprintf("%s/config/config.json", common.DataDir))}, nil // Return response
}

// ReadChainConfigFromMemory - config.ReadChainConfigFromMemory RPC handler
func (server *Server) ReadChainConfigFromMemory(ctx context.Context, req *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	chainConfig, err := config.ReadChainConfigFromMemory() // Read chain config from memory
	if err != nil {                                        // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	return &configProto.GeneralResponse{Message: fmt.Sprintf("\n%s", chainConfig.String())}, nil // Return response
}

// GetInflationRate - config.GetInflationRate RPC handler
func (server *Server) GetInflationRate(ctx context.Context, req *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	chainConfig, err := config.ReadChainConfigFromMemory() // Read chain config from memory
	if err != nil {                                        // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	return &configProto.GeneralResponse{Message: fmt.Sprintf("\n%f", chainConfig.InflationRate)}, nil // Return response
}

// GetTotalSupply - config.GetTotalSupply RPC handler
func (server *Server) GetTotalSupply(ctx context.Context, req *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	chainConfig, err := config.ReadChainConfigFromMemory() // Read chain config from memory
	if err != nil {                                        // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	supply := big.NewFloat(0) // Init supply

	for _, address := range chainConfig.AllocAddresses { // Iterate through alloc
		supply.Add(supply, chainConfig.Alloc[address.String()]) // Add alloc value
	}

	return &configProto.GeneralResponse{Message: fmt.Sprintf("\n%f", supply)}, nil // Return response
}
