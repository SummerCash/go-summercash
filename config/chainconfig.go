package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/crypto"
)

// ChainConfig - chain configuration
type ChainConfig struct {
	Alloc map[string]float64 `json:"alloc"` // Account balances at genesis

	NetworkID uint        `json:"network"` // Network ID (0: mainnet, 1: testnet, etc...)
	ChainID   common.Hash `json:"id"`      // Hashed networkID, genesisSignature
}

const (
	// Version - dist version def
	Version = "0.1"
)

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

	alloc := make(map[string]float64) // Init alloc map

	for key, value := range readJSON["alloc"].(map[string]interface{}) { // Iterate through genesis addresses
		intVal, err := strconv.Atoi(value.(map[string]interface{})["balance"].(string)) // Get int val

		if err != nil { // Check for errors
			return &ChainConfig{}, err // Return error
		}

		alloc[key] = float64(intVal) // Set int val
	}

	config := &ChainConfig{ // Init config
		Alloc:     alloc,
		NetworkID: uint(readJSON["networkID"].(float64)),
		ChainID:   common.NewHash(crypto.Sha3(append(rawJSON, []byte(strconv.Itoa(int(readJSON["networkID"].(float64))))...))), // Generate chainID
	}

	return config, nil // Return initialized chainConfig
}

// Bytes - convert given chainConfig to byte array
func (chainConfig *ChainConfig) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*chainConfig) // Serialize config

	return buffer.Bytes() // Return serialized
}

// String - convert given chainConfig to string
func (chainConfig *ChainConfig) String() string {
	marshaled, _ := json.MarshalIndent(*chainConfig, "", "  ") // Marshal config

	return string(marshaled) // Return marshaled
}

// WriteToMemory - write given chainConfig to memory
func (chainConfig *ChainConfig) WriteToMemory() error {
	json, err := json.MarshalIndent(*chainConfig, "", "  ") // Marshal config

	if err != nil { // Check for errors
		return err // Return error
	}

	err = common.CreateDirIfDoesNotExit(fmt.Sprintf("%s/config", common.DataDir)) // Create dir if necessary

	if err != nil { // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/config/config.json", common.DataDir)), json, 0644) // Write chainConfig to JSON

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}

// ReadChainConfigFromMemory - read chain configuration from chain config json file
func ReadChainConfigFromMemory() (*ChainConfig, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/config/config.json", common.DataDir))) // Read file

	if err != nil { // Check for errors
		return &ChainConfig{}, err // Return error
	}

	buffer := &ChainConfig{} // Initialize buffer

	err = json.Unmarshal(data, buffer) // Read json into buffer

	if err != nil { // Check for errors
		return &ChainConfig{}, err // Return error
	}

	return buffer, nil // No error occurred, return read config
}
