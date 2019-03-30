// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"context"
	"strings"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var (
	// BootstrapNodes represents all default bootstrap nodes on the given network.
	BootstrapNodes = []string{
		"/ip4/108.41.124.60/udp/3000/quic/ipfs/QmeUXECRX2awd3WeJ7xVvaYJ8u1UGKDoiPSS4XTcTBZ5YF",
	}
)

/* BEGIN EXPORTED METHODS */

// GetBestBootstrapAddress attempts to fetch the best bootstrap node.
func GetBestBootstrapAddress(ctx context.Context, host *routed.RoutedHost) string {
	for _, bootstrapAddress := range BootstrapNodes { // Iterate through bootstrap nodes
		multiaddr, err := multiaddr.NewMultiaddr(bootstrapAddress) // Parse address

		if err != nil { // Check for errors
			continue // Continue
		}

		peerID, err := peer.IDB58Decode(strings.Split(bootstrapAddress, "ipfs/")[1]) // Get peer ID

		if err != nil { // Check for errors
			continue // Continue
		}

		host.Peerstore().AddAddr(peerID, multiaddr, 10*time.Second) // Add bootstrap peer

		peerInfo, err := peerstore.InfoFromP2pAddr(multiaddr) // Get peer info

		if err != nil { // Check for errors
			continue // Continue
		}

		bootstrapCheckCtx, cancel := context.WithCancel(ctx) // Get context

		err = host.Connect(bootstrapCheckCtx, *peerInfo) // Connect to peer

		if err != nil { // Check for errors
			cancel() // Cancel
			continue // Continue
		}

		_, err = ping.Ping(bootstrapCheckCtx, host, peerID) // Attempt to ping

		if err == nil { // Check no errors
			cancel()                // Cancel
			return bootstrapAddress // Return bootstrap address
		}

		cancel() // Cancel
	}

	return "localhost" // Return localhost
}

/* END EXPORTED METHODS */
