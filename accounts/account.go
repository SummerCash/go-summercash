package accounts

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
)

// Account - container holding account metadata, private key
type Account struct {
	Address              common.Address    `json:"address"`      // Account address
	PrivateKey           *ecdsa.PrivateKey `json:"privateKey"`   // Account private key
	SerializedPrivateKey []byte            `json:"s_privateKey"` // Serialized account private key
}

/* BEGIN EXPORTED METHODS */

// NewAccount - create new account
func NewAccount() (*Account, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	chain, err := types.NewChain(account.Address) // Init account chain

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	err = chain.WriteToMemory() // Write chain to memory

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	return account, nil // Return initialized account
}

// NewContractAccount - create new account for contract
func NewContractAccount(contractSource []byte) (*Account, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	chain, err := types.NewContractChain(account.Address, contractSource) // Init contract chain

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	err = chain.WriteToMemory() // Write chain to memory

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	return account, nil // Return initialized account
}

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

// GetAllAccounts - get list of local account addresses
func GetAllAccounts() ([]string, error) {
	buffer := []string{} // Init buffer

	files, err := ioutil.ReadDir(filepath.FromSlash(fmt.Sprintf("%s/keystore", common.DataDir))) // Walk keystore dir

	if err != nil { // Check for errors
		return []string{}, err // Return found error
	}

	for x, file := range files { // Iterate through files
		if x == 0 { // Check is first index
			buffer = []string{strings.Split(strings.Split(file.Name(), "account_")[1], ".json")[0]} // Init buffer
		} else {
			buffer = append(buffer, strings.Split(strings.Split(file.Name(), "account_")[1], ".json")[0]) // Append to buffer
		}
	}

	return buffer, nil // No error occurred, return success
}

// GetAllContracts - get list of all deployed contracts from account
func GetAllContracts(deployingAccount common.Address) ([]string, error) {
	buffer := []string{} // Init buffer

	files, err := ioutil.ReadDir(filepath.FromSlash(fmt.Sprintf("%s/db/chain", common.DataDir))) // Walk chain dir

	if err != nil { // Check for errors
		return []string{}, err // Return found error
	}

	for _, file := range files { // Iterate through files
		chainBytes, err := ioutil.ReadFile(file.Name()) // Read file

		if err != nil { // Check for errors
			return []string{}, err // Return found error
		}

		chain, err := types.FromBytes(chainBytes) // Get chain

		if err != nil { // Check for errors
			return []string{}, err // Return found error
		}

		if chain.ContractSource != nil { // Check is contract
			buffer = append(buffer, strings.Split(strings.Split(file.Name(), "chain_")[1], ".json")[0]) // Append to buffer
		}
	}

	return buffer, nil // No error occurred, return success
}

// MakeEncodingSafe - make account safe for encoding
func (account *Account) MakeEncodingSafe() error {
	if account.PrivateKey != nil { // Check has private key
		marshaledPrivateKey, err := x509.MarshalECPrivateKey(account.PrivateKey) // Marshal private key

		if err != nil { // Check for errors
			return err // Return found error
		}

		pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshaledPrivateKey}) // Encode to memory

		(*account).SerializedPrivateKey = pemEncoded // Set serialized

		*(*account).PrivateKey = ecdsa.PrivateKey{} // Set nil
	}

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

// Bytes - convert given account to byte array
func (account *Account) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*account) // Serialize account

	return buffer.Bytes() // Return serialized
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

/* END INTERNAL METHODS */
