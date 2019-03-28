// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"github.com/SummerCash/go-summercash/validator"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
)

// Client represents a peer on the network with a known routed libp2p host.
type Client struct {
	Host *routed.RoutedHost `json:"host"` // Host

	Validator *validator.Validator `json:"validator"` // Validator
}

/* BEGIN EXPORTED METHODS */

// NewClient initializes a new client with a given host.
func NewClient(host *routed.RoutedHost, validator *validator.Validator) *Client {
	return &Client{
		Host:      host,      // Set host
		Validator: validator, // Set validator
	} // Return initialized client
}

/* END EXPORTED METHODS */
