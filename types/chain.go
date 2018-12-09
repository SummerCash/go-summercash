package types

import "github.com/space55/summertech-blockchain/common"

// Chain - account transactions chain
type Chain struct {
	Account common.Address `json:"account"` // Chain account

	Transactions []*Transaction `json:"transactions"` // Transactions in chain

	NetworkID uint        `json:"network"` // Network ID (mainnet: 0, testnet: 1, etc...)
	ID        common.Hash `json:"ID"`      // Chain ID
}
