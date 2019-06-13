// Package db defines the standard go-summercash transaction database.
package db

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/crypto"
	"github.com/SummerCash/go-summercash/types"
)

/* BEGIN EXPORTED METHODS */

// ImportBlockmesh attempts to reconstruct the dag from a given blockmesh.
func ImportBlockmesh(dbPath string, inflationRate float64, networkID uint) (*Dag, error) {
	files, err := ioutil.ReadDir(filepath.FromSlash(dbPath)) // Walk keystore dir
	if err != nil {                                          // Check for errors
		return &Dag{}, err // Return found error
	}

	alloc := make(map[string]*big.Float) // Init address list buffer
	addresses := []common.Address{}      // Init address list buffer
	concatenatedAddresses := []byte{}    // Init concat addr buffer

	for _, file := range files { // Iterate through files
		parsed, err := common.StringToAddress(strings.Split(strings.Split(file.Name(), "chain_")[1], ".json")[0]) // Parse file name
		if err != nil {                                                                                           // Check for errors
			return &Dag{}, err // Return found error
		}

		chain, err := types.ReadChainFromMemory(parsed) // Read chain from persistent memory
		if err != nil {                                 // Check for errors
			return &Dag{}, err // Return found error
		}

		alloc[strings.Split(strings.Split(file.Name(), "chain_")[1], ".json")[0]] = chain.CalculateBalance() // Calculate balance
		addresses = append(addresses, parsed)                                                                // Append parsed to address list
		concatenatedAddresses = append(concatenatedAddresses, parsed[:]...)                                  // Append parsed
	}

	config := &config.ChainConfig{
		Alloc:          alloc,                                              // Set alloc
		AllocAddresses: addresses,                                          // Set alloc addresses
		InflationRate:  inflationRate,                                      // Set inflation rate
		NetworkID:      networkID,                                          // Set network ID
		ChainID:        common.NewHash(crypto.Sha3(concatenatedAddresses)), // Set chain ID
		ChainVersion:   config.Version,                                     // Set version
	} // Initialize chain config

	err = config.WriteToMemory() // Write config to persistent memory
	if err != nil {              // Check for errors
		return &Dag{}, err // Return found error
	}

	dag := NewDag() // Initialize dag

	err = dag.MakeGenesis(config) // Make genesis
	if err != nil {               // Check for errors
		return &Dag{}, err // Return found error
	}

	return dag, nil // Return initialized dag
}

// WriteToMemory writes the working dag to persistent memory.
func (dag *Dag) WriteToMemory(network string) error {
	err := common.CreateDirIfDoesNotExist(common.DagDir) // Create dag dir
	if err != nil {                                      // Check for errors
		return err // Return found error
	}

	file, err := os.Create(filepath.FromSlash(fmt.Sprintf("%s/dag_%s.gob", common.DagDir, network))) // Create dag file

	defer file.Close() // Close file

	if err != nil { // Check for errors
		return err // Return found error
	}

	encoder := gob.NewEncoder(file) // Initialize encoder

	err = encoder.Encode(dag) // Encode dag

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// ReadDagFromMemory attempts to reconstruct a dag from
// the local persisted dag db.
func ReadDagFromMemory(network string) (*Dag, error) {
	file, err := os.Open(filepath.FromSlash(fmt.Sprintf("%s/dag_%s.gob", common.DagDir, network))) // Open file

	defer file.Close() // Close file

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	decoder := gob.NewDecoder(file) // Initialize decoder

	buffer := &Dag{} // Init dag buffer

	err = decoder.Decode(buffer) // Decode file
	if err != nil {              // Check for errors
		return &Dag{}, err // Return found error
	}

	return buffer, nil // Return decoded dag
}

/* END EXPORTED METHODS */
