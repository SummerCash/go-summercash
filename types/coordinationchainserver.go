package types

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

	err = coordinationChain.AddNode(coordinationNode, false) // Add node

	if err != nil { // Check for errors
		return err // Return found error
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
	}

	byteVal := coordinationChain.Bytes() // Get byte val

	return byteVal, nil // Return found byte value
}
