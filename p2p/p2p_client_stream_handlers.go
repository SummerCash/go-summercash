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

// HandleReceiveBestTransaction handles an incoming req_best_tx stream.
func (client *Client) HandleReceiveBestTransaction(stream inet.Stream) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer

	accountString, err := readWriter.ReadBytes('\f') // Read

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading req_best_tx stream: %s", err.Error()) // Log error
	}

	address, err := common.StringToAddress(string(accountString)) // Get address

	if err != nil { // Check for errors
		common.Logf("== P2P == error while parsing req_best_tx stream: %s", err.Error()) // Log error
	}

	chain, err := types.ReadChainFromMemory(address) // Read chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading chain from req_best_tx stream: %s", err.Error()) // Log error
	}

	readWriter.Write(append(chain.Bytes(), '\f')) // Write chain bytes

	readWriter.Flush() // Flush
}

// HandleReceiveNextTransactionRequest handles an incoming req_next_tx stream.
func (client *Client) HandleReceiveNextTransactionRequest(stream inet.Stream) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer

	lastTxAccount, err := readWriter.ReadBytes('\f') // Read

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading req_next_tx stream: %s", err.Error()) // Log error
	}

	address, err := common.StringToAddress(strings.Split(string(lastTxAccount), "_")[0]) // Get address

	if err != nil { // Check for errors
		common.Logf("== P2P == error while parsing req_next_tx stream: %s", err.Error()) // Log error
	}

	hash, err := common.StringToHash(strings.Split(string(lastTxAccount), "_")[1]) // Get hash

	if err != nil { // Check for errors
		common.Logf("== P2P == error while parsing req_next_tx stream: %s", err.Error()) // Log error
	}

	accountChain, err := types.ReadChainFromMemory(address) // Read account chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error reading req_next_tx stream: %s", err.Error()) // Log error
	}

	for x, transaction := range accountChain.Transactions { // Iterate through transactions
		if bytes.Equal(transaction.Hash.Bytes(), hash.Bytes()) { // Check hashes equal
			readWriter.Write(append(accountChain.Transactions[x+1].Bytes(), '\f')) // Write next transaction

			readWriter.Flush() // Flush

			break // Break
		}
	}
}

// HandleReceiveAllChainsRequest handles an incoming req_all_chains stream.
func (client *Client) HandleReceiveAllChainsRequest(stream inet.Stream) {
	writer := bufio.NewWriter(stream) // Initialize writer

	allLocalChains, err := types.GetAllLocalizedChains() // Get all localized chains

	if err != nil { // Check for errors
		common.Logf("== P2P == error while fetching local chains tx from pub_tx stream to recipient chain: %s", err.Error()) // Log error
	}

	_, err = writer.Write(append([]byte(strings.Join(allLocalChains, "_")), '\f')) // Write all local chains

	if err != nil { // Check for errors
		common.Logf("== P2P == error while writing req_chain stream: %s", err.Error()) // Log error
	}

	writer.Flush() // Flush
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
