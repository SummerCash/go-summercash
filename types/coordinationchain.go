package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/crypto"
	gop2pCommon "github.com/dowlandaiello/GoP2P/common"
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

	// ErrNilCoordinationChain - error definition describing a coordination chain that is nil in value
	ErrNilCoordinationChain = errors.New("nil coordination chain")

	// ErrClientOutOfDate - error definition describing a client that is out of date
	ErrClientOutOfDate = errors.New("client out of date (must upgrade client)")

	// ErrNoBootstrapNodes - error definition describing a network having no bootstrap nodes
	ErrNoBootstrapNodes = errors.New("== WARNING == no available bootstrap nodes")

	// ErrNilCoordinationChainCache - error definition describing a coordination chain that is nil in value, but must have its cache cleared
	ErrNilCoordinationChainCache = errors.New("nil coordination chain; nothing to clear")

	// ErrNilCoordinationChainPush - error definition describing a coordination chain that is nil in value, but must be pushed
	ErrNilCoordinationChainPush = errors.New("nil coordination chain; nothing to push")
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
		(*coordinationNode).Genesis = true                                 // Set genesis
		(*coordinationChain).Nodes = []*CoordinationNode{coordinationNode} // Initialize node list

		return coordinationChain.WriteToMemory() // No error occurred, return nil
	}

	node, err := coordinationChain.QueryAddress(coordinationNode.Address) // Check node already exists

	if err != nil { // Check for errors
		(*coordinationChain).Nodes = append((*coordinationChain).Nodes, coordinationNode) // Append node
	} else { // Node already exists
		(*node).Addresses = append((*node).Addresses, coordinationNode.Addresses[len(coordinationNode.Addresses)-1]) // Append node
	}

	if updateRemote { // Check should update remote db
		err := coordinationChain.PushNode(coordinationNode) // Push to remote chains

		if err != nil { // Check for errors
			return err // Return error
		}
	}

	err = coordinationChain.WriteToMemory() // Save for persistency

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// ClearCache - remove unresponsive nodes from given coordinationChain
func (coordinationChain *CoordinationChain) ClearCache() error {
	if coordinationChain.Nodes == nil || len(coordinationChain.Nodes) == 0 { // Check nothing to clear
		return ErrNilCoordinationChain // Return error
	}

	for _, node := range coordinationChain.Nodes { // Iterate through nodes
		verifiedNodes := []string{} // Init buffer

		for _, address := range node.Addresses { // Iterate through providing addresses
			if !gop2pCommon.StringInSlice(verifiedNodes, address) { // Check must be tested
				coordinationChainBytes, err := gop2pCommon.SendBytesResult([]byte("cChainRequest"), address) // Get coordination chain

				if err == nil && coordinationChainBytes != nil { // Check no errors
					verifiedNodes = append(verifiedNodes, address) // Append verified node address
				}
			}
		}

		(*node).Addresses = verifiedNodes // Set to cleared
	}

	err := coordinationChain.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = coordinationChain.UpdateRemotes() // Update remote instances

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// UpdateRemotes - update remote coordinationChain instances
func (coordinationChain *CoordinationChain) UpdateRemotes() error {
	if coordinationChain.Nodes == nil || len(coordinationChain.Nodes) == 0 { // Check for errors
		return ErrNilCoordinationChainPush // Return error
	}

	pushedNodes := []string{} // Init buffer

	for _, node := range coordinationChain.Nodes { // Iterate through nodes
		for _, address := range node.Addresses { // Iterate through providing addresses
			if !gop2pCommon.StringInSlice(pushedNodes, address) { // Check must be tested
				err := common.SendBytes(coordinationChain.Bytes(), address) // Push

				if err == nil { // Check no errors
					pushedNodes = append(pushedNodes, address) // Append verified node address
				}
			}
		}
	}

	return nil // No error occurred, return nil
}

// JoinNetwork - join given network with bootstrap node address
func JoinNetwork(bootstrapNode string, archivalNode bool) error {
	common.Logf("== NETWORK == requesting coordination chain from bootstrap node %s\n", bootstrapNode) // Log init

	coordinationChainBytes, err := gop2pCommon.SendBytesResult([]byte("cChainRequest"), bootstrapNode) // Get coordination chain

	if err != nil { // Check for errors
		return err // Return found error
	}

	coordinationChain, err := CoordinationChainFromBytes(coordinationChainBytes) // Decode result

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = coordinationChain.WriteToMemory() // Write to persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	common.Logf("== SUCCESS == received coordination chain %s from bootstrap node %s\n", coordinationChain.ChainID.String(), bootstrapNode) // Log success
	common.Logf("== NETWORK == requesting chain config from bootstrap node %s\n", bootstrapNode)                                            // Log request config

	configBytes, err := gop2pCommon.SendBytesResult([]byte("configReq"), bootstrapNode) // Get chain config

	if err != nil { // Check for errors
		return err // Return found error
	}

	config, err := config.FromBytes(configBytes) // Decode config

	if err != nil { // Check for errors
		return err // Return found error
	}

	common.Logf("== SUCCESS == received chain config with network ID %d from bootstrap node %s\n", config.NetworkID, bootstrapNode) // Log success

	err = config.WriteToMemory() // Write config to persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	if archivalNode { // Check is registering archival node
		err = RegisterArchivalNode() // Register archival node

		if err != nil { // Check for errors
			return err // Return found error
		}

		return SyncNetwork(archivalNode, false) // Register archival node
	}

	return nil // No error occurred, return nil
}

