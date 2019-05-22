package p2p

import (
	"context"
	"fmt"

	p2pProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/p2p"
	p2pPkg "github.com/SummerCash/go-summercash/p2p"
)

// Server - RPC server
type Server struct{}

// ConnectedPeers - p2p.ConnectedPeers RPC handler
func (server *Server) ConnectedPeers(ctx context.Context, req *p2pProto.GeneralRequest) (*p2pProto.GeneralResponse, error) {
	if p2pPkg.WorkingHost == nil { // Check no working host
		return &p2pProto.GeneralResponse{}, p2pPkg.ErrNoWorkingHost // Return error
	}

	return &p2pProto.GeneralResponse{Message: fmt.Sprintf("%d", len(p2pPkg.WorkingHost.Peerstore().Peers()))}, nil // Return num of peers
}
