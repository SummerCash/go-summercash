// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"fmt"

	"github.com/SummerCash/go-summercash/common"

	inet "github.com/libp2p/go-libp2p-net"
	protocol "github.com/libp2p/go-libp2p-protocol"
)

// Stream header protocol definitions
const (
	PublishTransaction StreamHeaderProtocol = iota

	RequestConfig

	RequestBestTransaction

	RequestTransaction

	RequestGenesisHash

	RequestChildHashes

	RequestChain

	RequestAllChains

	RequestNextTransaction
)

var (
	// StreamHeaderProtocolNames represents all stream header protocol names.
	StreamHeaderProtocolNames = []string{
		"pub_transaction",
		"req_config",
		"req_best_transaction",
		"req_transaction",
		"req_genesis_hash",
		"req_transaction_children_hashes",
		"req_chain",
		"req_all_chains",
		"req_next_transaction",
	}
)

// StreamHeaderProtocol represents the stream protocol type enum.
type StreamHeaderProtocol int

/* BEGIN EXPORTED METHODS */

// StartServingStreams starts serving all necessary strings.
func (client *Client) StartServingStreams() error {
	common.Logf("== P2P == starting node stream handlers") // Log start handlers

	network := client.Network // Get network

	err := client.StartServingStream(GetStreamHeaderProtocolPath(network, PublishTransaction), client.HandleReceiveTransaction) // Start serving pub tx

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestChain), client.HandleReceiveChainRequest) // Start serving request chain

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestAllChains), client.HandleReceiveAllChainsRequest) // Start serving request all chains

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestNextTransaction), client.HandleReceiveNextTransactionRequest) // Start serving request next tx

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestBestTransaction), client.HandleReceiveBestTransaction) // Start serving request best tx

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// StartServingStream starts serving a given stream.
func (client *Client) StartServingStream(streamHeaderProtocolPath string, handler func(inet.Stream)) error {
	if WorkingHost == nil { // Check no host
		return ErrNoWorkingHost // Return found error
	}

	WorkingHost.SetStreamHandler(protocol.ID(streamHeaderProtocolPath), handler) // Set handler

	return nil // No error occurred, return nil
}

// GetStreamHeaderProtocolPath attempts to determine the libp2p stream header protocol URI from a given stream protocol and network.
func GetStreamHeaderProtocolPath(network string, streamProtocol StreamHeaderProtocol) string {
	return fmt.Sprintf("/%s/%s", network, StreamHeaderProtocolNames[streamProtocol]) // Return URI
}

/* END EXPORTED METHODS */
