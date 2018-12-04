package types

import "github.com/space55/summertech-blockchain/common"

// CoordinationChain - "master" chain holding metadata regarding all address-spaces
type CoordinationChain struct {
	Nodes []CoordinationNode `json:"nodes"` // List of coordination nodes holding metadata regarding a specific address-space (e.g. 0x000-0x123)
}

// CoordinationNode - node holding metadata regarding a certain address-space
type CoordinationNode struct {
	AddressSpace *common.AddressSpace `json:"scope"` // Address focus
}
