// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"context"
	"strconv"

	"github.com/SummerCash/go-summercash/common"

	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var (
	// WorkingHost represents the global routed host.
	WorkingHost *routed.RoutedHost
)

/* BEGIN EXPORTED METHODS */

// NewHost initializes a new routed libp2p host with a given context.
func NewHost(ctx context.Context, port int) (*routed.RoutedHost, error) {
	identity, err := GetLibp2pPeerIdentity() // Get peer identity

	if err != nil { // Check for errors
		return &routed.RoutedHost{}, err // Return found error
	}

	host, err := libp2p.New(
		ctx,
		libp2p.NATPortMap(),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/udp/"+strconv.Itoa(port),
			"/ip6/::1/udp/"+strconv.Itoa(port)),
		libp2p.Identity(*identity),
		libp2p.Transport(quic.NewTransport),
	) // Initialize host

	if err != nil { // Check for errors
		return &routed.RoutedHost{}, err // Return found error
	}

	bootstrapCtx, cancel := context.WithCancel(ctx) // Get bootstrap context

	defer cancel() // Cancel

	dht, err := BootstrapDht(bootstrapCtx, host) // Bootstrap DHT

	if err != nil { // Check for errors
		return &routed.RoutedHost{}, err // Return found error
	}

	routedHost := routed.Wrap(host, dht) // Wrap host with DHT

	WorkingHost = routedHost // Set working host

	common.Logf("== P2P == initialized host with ID: %s on listening port: %d with multiaddr: %s", host.ID().Pretty(), strconv.Itoa(port), host.Addrs()[0].String()) // Log host

	return WorkingHost, nil // Return working routed host
}

// BootstrapDht bootstraps a KadDht to the list of bootstrap nodes.
func BootstrapDht(ctx context.Context, host host.Host) (*dht.IpfsDHT, error) {
	dht, err := dht.New(ctx, host) // Initialize DHT with host and context

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	err = dht.Bootstrap(ctx) // Bootstrap

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	for _, addr := range BootstrapNodes { // Iterate through bootstrap nodes
		address, err := multiaddr.NewMultiaddr(addr) // Parse multi address

		if err != nil { // Check for errors
			continue // Continue to next peer
		}

		peerInfo, err := peerstore.InfoFromP2pAddr(address) // Get peer info

		if err != nil { // Check for errors
			continue // Continue to next peer
		}

		err = host.Connect(ctx, *peerInfo) // Connect to discovered peer

		if err != nil { // Check for errors
			continue // Continue to next peer
		}
	}

	return dht, nil // Return DHT
}

/* END EXPORTED METHODS */
