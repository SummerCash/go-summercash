package config

import (
	"encoding/json"
	"io/ioutil"
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

	json.Unmarshal(rawJSON, &readJSON) // Unmarshal to buffer

	config := &ChainConfig{ // Init config
		Origin:    time.Now().UTC(),
		NetworkID: readJSON["networkID"].(uint),
		ChainID:   common.NewHash(crypto.Sha3(append(genesisSignature.Bytes(), []byte(strconv.Itoa(readJSON["networkID"].(int)))...))),
	}

	return config, nil // Return initialized chainConfig
}
