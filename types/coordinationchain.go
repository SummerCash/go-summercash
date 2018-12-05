package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/crypto"
)

// CoordinationChain - "master" chain holding metadata regarding all address-spaces
type CoordinationChain struct {
	Nodes []*CoordinationNode `json:"nodes"` // List of coordination nodes holding metadata regarding a specific address-space (e.g. 0x000-0x123)

	NetworkID uint        `json:"network"` // Network ID (e.g. mainnet: 0, testnet: 1, etc...)
	ChainID   common.Hash `json:"ID"`      // Chain ID
}

// CoordinationNode - node holding metadata regarding a certain address-space
type CoordinationNode struct {
	AddressSpace *common.AddressSpace `json:"scope"`     // Address focus
	Addresses    []string             `json:"addresses"` // Node addresses in coordination node

	Origin time.Time `json:"origin"` // Time at initialization of coordination node

	ID common.Hash `json:"id"` // Node ID
}

var (
	// ErrNilAddress - error definition describing an input of addresses of length 0
	ErrNilAddress = errors.New("nil address")
)

/* BEGIN EXPORTED METHODS */

/*
	BEGIN COORDINATIONCHAIN METHODS
*/

// NewCoordinationChain - initialize new CoordinationChain
func NewCoordinationChain(networkID uint, bootstrapNode *CoordinationNode) *CoordinationChain {
	coordinationChain := &CoordinationChain{
		Nodes:     []*CoordinationNode{bootstrapNode},
		NetworkID: networkID,
	} // Init chain

	(*coordinationChain).ChainID = common.NewHash(crypto.Sha3(coordinationChain.Bytes())) // Set chain ID

	return coordinationChain // Return chain
}

// Bytes - convert given coordinationChain to byte array
func (coordinationChain *CoordinationChain) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*coordinationChain) // Serialize chain

	return buffer.Bytes() // Return serialized
}

// String - convert given coordinationChain to string
func (coordinationChain *CoordinationChain) String() string {
	marshaled, _ := json.MarshalIndent(*coordinationChain, "", "  ") // Marshal coordination chain

	return string(marshaled) // Return marshaled
}

/*
	END COORDINATIONCHAIN METHODS
*/

/*
	BEGIN COORDINATIONNODE METHODS
*/

// NewCoordinationNode - initialize new coordinationNode
func NewCoordinationNode(addressSpace *common.AddressSpace, foundingAddresses []string) (*CoordinationNode, error) {
	if len(foundingAddresses) == 0 { // Check for invalid node
		return &CoordinationNode{}, ErrNilAddress // Return error
	}

	coordinationNode := &CoordinationNode{ // Init node
		AddressSpace: addressSpace,
		Addresses:    foundingAddresses,
		Origin:       time.Now().UTC(),
	}

	coordinationNode.ID = common.NewHash(crypto.Sha3(coordinationNode.Bytes())) // Set ID

	return coordinationNode, nil // Return initialized node
}

// Bytes - convert given coordinationNode to byte array
func (coordinationNode *CoordinationNode) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*coordinationNode) // Serialize node

	return buffer.Bytes() // Return serialized
}

// String - convert given coordinationNode to string
func (coordinationNode *CoordinationNode) String() string {
	marshaled, _ := json.MarshalIndent(*coordinationNode, "", "  ") // Marshal coordination node

	return string(marshaled) // Return marshaled
}

/*
	END COORDINATIONNODE METHODS
*/

/* END EXPORTED METHODS */
