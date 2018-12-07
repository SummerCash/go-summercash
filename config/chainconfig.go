package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/crypto"
)

// ChainConfig - chain configuration
type ChainConfig struct {
	Origin time.Time `json:"origin"` // Time at chain initialization

	GenesisSignature *crypto.Signature `json:"genesis"` // Signature of genesis address

	NetworkID uint        `json:"network"` // Network ID (0: mainnet, 1: testnet, etc...)
	ChainID   common.Hash `json:"id"`      // Hashed networkID, genesisSignature
}

// NewChainConfig - generate new ChainConfig from genesis.json file
func NewChainConfig(genesisFilePath string) (*ChainConfig, error) {
	rawJSON, err := ioutil.ReadFile(genesisFilePath) // Read genesis file

	if err != nil { // Check for errors
		return &ChainConfig{}, err // Return error
	}

	var readJSON map[string]interface{} // Init buffer

	err = json.Unmarshal(rawJSON, &readJSON) // Unmarshal to buffer

	if err != nil { // Check for errors
		return &ChainConfig{}, err // Return error
	}

	config := &ChainConfig{ // Init config
		Origin:    time.Now().UTC(),
		NetworkID: uint(readJSON["networkID"].(float64)),
		ChainID:   common.NewHash(crypto.Sha3(append(rawJSON, []byte(strconv.Itoa(int(readJSON["networkID"].(float64))))...))), // Generate chainID
	}

	return config, nil // Return initialized chainConfig
}

// WriteToMemory - write given chainConfig to memory
func (chainConfig *ChainConfig) WriteToMemory() error {
	json, err := json.MarshalIndent(*chainConfig, "", "  ") // Marshal config

	if err != nil { // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/config/config.json", common.DataDir)), json, 0644) // Write chainConfig to JSON

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}
