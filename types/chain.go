package types

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/config"
	"github.com/space55/summertech-blockchain/crypto"
)

// Chain - account transactions chain
type Chain struct {
	Account common.Address `json:"account"` // Chain account

	Transactions []*Transaction `json:"transactions"` // Transactions in chain

	Genesis common.Hash `json:"genesis"` // Genesis block hash

	NetworkID uint        `json:"network"` // Network ID (mainnet: 0, testnet: 1, etc...)
	ID        common.Hash `json:"ID"`      // Chain ID
}

var (
	// ErrChainAlreadyExists - error definition describing a given chain that has already been registered in the coordinationChain
	ErrChainAlreadyExists = errors.New("chain already exists for given account")

	// ErrGenesisAlreadyExists - error definition describing a given chain with an existing genesis block
	ErrGenesisAlreadyExists = errors.New("chain already has existing genesis")
)

/* BEGIN EXPORTED METHODS */

// NewChain - initialize new chain
func NewChain(account common.Address) (*Chain, error) {
	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain from memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	config, err := config.ReadChainConfigFromMemory() // Read config from memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	_, err = coordinationChain.QueryAddress(account) // Query address

	if err != nil && err != ErrNilNode { // Check chain with address does not exist
		return &Chain{}, err // Return error
	} else if err == nil { // Check exists
		return &Chain{}, ErrAlreadySigned // Return error
	}

	if coordinationChain.Nodes == nil || len(coordinationChain.Nodes) == 0 { // Check genesis

	}

	chain := &Chain{ // Init chain
		Account:      account,
		Transactions: []*Transaction{},
		NetworkID:    config.NetworkID,
	}

	(*chain).ID = common.NewHash(crypto.Sha3(chain.Bytes())) // Set ID

	return chain, nil // Return initialized chain
}

// Bytes - convert given chain to byte array
func (chain *Chain) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*chain) // Serialize chain

	return buffer.Bytes() // Return serialized
}

// String - convert given chain to string
func (chain *Chain) String() string {
	marshaled, _ := json.MarshalIndent(*chain, "", "  ") // Marshal chain

	return string(marshaled) // Return marshaled
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// makeGenesis - generate genesis blocks from genesis file
func (chain *Chain) makeGenesis(genesis *config.ChainConfig) (common.Hash, error) {
	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain from memory

	if err != nil { // Check for errors
		return common.Hash{}, err // Return error
	}

	genesisNode, _ := coordinationChain.GetGenesis() // Get genesis

	if genesisNode != nil { // Check genesis already exists
		return common.Hash{}, ErrGenesisAlreadyExists // Return error
	}

	genesisTx, err := NewTransaction(0, nil, nil, &genesis.AllocAddresses[0], genesis.Alloc[genesis.AllocAddresses[0].String()], []byte("genesis")) // Init transaction

	if err != nil { // Check for errors
		return common.Hash{}, err // Return error
	}

	return *genesisTx.Hash, nil // Return genesis
}

/* END INTERNAL METHODS */
