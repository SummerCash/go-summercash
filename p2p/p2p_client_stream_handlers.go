// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/crypto"

	inet "github.com/libp2p/go-libp2p-net"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
)

/* BEGIN EXPORTED METHODS */

// HandleReceiveConfigRequest handles an incoming req_config stream.
func (client *Client) HandleReceiveConfigRequest(stream inet.Stream) {
	common.Logf("== P2P == handling req_config stream\n") // Log handle stream

	writer := bufio.NewWriter(stream) // Initialize writer

	config, _ := config.ReadChainConfigFromMemory() // Read config from memory

	writer.Write(config.Bytes()) // Write config bytes

	writer.Flush() // Flush
}

// HandleReceiveTransaction handles an incoming pub_tx stream.
func (client *Client) HandleReceiveTransaction(stream inet.Stream) {
	common.Logf("== P2P == handling pub_tx stream\n") // Log handle stream

	scanner := bufio.NewScanner(stream) // Initialize scanner

	var b []byte // Init response buffer

	for scanner.Scan() { // Scan
		b = append(b, scanner.Bytes()...) // Append scanned
	}

	err := scanner.Err() // Get error

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading pub_tx stream: %s\n", err.Error()) // Log error

		return // Return
	}

	tx, err := types.TransactionFromBytes(b) // Marshal bytes to transaction

	if err != nil { // Check for errors
		common.Logf("== P2P == error while deserializing tx read from pub_tx stream: %s\n", err.Error()) // Log error

		return // Return
	}

	tx.RecoverSafeEncoding() // Recover safe encoding

	err = (*client.Validator).ValidateTransaction(tx) // Validate tx

	if err != nil { // Check for errors
		common.Logf("== P2P == error while validating given tx read from pub_tx stream: %s\n", err.Error()) // Log error

		return // Return
	}

	senderChain, err := types.ReadChainFromMemory(*tx.Sender) // Read chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading sender chain from pub_tx stream: %s\n", err.Error()) // Log error

		return // Return
	}

	err = senderChain.AddTransaction(tx) // Add transaction

	if err != nil { // Check for errors
		common.Logf("== P2P == error while adding tx from pub_tx stream to sender chain: %s\n", err.Error()) // Log error

		return // Return
	}

	chain, err := types.ReadChainFromMemory(*tx.Recipient) // Read chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading recipient chain from pub_tx stream: %s\n", err.Error()) // Log error

		return // Return
	}

	err = chain.AddTransaction(tx) // Add transaction

	if err != nil { // Check for errors
		common.Logf("== P2P == error while adding tx from pub_tx stream to recipient chain: %s\n", err.Error()) // Log error

		return // Return
	}
}

// HandleReceiveBestTransaction handles an incoming req_best_tx stream.
func (client *Client) HandleReceiveBestTransaction(stream inet.Stream) {
	common.Logf("== P2P == handling req_best_tx stream\n") // Log handle stream

	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer

	scanner := bufio.NewScanner(stream) // Init scanner

	var accountString []byte // Init response buffer

	for scanner.Scan() { // Scan
		accountString = append(accountString, scanner.Bytes()...) // Append scanned
	}

	err := scanner.Err() // Get error

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading req_best_tx stream: %s\n", err.Error()) // Log error

		return // Return
	}

	address, err := common.StringToAddress(string(accountString)) // Get address

	if err != nil { // Check for errors
		common.Logf("== P2P == error while parsing req_best_tx stream: %s\n", err.Error()) // Log error
	}

	chain, err := types.ReadChainFromMemory(address) // Read chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading chain from req_best_tx stream: %s\n", err.Error()) // Log error
	}

	if len(chain.Transactions) > 0 { // Check has txs
		readWriter.Write(chain.Transactions[len(chain.Transactions)-1].Hash.Bytes()) // Write tx hash
	} else { // No txs
		readWriter.Write(common.NewHash(crypto.Sha3(nil)).Bytes()) // Write nil hash
	}

	readWriter.Flush() // Flush
}

