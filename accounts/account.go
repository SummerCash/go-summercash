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

// AccountFromKey - generate account from given private key
func AccountFromKey(privateKey *ecdsa.PrivateKey) (*Account, error) {
	address, err := common.NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	account := Account{ // Init account from creds
		Address:    address,
		PrivateKey: privateKey,
	}

	return &account, nil // Return found account
}
