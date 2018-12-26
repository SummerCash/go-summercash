package accounts

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"

	"github.com/space55/summertech-blockchain/common"
)

// Account - container holding account metadata, private key
type Account struct {
	Address              common.Address    `json:"address"`      // Account address
	PrivateKey           *ecdsa.PrivateKey `json:"privateKey"`   // Account private key
	SerializedPrivateKey []byte            `json:"s_privateKey"` // Serialized account private key
}

/* BEGIN EXPORTED METHODS */

// AccountFromKey - generate account from given private key
func AccountFromKey(privateKey *ecdsa.PrivateKey) (*Account, error) {
	address, err := common.NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	account := Account{ // Init account from creds
		Address:    address,    // Initialize with address
		PrivateKey: privateKey, // Initialize with private key
	}

	return &account, nil // Return found account
}

// MakeEncodingSafe - make account safe for encoding
func (account *Account) MakeEncodingSafe() error {
	marshaledPrivateKey, err := x509.MarshalECPrivateKey(account.PrivateKey) // Marshal private key

	if err != nil { // Check for errors
		return err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshaledPrivateKey}) // Encode to memory

	(*account).SerializedPrivateKey = pemEncoded // Set serialized

	*(*account).PrivateKey = ecdsa.PrivateKey{} // Set nil

	return nil // No error occurred, return nil
}

// RecoverSafeEncoding - recover full data from safely encoded type
func (account *Account) RecoverSafeEncoding() error {
	blockPub, _ := pem.Decode([]byte(account.SerializedPrivateKey)) // Decode

	x509EncodedPrivateKey := blockPub.Bytes // Get x509 byte val

	privateKey, err := x509.ParseECPrivateKey(x509EncodedPrivateKey) // Parse private key

	if err != nil { // Check for errors
		return err // Return found error
	}

	*(*account).PrivateKey = *privateKey // Set private key

	return nil // No error occurred, return nil
}

// String - convert given account to string
func (account *Account) String() string {
	marshaled, _ := json.MarshalIndent(*account, "", "  ") // Marshal account

	return string(marshaled) // Return marshaled
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

/* END INTERNAL METHODS */
