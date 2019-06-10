package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
)

/* BEGIN EXPORTED METHODS */

// TestNewChain - test chain initializer
func TestNewChain(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := common.NewAddress(privateKey) // Generate address
	if err != nil {                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	coordinationChain, err := NewCoordinationChain() // Init coordinationChain
	if err != nil {                                  // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(address) // Initialize chain
	if err != nil {                 // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*chain) // Log initialized chain
}

// TestAddTransaction - test functionality of transaction appending
func TestAddTransaction(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	sender, err := common.NewAddress(privateKey) // Initialize address from private key
	if err != nil {                              // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = makeChainConfig(sender) // Make config

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created transaction: %s", transaction.Hash.String()) // Log issued tx

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(sender) // Initialize chain
	if err != nil {                // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = common.CreateDirIfDoesNotExist(fmt.Sprintf("%s/db/chain", common.DataDir)) // Create dir if necessary

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(chain.Account.String())

	err = common.WriteGob(fmt.Sprintf("%s/db/chain/chain_%s.gob", common.DataDir, chain.Account.String()), *chain) // Write gob

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created chain: %s", chain.ID.String()) // Log init

	err = chain.AddTransaction(transaction) // Add transaction

	if err != nil && !strings.Contains(err.Error(), "timed out") { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("added transaction: %s", transaction.Hash.String()) // Log signed

	coordinationChain, err := ReadCoordinationChainFromMemory() // Read chain from memory
	if err != nil {                                             // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	balance, err := coordinationChain.GetBalance(sender) // Get balance

	if err != nil && !strings.Contains(err.Error(), "timed out") { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("genesis: %s, chain: %s", balance.String(), chain.String()) // Log chain state
}

// TestQueryTransaction - test tx querying
func TestQueryTransaction(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	sender, err := common.NewAddress(privateKey) // Initialize address from private key
	if err != nil {                              // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = makeChainConfig(sender) // Make config

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created transaction: %s", transaction.Hash.String()) // Log issued tx

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(sender) // Initialize chain
	if err != nil {                // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = common.CreateDirIfDoesNotExist(fmt.Sprintf("%s/db/chain", common.DataDir)) // Create dir if necessary

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(chain.Account.String())

	err = common.WriteGob(fmt.Sprintf("%s/db/chain/chain_%s.gob", common.DataDir, chain.Account.String()), *chain) // Write gob

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created chain: %s", chain.ID.String()) // Log init

	err = chain.AddTransaction(transaction) // Add transaction

	if err != nil && !strings.Contains(err.Error(), "timed out") { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	transaction, err = chain.QueryTransaction(*transaction.Hash) // Query hash

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(transaction.String()) // Log success
}

// TestBytesChain - test chain to bytes conversion
func TestBytesChain(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := common.NewAddress(privateKey) // Generate address
	if err != nil {                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(address) // Initialize chain
	if err != nil {                 // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := chain.Bytes() // Get byte val

	if byteVal == nil { // Check nil byte val
		t.Errorf("invalid byte val") // Log error
		t.FailNow()                  // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringChain - test chain to string conversion
func TestStringChain(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := common.NewAddress(privateKey) // Generate address
	if err != nil {                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(address) // Initialize chain
	if err != nil {                 // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := chain.String() // Get string val

	if stringVal == "" { // Check nil string val
		t.Errorf("invalid string val") // Log error
		t.FailNow()                    // Panic
	}

	t.Log(stringVal) // Log success
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

type genesis struct {
	NetworkID uint `json:"networkID"`

	Alloc map[string]map[string]string `json:"alloc"`
}

// makeChainConfig - generate necessary config files
func makeChainConfig(address common.Address) error {
	alloc := make(map[string]map[string]string) // Init map

	alloc[address.String()] = make(map[string]string) // Init map

	alloc[address.String()]["balance"] = "500000000000000" // Set balance

	genesis := genesis{NetworkID: 0, Alloc: alloc} // Init genesis

	json, err := json.MarshalIndent(genesis, "", "  ") // Marshal genesis
	if err != nil {                                    // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile("genesis.json", json, 0644) // Write genesis to JSON

	if err != nil { // Check for errors
		return err // Return error
	}

	config, err := config.NewChainConfig("genesis.json") // Generate config
	if err != nil {                                      // Check for errors
		return err // Return error
	}

	err = config.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		return err // Return error
	}

	coordinationChain, err := NewCoordinationChain() // Init coordinationChain
	if err != nil {                                  // Check for errors
		return err // Return error
	}

	return coordinationChain.WriteToMemory() // Write to memory
}

/* END INTERNAL METHODS */
