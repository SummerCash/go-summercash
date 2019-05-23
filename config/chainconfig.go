package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"strconv"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/crypto"
)

// ChainConfig - chain configuration
type ChainConfig struct {
	Alloc map[string]*big.Float `json:"alloc"` // Account balances at genesis

	AllocAddresses []common.Address // Account addresses

	InflationRate float64 `json:"inflation"` // Inflation rate

	NetworkID    uint        `json:"network"` // Network ID (0: mainnet, 1: testnet, etc...)
	ChainID      common.Hash `json:"id"`      // Hashed networkID, genesisSignature
	ChainVersion string      `json:"version"` // Network version
}

const (
	// Version - dist version def
	Version = "0.6.95"
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

	alloc := make(map[string]*big.Float) // Init alloc map

	allocAddresses := []common.Address{}

	x := 0 // Init iterator

	for key, value := range readJSON["alloc"].(map[string]interface{}) { // Iterate through genesis addresses
		floatVal, _, _ := big.ParseFloat(value.(map[string]interface{})["balance"].(string), 10, 350, big.ToNearestEven) // Parse float

		address, err := common.StringToAddress(key) // Get address value

		if err != nil { // Check for errors
			return &ChainConfig{}, err // Return error
		}

		if x == 0 { // Check genesis
			allocAddresses = []common.Address{address} // Init slice
		} else { // Else
			allocAddresses = append(allocAddresses, address) // Append address
		}

		alloc[key] = floatVal // Set int val

		x++ // Increment iterator
	}

	config := &ChainConfig{ // Init config
		Alloc:          alloc,
		AllocAddresses: allocAddresses,
		NetworkID:      uint(readJSON["networkID"].(float64)),
		InflationRate:  readJSON["inflation"].(float64),
		ChainID:        common.NewHash(crypto.Sha3(append(rawJSON, []byte(strconv.Itoa(int(readJSON["networkID"].(float64))))...))), // Generate chainID
		ChainVersion:   Version,                                                                                                     // Set version
	}

	return config, nil // Return initialized chainConfig
}

// UpdateChainVersion updates the version of the given chain config.
func (chainConfig *ChainConfig) UpdateChainVersion() error {
	(*chainConfig).ChainVersion = Version // Update version

	return chainConfig.WriteToMemory() // Write chain config to persistent memory
}

// Bytes - convert given chainConfig to byte array
func (chainConfig *ChainConfig) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*chainConfig) // Serialize config

	return buffer.Bytes() // Return serialized
}

// FromBytes - decode byte array into chain config
func FromBytes(b []byte) (*ChainConfig, error) {
	buffer := &ChainConfig{} // Initialize buffer

	err := json.Unmarshal(b, buffer) // Read json into buffer

	if err != nil { // Check for errors
		return &ChainConfig{}, err // Return error
	}

	return buffer, nil // No error occurred, return read config
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

	err = common.CreateDirIfDoesNotExist(fmt.Sprintf("%s/config", common.DataDir)) // Create dir if necessary

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
