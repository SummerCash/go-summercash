// Package db defines the standard go-summercash transaction database.
package db

import (
	"encoding/gob"
	"errors"
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

// ErrNilDag is an error definition representing a root dag leaf of nil value.
var ErrNilDag = errors.New("dag has no root")

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

// Flatten flattens the working dag. If the working
// dag is nil, an error is returned.
func (dag *Dag) Flatten() (*Flattened, error) {
	if dag.Root == nil { // Check no root
		return &Flattened{}, ErrNilDag // Return error
	}

	leaves := []*Leaf{dag.Root} // Init leaves list

	children, err := dag.Root.GetChildren() // Get children
	if err != nil {                         // Check for errors
		return &Flattened{}, err // Return found error
	}

	leaves = append(leaves, children...) // Append children to leaves

	transactions := []*types.Transaction{} // Initialize transactions

	for _, leaf := range leaves { // Iterate through leaves
		transactionCopy := *leaf.Transaction // Get raw value

		err = transactionCopy.MakeEncodingSafe() // Make tx encoding safe

		if err != nil { // Check for errors
			return &Flattened{}, err // Return found error
		}

		transactions = append(transactions, &transactionCopy) // Append transaction
	}

	return &Flattened{
		Transactions: transactions,
	}, nil // Return flattened dag
}

// UnflattenDag attempts to unflatten the given flattened dag.
func UnflattenDag(flattened *Flattened) (*Dag, error) {
	dag := NewDag() // Initialize dag ref

	for _, transaction := range flattened.Transactions { // Iterate through transactions
		err := transaction.RecoverSafeEncoding() // Recover tx
		if err != nil {                          // Check for errors
			return &Dag{}, err // Return found error
		}

		err = dag.AddTransaction(transaction) // Add tx to dag
		if err != nil {                       // Check for errors
			return &Dag{}, err // Return found error
		}
	}

	return dag, nil // Return unflattened dag
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

	flattened, err := dag.Flatten() // Flatten dag
	if err != nil {                 // Check for errors
		return err // Return found error
	}

	err = encoder.Encode(flattened) // Encode dag
	if err != nil {                 // Check for errors
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

	buffer := &Flattened{} // Init flattened dag buffer

	err = decoder.Decode(buffer) // Decode file
	if err != nil {              // Check for errors
		return &Dag{}, err // Return found error
	}

	return UnflattenDag(buffer) // Return decoded, unflattened dag
}

/* END EXPORTED METHODS */
