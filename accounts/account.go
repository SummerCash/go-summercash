package accounts

import (
	"crypto/ecdsa"

	"github.com/space55/summertech-blockchain/common"
)

// Account - container holding account metadata, private key
type Account struct {
	Address    common.Address    `json:"address"`    // Account address
	PrivateKey *ecdsa.PrivateKey `json:"privateKey"` // Account private key
}
