package coordinationchain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SummerCash/go-summercash/common"
	coordinationChainProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/coordinationchain"
	"github.com/SummerCash/go-summercash/types"
	commonGoP2P "github.com/dowlandaiello/GoP2P/common"
)

// Server - RPC server
type Server struct{}

// SyncNetwork - coordinationChain.SyncNetwork RPC handler
func (server *Server) SyncNetwork(ctx context.Context, req *coordinationChainProto.GeneralRequest) (*coordinationChainProto.GeneralResponse, error) {
	err := types.SyncNetwork(true, true) // Sync network

	if err != nil { // Check for errors
		return &coordinationChainProto.GeneralResponse{}, err // Return found error
	}

	chain, err := types.ReadCoordinationChainFromMemory() // Read chain

	if err != nil { // Check for errors
		return &coordinationChainProto.GeneralResponse{}, err // Return found error
	}

	return &coordinationChainProto.GeneralResponse{Message: fmt.Sprintf("\nfinished sync with chain %s", chain.ChainID.String())}, nil // Return response
}

// GetPeers - get all peers in coordination chain
func (server *Server) GetPeers(ctx context.Context, req *coordinationChainProto.GeneralRequest) (*coordinationChainProto.GeneralResponse, error) {
	chain, err := types.ReadCoordinationChainFromMemory() // Read chain

	if err != nil { // Check for errors
		return &coordinationChainProto.GeneralResponse{}, err // Return found error
	}

	peerCount := 0                                   // Init buffer
	knownPeers := []string{common.BootstrapNodes[0]} // Init buffer

	for _, node := range chain.Nodes { // Iterate through nodes
		for _, peer := range node.Addresses { // Iterate through peers in node
			if !commonGoP2P.StringInSlice(knownPeers, peer) { // Check peer is not already known
				knownPeers = append(knownPeers, peer) // Add known peer

				peerCount++ // Increment
			}
		}
	}

	json, err := json.MarshalIndent(knownPeers, "", "  ") // Marshal

	if err != nil { // Check for errors
		return &coordinationChainProto.GeneralResponse{}, err // Return found error
	}

	return &coordinationChainProto.GeneralResponse{Message: fmt.Sprintf("\nfound %d connected peers: %s", peerCount, string(json))}, nil // Return response
}

// Bytes - coordinationChain.Bytes RPC handler
func (server *Server) Bytes(ctx context.Context, req *coordinationChainProto.GeneralRequest) (*coordinationChainProto.GeneralResponse, error) {
	chain, err := types.ReadCoordinationChainFromMemory() // Read chain

	if err != nil { // Check for errors
		return &coordinationChainProto.GeneralResponse{}, err // Return found error
	}

	hex, err := common.EncodeString(chain.Bytes()) // Encode chain byte value to hex

	if err != nil { // Check for errors
		return &coordinationChainProto.GeneralResponse{}, err // Return found error
	}

	return &coordinationChainProto.GeneralResponse{Message: fmt.Sprintf("\n%s", hex)}, nil // Return response
}

// String - coordinationChain.String RPC handler
func (server *Server) String(ctx context.Context, req *coordinationChainProto.GeneralRequest) (*coordinationChainProto.GeneralResponse, error) {
	chain, err := types.ReadCoordinationChainFromMemory() // Read chain

	if err != nil { // Check for errors
		return &coordinationChainProto.GeneralResponse{}, err // Return found error
	}

	return &coordinationChainProto.GeneralResponse{Message: fmt.Sprintf("\n%s", chain.String())}, nil // Return response
}
