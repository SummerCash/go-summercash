package p2p

import (
	"context"
	"fmt"
	"strings"

	"github.com/SummerCash/go-summercash/config"
	p2pProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/p2p"
	p2pPkg "github.com/SummerCash/go-summercash/p2p"
	"github.com/SummerCash/go-summercash/validator"
)

// Server - RPC server
type Server struct{}

// NumConnectedPeers - p2p.NumConnectedPeers RPC handler
func (server *Server) NumConnectedPeers(ctx context.Context, req *p2pProto.GeneralRequest) (*p2pProto.GeneralResponse, error) {
	if p2pPkg.WorkingHost == nil { // Check no working host
		return &p2pProto.GeneralResponse{}, p2pPkg.ErrNoWorkingHost // Return error
	}

	numPeers := 0 // Initialize peer num

	for _, peer := range p2pPkg.WorkingHost.Network().Peers() { // Iterate through peers
		if peer != p2pPkg.WorkingHost.ID() { // Check is foreign peer
			numPeers++ // Increment number of peers
		}
	}

	return &p2pProto.GeneralResponse{Message: fmt.Sprintf("\n%d", numPeers)}, nil // Return num of peers
}

// ConnectedPeers - p2p.ConnectedPeers RPC handler.
func (server *Server) ConnectedPeers(ctx context.Context, req *p2pProto.GeneralRequest) (*p2pProto.GeneralResponse, error) {
	if p2pPkg.WorkingHost == nil { // Check no working host
		return &p2pProto.GeneralResponse{}, p2pPkg.ErrNoWorkingHost // Return error
	}

	peers := []string{} // Initialize peer buffer

	for _, peerInfo := range p2pPkg.WorkingHost.Network().Peers() {
		if peerInfo != p2pPkg.WorkingHost.ID() { // Check is foreign peer
			peers = append(peers, peerInfo.String()) // Append peer
		}
	}

	return &p2pProto.GeneralResponse{Message: fmt.Sprintf("\n%s", strings.Join(peers, ", "))}, nil // Return peers
}

// SyncNetwork - p2p.SyncNetwork RPC handler
func (server *Server) SyncNetwork(ctx context.Context, req *p2pProto.GeneralRequest) (*p2pProto.GeneralResponse, error) {
	if p2pPkg.WorkingHost == nil { // Check no working host
		return &p2pProto.GeneralResponse{}, p2pPkg.ErrNoWorkingHost // Return error
	}

	config, err := config.ReadChainConfigFromMemory() // Read config from memory

	if err != nil { // Check for errors
		return &p2pProto.GeneralResponse{}, err // Return found error
	}

	validator := validator.Validator(validator.NewStandardValidator(config)) // Initialize validator

	client := p2pPkg.NewClient(p2pPkg.WorkingHost, &validator, req.Network) // Initialize p2p client

	err = client.SyncNetwork() // Sync network

	if err != nil { // Check for errors
		return &p2pProto.GeneralResponse{}, err // Return found error
	}

	return &p2pProto.GeneralResponse{Message: "\nSuccessful"}, nil // Return response
}
