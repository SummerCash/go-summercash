// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/libp2p/go-libp2p"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
	multiaddr "github.com/multiformats/go-multiaddr"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
)

var (
	// WorkingHost represents the global routed host.
	WorkingHost *routed.RoutedHost

	// ErrNoWorkingHost represents an error describing a WorkingHost value of nil.
	ErrNoWorkingHost = errors.New("no working host")
)

/* BEGIN EXPORTED METHODS */

// NewHost initializes a new routed libp2p host with a given context.
func NewHost(ctx context.Context, port int) (*routed.RoutedHost, error) {
	identity, err := GetPeerIdentity() // Get peer identity

	if err != nil { // Check for errors
		return &routed.RoutedHost{}, err // Return found error
	}

	host, err := libp2p.New(
		ctx,
		libp2p.NATPortMap(),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/"+strconv.Itoa(port),
			"/ip6/::1/tcp/"+strconv.Itoa(port),
		),
		libp2p.Identity(*identity),
	) // Initialize host

	if err != nil { // Check for errors
		return &routed.RoutedHost{}, err // Return found error
	}

	common.Logf("== P2P == initialized host with ID: %s on listening port: %d with multiaddr: %s\n", host.ID().Pretty(), port, host.Addrs()[0].String()) // Log host

	dht, err := BootstrapDht(ctx, host) // Bootstrap DHT

	if err != nil { // Check for errors
		return &routed.RoutedHost{}, err // Return found error
	}

	routingDiscovery := discovery.NewRoutingDiscovery(dht) // Initialize routing discovery

	common.Logf("== P2P == advertising network presence\n") // Log advertise

	discovery.Advertise(ctx, routingDiscovery, config.Version) // Advertise network presence

	routedHost := routed.Wrap(host, dht) // Wrap host with DHT

	WorkingHost = routedHost // Set working host

	peerChan, err := routingDiscovery.FindPeers(ctx, config.Version) // Look for peers

	if err != nil { // Check for errors
		return &routed.RoutedHost{}, err // Return found error
	}

	common.Logf("== P2P == searching for remote nodes via rendezvous discovery...\n") // Log search

	for peer := range peerChan { // Iterate through discovered peers
		fmt.Println("test")
		if peer.ID == host.ID() { // Check is self
			continue // Skip
		}

		err = WorkingHost.Connect(ctx, peer) // Connect to discovered peer

		if err != nil { // Check for errors
			continue // Continue to next peer
		}
	}

	return WorkingHost, nil // Return working routed host
}

// BootstrapConfig bootstraps the network's working config with a given host.
func BootstrapConfig(ctx context.Context, host *routed.RoutedHost, bootstrapAddress string, network string) (*config.ChainConfig, error) {
	common.Logf("== P2P == bootstrapping config with bootstrap node address %s\n", bootstrapAddress) // Log bootstrap config

	peerID, err := peer.IDB58Decode(strings.Split(bootstrapAddress, "ipfs/")[1]) // Get peer ID

	if err != nil { // Check for errors
		return &config.ChainConfig{}, err // Return found error
	}

	readCtx, cancel := context.WithCancel(ctx) // Get context

	stream, err := (*host).NewStream(readCtx, peerID, protocol.ID(GetStreamHeaderProtocolPath(network, RequestConfig))) // Initialize new stream

	if err != nil { // Check for errors
		cancel() // Cancel

		return &config.ChainConfig{}, err // Return found error
	}

	reader := bufio.NewReader(stream) // Initialize reader from stream

	dagConfigBytes, err := reader.ReadBytes('\r') // Read

	if err != nil { // Check for errors
		cancel() // Cancel

		return &config.ChainConfig{}, err // Return found error
	}

	dagConfigBytes = bytes.Trim(dagConfigBytes, "\r") // Trim delimiter

	deserializedConfig, err := config.FromBytes(dagConfigBytes) // Deserialize

	if err != nil { // Check for errors
		cancel() // Cancel

		return &config.ChainConfig{}, err // Return found error
	}

	cancel() // Cancel

	common.Logf("== P2P == finished bootstrapping config\n") // Log finish bootstrap config

	return deserializedConfig, nil // Return deserialized dag config
}

// BootstrapDht bootstraps a KadDht to the list of bootstrap nodes.
func BootstrapDht(ctx context.Context, host host.Host) (*dht.IpfsDHT, error) {
	dht, err := dht.New(ctx, host) // Initialize DHT with host and context

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

	err = dht.Bootstrap(ctx) // Bootstrap

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return dht, nil // Return DHT
}

/* END EXPORTED METHODS */