// SyncNetwork - download all chains
func SyncNetwork(archival bool, updateRemote bool) error {
	var coordinationChainBytes []byte // Init buffer
	var err error                     // Init error buffer

	ip, err := common.GetExtIPAddrWithoutUPnP() // Get IP

	if err != nil { // Check for errors
		return err // Return found error
	}

	coordinationChain, _ := ReadCoordinationChainFromMemory() // Read coordination chain

	if strings.Contains(ip, ":") { // Check is IPv6
		ip = "[" + ip + "]" + ":" + strconv.Itoa(common.NodePort) // Add port
	} else {
		ip = ip + ":" + strconv.Itoa(common.NodePort) // Add port
	}

	x := 0 // Init iterator

	for {
		if x >= len(common.BootstrapNodes) { // Check is not out of bounds
			if coordinationChain == nil { // Check nodes available
				return ErrNoBootstrapNodes // Panic
			}

			archivalNodes, err := coordinationChain.QueryAllArchivalNodes() // Query all archival nodes

			if err != nil { // Check for errors
				return err // Return found error
			}

			common.BootstrapNodes = append(common.BootstrapNodes, archivalNodes...) // Append archival nodes
		}

		common.Logf("== NETWORK == requesting coordination chain from bootstrap node %s\n", common.BootstrapNodes[x]) // Log req

		if common.BootstrapNodes[x] != ip { // Prevent recursion
			coordinationChainBytes, err = gop2pCommon.SendBytesResult([]byte("cChainRequest"), common.BootstrapNodes[x]) // Get coordination chain

			if err == nil { // Check for errors
				break // Break loop
			}

			x++ // Increment
		} else {
			x++ // Increment
		}
	}

	err = syncChainConfig() // Sync chain config

	if err != nil { // Check for errors
		return err // Return found error
	}

	coordinationChain, err = CoordinationChainFromBytes(coordinationChainBytes) // Decode result

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = coordinationChain.WriteToMemory() // Write to persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	common.Logf("== NODE == syncing with network %s\n", coordinationChain.ChainID.String()) // Log sync

	if archival { // Check is archival node
		for _, node := range coordinationChain.Nodes { // Iterate through nodes
			common.Logf("== NETWORK == requesting account chain for address %s\n", node.Address.String()) // Log req

			chainBytes := []byte{} // Init buffer

			var err error // Init error buffer

			for _, address := range node.Addresses { // Iterate through node providers
				if address != ip { // Prevent recursive node lookup
					chainBytes, err = gop2pCommon.SendBytesResult(append([]byte("chainRequest")[:], node.Address[:]...), address) // Get chain

					if err == nil { // Check for errors
						break // Break
					}
				}
			}

			if err != nil { // Check for errors
				return err // Return found error
			}

			chain, err := FromBytes(chainBytes) // Get chain

			if err != nil { // Check for errors
				return err // Return found error
			}

			err = chain.WriteToMemory() // Write chain to memory

			if err != nil { // Check for errors
				return err // Return found error
			}
		}

		if updateRemote { // Check needs to re-register archival node
			err = RegisterArchivalNode() // Register archival node

			if err != nil { // Check for errors
				return err // Return found error
			}
		}
	}

	files, err := ioutil.ReadDir(filepath.FromSlash(fmt.Sprintf("%s/chain", common.DataDir))) // Walk keystore dir

	if err == nil { // Check no error
		for _, file := range files { // Iterate through files
			data, err := ioutil.ReadFile(file.Name()) // Read file JSON bytes

			if err != nil { // Check for errors
				return err // Return found error
			}

			chain, err := FromBytes(data) // Read chain from bytes

			if err != nil { // Check for errors
				return err // Return found error
			}

			node, err := coordinationChain.QueryAddress(chain.Account) // Query address

			if err != nil { // Check for errors
				return err // Return found error
			}

			if node.Addresses[0] != ip+":"+strconv.Itoa(common.NodePort) { // Check not current node
				common.SendBytes(data, node.Addresses[0]) // Send chain
			}

			for x, address := range node.Addresses { // Iterate through addresses
				if x != 0 && address != ip+":"+strconv.Itoa(common.NodePort) { // Check not first index and not current addr
					go common.SendBytes(data, address) // Send chain
				}
			}
		}
	}

	err = coordinationChain.ClearCache() // Clear cooridnation chain cache

	if err != nil { // Check for errors
		return err // Return found error
	}

	common.Log("== NODE == finished syncing") // Log success

	return nil // No error occurred, return nil
}

