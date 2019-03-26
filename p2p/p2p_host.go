// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"context"
	"strconv"

	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
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

	for _, addr := range BootstrapNodes {

	}
}

/* END EXPORTED METHODS */
