package types

import "github.com/space55/summertech-blockchain/common"

// Chain - account blockchain
type Chain struct {
	Account common.Address `json:"account"` // Chain account

	Transactions []*Transaction `json:"transactions"` // Transactions in chain

	ID common.Hash `json:"ID"` // Chain ID
}
