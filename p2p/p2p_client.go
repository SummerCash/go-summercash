// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"bytes"
	"context"

	"github.com/SummerCash/go-summercash/crypto"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
	"github.com/SummerCash/go-summercash/validator"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
)

// Client represents a peer on the network with a known routed libp2p host.
type Client struct {
	Host *routed.RoutedHost `json:"host"` // Host

	Validator *validator.Validator `json:"validator"` // Validator

	Network string `json:"network"` // Network
}

/* BEGIN EXPORTED METHODS */

// NewClient initializes a new client with a given host.
func NewClient(host *routed.RoutedHost, validator *validator.Validator, network string) *Client {
	return &Client{
		Host:      host,      // Set host
		Validator: validator, // Set validator
		Network:   network,   // Set network
	} // Return initialized client
}

// SyncNetwork syncs all available chains and state roots.
func (client *Client) SyncNetwork() error {
	// TODO: Implement
	return nil // No error occurred, return nil
}

// RequestChain requests a chain from the working network with a given sample size.
func (client *Client) RequestChain(account common.Address, sampleSize uint) (*types.Chain, error) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	responses, err := BroadcastDhtResult(ctx, client.Host, account.Bytes(), GetStreamHeaderProtocolPath(client.Network, RequestChain), client.Network, 16) // Broadcast, get result

	if err != nil { // Check for errors
		return &types.Chain{}, err // Return found error
	}

	occurrences := make(map[common.Hash]int64) // Init occurrences buffer

	var bestResponse []byte // Init best response buffer

	for _, response := range responses { // Iterate through responses
		if len(response) == 0 || response == nil || bytes.Equal(response, make([]byte, len(response))) { // Check is nil
			continue // Continue
		}

		occurrences[common.NewHash(crypto.Sha3(response))]++ // Increment occurrences

		if occurrences[common.NewHash(crypto.Sha3(response))] > occurrences[common.NewHash(crypto.Sha3(bestResponse))] { // Check is better response
			bestResponse = response // Set best response
		}
	}

	chain, err := types.FromBytes(bestResponse) // Deserialize chain

	if err != nil { // Check for errors
		return &types.Chain{}, err // Return found error
	}

	return chain, nil // Return chain
}

/* END EXPORTED METHODS */