// HandleReceiveNextTransactionRequest handles an incoming req_next_tx stream.
func (client *Client) HandleReceiveNextTransactionRequest(stream inet.Stream) {
	common.Logf("== P2P == handling req_next_tx stream\n") // Log handle stream

	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer

	scanner := bufio.NewScanner(stream) // Init scanner

	var lastTxAccount []byte // Init response buffer

	for scanner.Scan() { // Scan
		lastTxAccount = append(lastTxAccount, scanner.Bytes()...) // Append scanned
	}

	err := scanner.Err() // Get error

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading req_next_tx stream: %s\n", err.Error()) // Log error

		return // Return
	}

	address, err := common.StringToAddress(strings.Split(string(lastTxAccount), "_")[0]) // Get address

	if err != nil { // Check for errors
		common.Logf("== P2P == error while parsing req_next_tx stream: %s\n", err.Error()) // Log error
	}

	hash, err := common.StringToHash(strings.Split(string(lastTxAccount), "_")[1]) // Get hash

	if err != nil { // Check for errors
		common.Logf("== P2P == error while parsing req_next_tx stream: %s\n", err.Error()) // Log error
	}

	accountChain, err := types.ReadChainFromMemory(address) // Read account chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error reading req_next_tx stream: %s\n", err.Error()) // Log error
	}

	if bytes.Equal(hash.Bytes(), common.NewHash(crypto.Sha3(nil)).Bytes()) { // Check is nil request
		readWriter.Write(accountChain.Transactions[0].Bytes()) // Write genesis bytes

		readWriter.Flush() // Flush

		return // Return
	}

	for x, transaction := range accountChain.Transactions { // Iterate through transactions
		if bytes.Equal(transaction.Hash.Bytes(), hash.Bytes()) { // Check hashes equal
			if len(accountChain.Transactions) == x+1 { // Check no next
				readWriter.Write(accountChain.Transactions[x].Bytes()) // Write current transaction

				readWriter.Flush() // Flush

				break // Break
			}

			readWriter.Write(accountChain.Transactions[x+1].Bytes()) // Write next transaction

			readWriter.Flush() // Flush

			break // Break
		}
	}
}

// HandleReceiveAllChainsRequest handles an incoming req_all_chains stream.
func (client *Client) HandleReceiveAllChainsRequest(stream inet.Stream) {
	common.Logf("== P2P == handling req_all_chains stream\n") // Log handle stream

	writer := bufio.NewWriter(stream) // Initialize writer

	allLocalChains, err := types.GetAllLocalizedChains() // Get all localized chains

	if err != nil { // Check for errors
		common.Logf("== P2P == error while fetching local chains tx from pub_tx stream to recipient chain: %s\n", err.Error()) // Log error
	}

	_, err = writer.Write([]byte(strings.Join(allLocalChains, "_"))) // Write all local chains

	if err != nil { // Check for errors
		common.Logf("== P2P == error while writing req_chain stream: %s\n", err.Error()) // Log error
	}

	writer.Flush() // Flush
}

// HandleReceiveChainRequest handles an incoming req_chain stream.
func (client *Client) HandleReceiveChainRequest(stream inet.Stream) {
	common.Logf("== P2P == handling req_chain stream\n") // Log handle stream

	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer

	scanner := bufio.NewScanner(stream) // Init scanner

	var addressBytes []byte // Init response buffer

	for scanner.Scan() { // Scan
		addressBytes = append(addressBytes, scanner.Bytes()...) // Append scanned
	}

	err := scanner.Err() // Get error

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading req_chain stream: %s\n", err.Error()) // Log error

		return // Return
	}

	var address common.Address // Init buffer

	copy(address[:], addressBytes) // Write to buffer

	chain, err := types.ReadChainFromMemory(address) // Read chain

	if err != nil { // Check for errors
		common.Logf("== P2P == error while reading req_chain stream: %s\n", err.Error()) // Log error
	}

	_, err = readWriter.Write(chain.Bytes()) // Write chain bytes

	if err != nil { // Check for errors
		common.Logf("== P2P == error while writing req_chain stream: %s\n", err.Error()) // Log error
	}

	readWriter.Flush() // Flush writer
}

// HandleReceiveAliveRequest handles an incoming req_not_dead_lol stream.
func (client *Client) HandleReceiveAliveRequest(stream inet.Stream) {
	common.Logf("== P2P == handling req_not_dead_lol stream\n") // Log handle stream

	writer := bufio.NewWriter(stream) // Init writer

	_, err := writer.Write([]byte("despacito")) // Write alive

	if err != nil { // Check for errors
		common.Logf("== P2P == error while writing req_not_dead_lol stream: %s\n", err.Error()) // Log error
	}

	writer.Flush() // Flush writer
}

/* END EXPORTED METHODS */
