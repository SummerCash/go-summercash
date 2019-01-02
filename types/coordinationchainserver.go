package types

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	commonGoP2P "github.com/dowlandaiello/GoP2P/common"
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
		if !commonGoP2P.StringInSlice(node.Addresses, coordinationNode.Addresses[len(coordinationNode.Addresses)-1]) { // Check not already in node
			(*node).Addresses = append((*node).Addresses, coordinationNode.Addresses[len(coordinationNode.Addresses)-1]) // Append node

			return coordinationChain.WriteToMemory() // Write coordinationChain to memory
		}
	}

	ip, err := common.GetExtIPAddrWithoutUPnP() // Get IP

	if err != nil { // Check for errors
		return err // Return error
	}

	ipPortIncluded := "" // Init buffer

	if strings.Contains(ip, ":") { // Check is IPv6
		ipPortIncluded = "[" + ip + "]" + ":" + strconv.Itoa(common.NodePort) // Add port
	} else {
		ipPortIncluded = ip + ":" + strconv.Itoa(common.NodePort) // Add port
	}

	if !commonGoP2P.StringInSlice(coordinationNode.Addresses, ipPortIncluded) { // Check is not in node
		common.Logf("== NETWORK == adding self %s to coordination node %s\n", ipPortIncluded, coordinationNode.Address.String()) // Log add self

		(*coordinationNode).Addresses = append((*coordinationNode).Addresses, ipPortIncluded) // Append current IP

		common.Logf("== CHAIN == node addresses %s\n", (*coordinationNode).Addresses) // Log addrs

		common.Logf("== NETWORK == pushing coordination node %s\n", coordinationNode.Address.String()) // Log push

		err = coordinationChain.AddNode(coordinationNode, true) // Add node

		if err != nil { // Check for errors
			common.Logf("== ERROR == error pushing coordination node %s\n", err.Error()) // Log error

			return err // Return found error
		}

		common.Logf("== SUCCESS == successfully pushed coordination node %s to network\n", coordinationNode.Address.String()) // Log success
	} else {
		common.Logf("== NETWORK == added coordination node %s to local coordination chain %s\n", coordinationNode.Address.String(), coordinationChain.ChainID.String()) // Log add to local chain

		err = coordinationChain.AddNode(coordinationNode, false) // Add node

		if err != nil { // Check for errors
			common.Logf("== ERROR == error adding coordination node to local coordination chain %s\n", err.Error()) // Log error

			return err // Return found error
		}

		common.Logf("== SUCCESS == successfully pushed coordination node %s to local coordination chain\n", coordinationNode.Address.String()) // Log success
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
