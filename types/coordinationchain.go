package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	gop2pCommon "github.com/dowlandaiello/GoP2P/common"
	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/config"
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
	Address   common.Address `json:"address"`   // Address
	Addresses []string       `json:"addresses"` // Node addresses in coordination node

	Origin time.Time `json:"origin"` // Time at initialization of coordination node

	Genesis bool `json:"genesis"` // Has genesis

	ID common.Hash `json:"id"` // Node ID
}

var (
	// ErrNilAddress - error definition describing an input of addresses of length 0
	ErrNilAddress = errors.New("nil address")

	// ErrNilNode - error definition describing a coordinationNode input of nil value
	ErrNilNode = errors.New("nil node")
)

/* BEGIN EXPORTED METHODS */

/*
	BEGIN COORDINATIONCHAIN METHODS
*/

// NewCoordinationChain - initialize new CoordinationChain
func NewCoordinationChain() (*CoordinationChain, error) {
	config, err := config.ReadChainConfigFromMemory() // Read config from memory

	if err != nil { // Check for errors
		return &CoordinationChain{}, err // Return error
	}

	coordinationChain := &CoordinationChain{ // Init chain
		Nodes:     []*CoordinationNode{},
		NetworkID: config.NetworkID,
		ChainID:   config.ChainID,
	}

	return coordinationChain, nil // Return chain
}

// AddNode - append given coordination node to coordinationChain
func (coordinationChain *CoordinationChain) AddNode(coordinationNode *CoordinationNode, updateRemote bool) error {
	if coordinationNode == nil { // Check for errors
		return ErrNilNode // Return error
	}

	if len(coordinationChain.Nodes) == 0 { // Check genesis
		(*coordinationChain).Nodes = []*CoordinationNode{coordinationNode} // Initialize node list

		return coordinationChain.WriteToMemory() // No error occurred, return nil
	}

	(*coordinationChain).Nodes = append((*coordinationChain).Nodes, coordinationNode) // Append node

	if updateRemote {
		err := coordinationChain.PushNode(coordinationNode) // Push to remote chains

		if err != nil { // Check for errors
			return err // Return error
		}
	}

	err := coordinationChain.WriteToMemory() // Save for persistency

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// QueryAddress - query for address in coordination chain
func (coordinationChain *CoordinationChain) QueryAddress(queryAddress common.Address) (*CoordinationNode, error) {
	if coordinationChain.Nodes == nil { // Check for nil nodes
		return &CoordinationNode{}, ErrNilNode // Return error
	}

	for _, node := range coordinationChain.Nodes { // Iterate through nodes
		if node != nil { // Ensure safe pointer
			if node.Address == queryAddress { // Check for match
				return node, nil // Return result
			}
		}
	}

	return &CoordinationNode{}, ErrNilNode // Return error
}

// PushNode - send new node to addresses in coordination chain
func (coordinationChain *CoordinationChain) PushNode(coordinationNode *CoordinationNode) error {
	localIP, err := common.GetExtIPAddrWithoutUPnP() // Get IP address

	if err != nil { // Check for errors
		return err // Return error
	}

	for _, node := range coordinationChain.Nodes { // Iterate through nodes
		if node != coordinationNode { // Plz no recursion
			for _, address := range node.Addresses { // Iterate through node addresses
				if address != localIP { // Plz, plz no recursion
					go common.SendBytes(coordinationNode.Bytes(), address+":"+strconv.Itoa(common.DefaultNodePort)) // Send new node
				}
			}
		}
	}

	return nil // No error occurred, return nil
}

// GetGenesis - iterate through coordination nodes, return genesis node
func (coordinationChain *CoordinationChain) GetGenesis() (*CoordinationNode, error) {
	for _, node := range coordinationChain.Nodes { // Iterate through nodes
		if node.Genesis == true { // Check genesis
			return node, nil // No error occurred, return nil
		}
	}

	return &CoordinationNode{}, ErrNilNode // Couldn't find node, return error
}

// GetBalance - attempt to get balance of account
func (coordinationChain *CoordinationChain) GetBalance(address common.Address) (float64, error) {
	node, err := coordinationChain.QueryAddress(address) // Get node

	if err != nil { // Check for errors
		return 0, err // Return found error
	}

	result, err := gop2pCommon.SendBytesResult(append([]byte("chainRequest")[:], node.Address[:]...), node.Addresses[len(node.Addresses)-1]+":"+strconv.Itoa(common.DefaultNodePort)) // Get chain

	if err != nil { // Check for errors
		return 0, err // Return found error
	}

	chain, err := FromBytes(result) // Get chain from bytes

	if err != nil { // Check for errors
		return 0, err // Return found error
	}

	return chain.CalculateBalance(), nil // No error occurred, return balance
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
func NewCoordinationNode(address common.Address, foundingAddresses []string) (*CoordinationNode, error) {
	if len(foundingAddresses) == 0 { // Check for invalid node
		return &CoordinationNode{}, ErrNilAddress // Return error
	}

	coordinationNode := &CoordinationNode{ // Init node
		Address:   address,
		Addresses: foundingAddresses,
		Origin:    time.Now().UTC(),
	}

	coordinationNode.ID = common.NewHash(crypto.Sha3(coordinationNode.Bytes())) // Set ID

	return coordinationNode, nil // Return initialized node
}

// CoordinationNodeFromBytes - convert byte array to coordinationNode
func CoordinationNodeFromBytes(b []byte) (*CoordinationNode, error) {
	coordinationNode := CoordinationNode{} // Init buffer

	err := json.NewDecoder(bytes.NewReader(b)).Decode(&coordinationNode) // Decode into buffer

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &coordinationNode, nil // No error occurred, return read value
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
