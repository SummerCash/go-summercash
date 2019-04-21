package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/SummerCash/go-summercash/common"
)

// WriteToMemory - write given account to persistent memory
func (account *Account) WriteToMemory() error {
	err := account.MakeEncodingSafe() // Make safe for encoding

	if err != nil { // Check for errors
		return err // Return error
	}

	err = common.CreateDirIfDoesNotExist(filepath.FromSlash(fmt.Sprintf("%s/keystore", common.DataDir))) // Create dir if necessary

	if err != nil { // Check for errors
		return err // Return error
	}

	json, err := json.MarshalIndent(*account, "", "  ") // Marshal account

	if err != nil { // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/keystore/account_%s.json", common.DataDir, account.Address.String())), json, 0644) // Write json

	if err != nil { // Check for errors
		return err // Return found error
	}

	return account.RecoverSafeEncoding() // Recover
}

// ReadAccountFromMemory - read account with address from persistent memory
func ReadAccountFromMemory(address common.Address) (*Account, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/keystore/account_%s.json", common.DataDir, address.String()))) // Read account file

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	buffer := &Account{} // Init buffer

	err = json.Unmarshal(data, buffer) // Unmarshal into buffer

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	err = buffer.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return &Account{}, err // Return error
	}

	return buffer, nil // Return read account
}
