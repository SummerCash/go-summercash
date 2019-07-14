package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
)

/* BEGIN EXPORTED METHODS */

// GetAllLocalizedChains gets a list of the locally-provided chains, and their addresses.
func GetAllLocalizedChains() ([]string, error) {
	err := common.CreateDirIfDoesNotExist(filepath.FromSlash(fmt.Sprintf("%s/db/chain", common.DataDir))) // Make chain dir
	if err != nil {                                                                                       // Check for errors
		return []string{}, err // Return found error
	}

	buffer := []string{} // Init buffer

	files, err := ioutil.ReadDir(filepath.FromSlash(fmt.Sprintf("%s/db/chain", common.DataDir))) // Walk keystore dir
	if err != nil {                                                                              // Check for errors
		return []string{}, err // Return found error
	}

	for _, file := range files { // Iterate through files
		buffer = append(buffer, strings.Split(strings.Split(file.Name(), "chain_")[1], ".json")[0]) // Append to buffer
	}

	return buffer, nil // No error occurred, return success
}

// WriteToMemory - write given chain to memory
func (chain *Chain) WriteToMemory() error {
	err := common.CreateDirIfDoesNotExist(fmt.Sprintf("%s/db/chain", common.DataDir)) // Create dir if necessary
	if err != nil {                                                                   // Check for errors
		return err // Return error
	}

	err = chain.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		return err // Return error
	}

	json, err := json.Marshal(*chain) // Marshal chain
	if err != nil {                   // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/db/chain/chain_%s.json", common.DataDir, chain.Account.String())), json, 0644) // Write JSON

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = chain.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}

// ReadGenesisChainFromMemory reads a genesis chain based on a given chain config.
func ReadGenesisChainFromMemory(config *config.ChainConfig) (*Chain, error) {
	genesis := config.AllocAddresses[0] // Get genesis

	genesisChain, err := ReadChainFromMemory(genesis) // Read genesis
	if err != nil {                                   // Check for errors
		return &Chain{}, err // Return found error
	}

	return genesisChain, nil // Return read chain
}

// ReadChainFromMemory - read chain from memory
func ReadChainFromMemory(address common.Address) (*Chain, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/db/chain/chain_%s.json", common.DataDir, address.String()))) // Read chain
	if err != nil {                                                                                                              // Check for errors
		return &Chain{}, err // Return error
	}

	buffer := &Chain{} // Initialize buffer

	err = json.Unmarshal(data, buffer) // Read json into buffer

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	err = buffer.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	return buffer, nil // No error occurred, return read coordinationChain
}

/* END EXPORTED METHODS */
