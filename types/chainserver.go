package types

import (
	"github.com/space55/summertech-blockchain/common"
)

/* BEGIN EXPORTED METHODS */

// HandleReceivedChainRequest - handle chain request
func HandleReceivedChainRequest(b []byte) (*Chain, error) {
	var address common.Address // Init buffer

	copy(address[:], b[12:len(b)][:]) // Copy read address

	chain, err := ReadChainFromMemory(address) // Read chain from memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	return chain, nil // Return read chain
}

/* END EXPORTED METHODS */
