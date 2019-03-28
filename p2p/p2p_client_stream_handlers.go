// Package p2p outlines helper methods and types for p2p communications.
package p2p

import inet "github.com/libp2p/go-libp2p-net"

/* BEGIN EXPORTED METHODS */

// HandleReceiveTransaction handles an incoming pub_tx stream.
func (client *Client) HandleReceiveTransaction(inet.Stream) {
	tx, err := TransactionFromBytes(b) // Marshal bytes to transaction

	if err != nil { // Check for errors
		return err // Return error
	}

	tx.RecoverSafeEncoding() // Recover safe encoding

	chain, err := ReadChainFromMemory(*tx.Recipient) // Read chain

	if err != nil { // Check for errors
		return err // Return found error
	}

	if tx.AccountNonce != uint64(len(chain.Transactions)) { // Check invalid nonce
		return
	}
}

/* END EXPORTED METHODS */
