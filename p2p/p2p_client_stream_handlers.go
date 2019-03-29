// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"bufio"
	"bytes"
	"strings"

	inet "github.com/libp2p/go-libp2p-net"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
)

/* BEGIN EXPORTED METHODS */

// HandleReceiveTransaction handles an incoming pub_tx stream.
func (client *Client) HandleReceiveTransaction(stream inet.Stream) {
	reader := bufio.NewReader(stream) // Initialize reader

	b, err := reader.ReadBytes('\f') // Read up to delimiter

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading pub_tx stream: %s", err.Error()) // Log error

		return // Return
	}

	b = bytes.Trim(b, "\f") // Trim delimiter

	tx, err := types.TransactionFromBytes(b) // Marshal bytes to transaction

	if err != nil { // Check for errors
		common.Logf("== P2P == error while deserializing tx read from pub_tx stream: %s", err.Error()) // Log error

		return // Return
	}

	tx.RecoverSafeEncoding() // Recover safe encoding

	err = (*client.Validator).ValidateTransaction(tx) // Validate tx

	if err != nil { // Check for errors
		common.Logf("== P2P == error while validating given tx read from pub_tx stream: %s", err.Error()) // Log error

		return // Return
	}

	senderChain, err := types.ReadChainFromMemory(*tx.Sender) // Read chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading sender chain from pub_tx stream: %s", err.Error()) // Log error

		return // Return
	}

	err = senderChain.AddTransaction(tx) // Add transaction

	if err != nil { // Check for errors
		common.Logf("== P2P == error while adding tx from pub_tx stream to sender chain: %s", err.Error()) // Log error

		return // Return
	}

	chain, err := types.ReadChainFromMemory(*tx.Recipient) // Read chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading recipient chain from pub_tx stream: %s", err.Error()) // Log error

		return // Return
	}

	err = chain.AddTransaction(tx) // Add transaction

	if err != nil { // Check for errors
		common.Logf("== P2P == error while adding tx from pub_tx stream to recipient chain: %s", err.Error()) // Log error

		return // Return
	}
}

// HandleReceiveAllChainsRequest handles an incoming req_all_chains stream.
func (client *Client) HandleReceiveAllChainsRequest(stream inet.Stream) {
	writer := bufio.NewWriter(stream) // Initialize writer

	allLocalChains, err := types.GetAllLocalizedChains() // Get all localized chains

	if err != nil { // Check for errors
		common.Logf("== P2P == error while fetching local chains tx from pub_tx stream to recipient chain: %s", err.Error()) // Log error

		return // Return
	}

	writer.Write([]byte(strings.Join(allLocalChains, "_"))) // Write all local chains
}

// HandleReceiveChainRequest handles an incoming req_chain stream.
func (client *Client) HandleReceiveChainRequest(stream inet.Stream) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer

	addressBytes, err := readWriter.ReadBytes('\f') // Read up to delimiter

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading req_chain stream: %s", err.Error()) // Log error
	}

	var address common.Address // Init buffer

	copy(address[:], addressBytes) // Write to buffer

	chain, err := types.ReadChainFromMemory(address) // Read chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading req_chain stream: %s", err.Error()) // Log error
	}

	_, err = readWriter.Write(append(chain.Bytes(), '\f')) // Write chain bytes

	if err != nil { // Check for errors
		common.Logf("== P2P == error while writing req_chain stream: %s", err.Error()) // Log error
	}

	readWriter.Flush() // Flush writer
}

/* END EXPORTED METHODS */