// RegisterArchivalNode - register archival node on network
func RegisterArchivalNode() error {
	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain from persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	ip, err := common.GetExtIPAddrWithoutUPnP() // Get IP

	if err != nil { // Check for errors
		return err // Return found error
	}

	common.Logf("== NODE == registering local archival node with external IP %s\n", ip) // Log register

	if strings.Contains(ip, ":") { // Check is IPv6
		ip = "[" + ip + "]" + ":" + strconv.Itoa(common.NodePort) // Add port
	} else {
		ip = ip + ":" + strconv.Itoa(common.NodePort) // Add port
	}

	_, err = coordinationChain.QueryArchivalNode(ip) // Check node already in network

	if err != nil { // Check for errors
		for _, node := range coordinationChain.Nodes { // Iterate through nodes
			nodeInstance, err := NewCoordinationNode(node.Address, []string{ip}) // Init node

			if err != nil { // Check for errors
				return err // Return found error
			}

			err = coordinationChain.AddNode(nodeInstance, true) // Add node

			if err != nil { // Check for errors
				return err // Return found error
			}
		}

		common.Log("== NETWORK == finished registering archval node") // Log success
	}

	return coordinationChain.WriteToMemory() // No error occurred, return nil
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

// QueryNode - query for node address in coordination chain
func (coordinationChain *CoordinationChain) QueryNode(address string) (*CoordinationNode, error) {
	if address == "" { // Check for nil address
		return nil, ErrNilAddress // Return error
	}

	for _, node := range coordinationChain.Nodes { // Iterate through nodes
		if node != nil { // Ensure safe pointer
			for _, currentAddress := range node.Addresses { // Iterate through addresses
				if currentAddress == address { // Check has address
					return node, nil // Found match, return node
				}
			}
		}
	}

	return &CoordinationNode{}, ErrNilNode // Return error
}

// QueryArchivalNode - query for archival node address in coordination chain
func (coordinationChain *CoordinationChain) QueryArchivalNode(address string) ([]*CoordinationNode, error) {
	if address == "" { // Check for nil address
		return nil, ErrNilAddress // Return error
	}

	matches := []*CoordinationNode{} // Init matches

	for _, node := range coordinationChain.Nodes { // Iterate through nodes
		if node != nil { // Ensure safe pointer
			for _, currentAddress := range node.Addresses { // Iterate through addresses
				if currentAddress == address { // Check has address
					if len(matches) == 0 { // Check must init
						matches = []*CoordinationNode{node} // Init matches
					}

					matches = append(matches, node) // Append found node
				}
			}
		}
	}

	if float64(len(matches)) > 0.25*float64(len(coordinationChain.Nodes)) { // Check enough matches
		return matches, nil // Return found matches
	}

	return []*CoordinationNode{}, ErrNilNode // Return error
}

// QueryAllArchivalNodes - get all archival nodes in coordination chain
func (coordinationChain *CoordinationChain) QueryAllArchivalNodes() ([]string, error) {
	if coordinationChain == nil { // Check nil pointer
		return []string{}, ErrNilCoordinationChain // Return error
	}

	matches := []string{} // Init matches

	for x, node := range coordinationChain.Nodes { // Iterate through nodes
		if node != nil { // Ensure safe pointer
			for _, currentAddress := range node.Addresses { // Iterate through addresses
				if len(coordinationChain.Nodes) > 2 {
					if len(matches) == 0 { // Check init
						if x+1 < len(coordinationChain.Nodes) && gop2pCommon.StringInSlice(coordinationChain.Nodes[x+1].Addresses, node.Addresses[0]) || len(coordinationChain.Nodes) == 0 { // Check can be + indexed
							matches = []string{node.Addresses[0]} // Init matches
						}
					} else {
						if x+1 < len(coordinationChain.Nodes) && gop2pCommon.StringInSlice(coordinationChain.Nodes[x+1].Addresses, currentAddress) || len(coordinationChain.Nodes) == 0 { // Check can be + indexed
							matches = append(matches, currentAddress) // Append to matches
						}
					}
				} else {
					if len(matches) == 0 { // Check init
						matches = []string{node.Addresses[0]} // Init matches
					} else {
						matches = append(matches, currentAddress) // Append to matches
					}
				}
			}
		}
	}

	return matches, nil // Return error
}

// PushNode - send new node to addresses in coordination chain
func (coordinationChain *CoordinationChain) PushNode(coordinationNode *CoordinationNode) error {
	common.Logf("== NETWORK == pushing coordination chain node %s to network\n", coordinationNode.Address.String()) // Log push

	localIP, err := common.GetExtIPAddrWithoutUPnP() // Get IP address

	if err != nil { // Check for errors
		return err // Return error
	}

	if strings.Contains(localIP, ":") { // Check is IPv6
		localIP = "[" + localIP + "]" + ":" + strconv.Itoa(common.NodePort) // Add port
	} else {
		localIP = localIP + ":" + strconv.Itoa(common.NodePort) // Add port
	}

	for _, node := range coordinationChain.Nodes { // Iterate through nodes
		if node != coordinationNode { // Plz no recursion
			for _, address := range node.Addresses { // Iterate through node addresses
				if address != localIP { // Plz, plz no recursion
					common.Logf("== NETWORK == pushing coordination chain node %s to peer %s\n", coordinationNode.Address.String(), address) // Log push

					go common.SendBytes(coordinationNode.Bytes(), address) // Send new node
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

	var result []byte // Init buffer

	for _, nodeAddress := range node.Addresses { // Iterate through node addresses
		result, err = gop2pCommon.SendBytesResult(append([]byte("chainRequest")[:], node.Address[:]...), nodeAddress) // Get chain

		if err == nil { // Check no errors
			break // Break
		}
	}

	chain, err := FromBytes(result) // Get chain from bytes

	if err != nil { // Check for errors
		return 0, err // Return found error
	}

	return chain.CalculateBalance(), nil // No error occurred, return balance
}

// CoordinationChainFromBytes - decode coordination chain from given byte array
func CoordinationChainFromBytes(b []byte) (*CoordinationChain, error) {
	coordinationChain := CoordinationChain{} // Init buffer

	err := json.NewDecoder(bytes.NewReader(b)).Decode(&coordinationChain) // Decode into buffer

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &coordinationChain, nil // No error occurred, return read value
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

/* BEGIN INTERNAL METHODS */

func syncChainConfig() error {
	var configBytes []byte // Init config buffer
	var err error          // Init error buffer

	for _, bootstrapNode := range common.BootstrapNodes { // Iterate through bootstrap nodes
		common.Logf("== NETWORK == requesting chain config from bootstrap node %s\n", bootstrapNode) // Log request config

		configBytes, err = gop2pCommon.SendBytesResult([]byte("configReq"), bootstrapNode) // Get chain config

		if err != nil { // Check for errors
			continue // Continue
		} else {
			break // Break
		}
	}

	remoteConfig, err := config.FromBytes(configBytes) // Get config from bytes

	if err != nil { // Check for errors
		return err // Return found error
	}

	localConfig, err := config.ReadChainConfigFromMemory() // Read local chain config

	if err != nil {
		return remoteConfig.WriteToMemory() // Write config to persistent memory
	}

	softForkRemote, _ := strconv.Atoi(strings.Split(remoteConfig.ChainVersion, ".")[1]) // Get remote soft fork version int val
	softForkLocal, _ := strconv.Atoi(strings.Split(localConfig.ChainVersion, ".")[1])   // Get local soft fork version int val

	hardForkRemote, _ := strconv.Atoi(strings.Split(remoteConfig.ChainVersion, ".")[0]) // Get remote hard fork version int val
	hardForkLocal, _ := strconv.Atoi(strings.Split(localConfig.ChainVersion, ".")[0])   // Get remote hard fork version int val

	if softForkRemote-softForkLocal > 4 && hardForkRemote == hardForkLocal || hardForkRemote > hardForkLocal { // Check too out-of-date
		return ErrClientOutOfDate // Return out-of-date error
	}

	return nil // No error occurred, return nil
}

/* END INTERNAL METHODS */
