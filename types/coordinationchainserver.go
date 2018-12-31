package types

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	commonGoP2P "github.com/dowlandaiello/GoP2P/common"
	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/config"
)

// HandleReceivedCoordinationNode - handle received node
func HandleReceivedCoordinationNode(b []byte) error {
	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain

	if err != nil { // Check for errors
		return err // Return found error
	}

	coordinationNode, err := CoordinationNodeFromBytes(b) // Convert to coordinationNode

	if err != nil { // Check for errors
		return err // Return found error
	}

	node, err := coordinationChain.QueryAddress(coordinationNode.Address) // Check node already exists

	if err == nil { // Check already exists
		(*node).Addresses = append((*node).Addresses, coordinationNode.Addresses[len(coordinationNode.Addresses)-1]) // Append node

		return coordinationChain.WriteToMemory() // Write coordinationChain to memory
	}

	ip, err := common.GetExtIPAddrWithoutUPnP() // Get IP

	if err != nil { // Check for errors
		return err // Return error
	}

	if !commonGoP2P.StringInSlice(node.Addresses, ip) { // Check is not in node
		(*node).Addresses = append((*node).Addresses, ip) // Append current IP

		err = coordinationChain.AddNode(coordinationNode, true) // Add node

		if err != nil { // Check for errors
			return err // Return found error
		}
	} else {
		err = coordinationChain.AddNode(coordinationNode, false) // Add node

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	err = coordinationChain.WriteToMemory() // Write coordinationChain to memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// HandleReceivedCoordinationChainRequest - handle received byte value for coordination chain request
func HandleReceivedCoordinationChainRequest() ([]byte, error) {
	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain

	if err != nil { // Check for errors
		coordinationChain, err = NewCoordinationChain() // Init coordination chain

		if err != nil { // Check for errors
			return nil, err // Return found error
		}

		err = coordinationChain.WriteToMemory() // Write to persistent memory

		if err != nil { // Check for errors
			return nil, err // Return found error
		}

		chainConfig, err := config.ReadChainConfigFromMemory() // Read chain config

		if err != nil { // Check for errors
			return nil, err // Return found error
		}

		for _, address := range chainConfig.AllocAddresses { // Iterate through alloc
			_, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/keystore/account_%s.json", common.DataDir, address.String()))) // Read account file

			if err == nil { // Check for errors
				_, err := NewChain(address) // Init chain

				if err != nil { // Check for err`ors
					return nil, err // Return found error
				}
			}
		}

		coordinationChain, err = ReadCoordinationChainFromMemory() // Sync

		if err != nil { // Check for errors
			return nil, err // Return found error
		}
	}

	byteVal := coordinationChain.Bytes() // Get byte val

	return byteVal, nil // Return found byte value
}
