package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/config"
)

/* BEGIN EXPORTED METHODS */

// TestNewChain - test chain initializer
func TestNewChain(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := common.NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(address) // Initialize chain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*chain) // Log initialized chain
}

// TestAddTransaction - test functionality of transaction appending
func TestAddTransaction(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	sender, err := common.NewAddress(privateKey) // Initialize address from private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = makeChainConfig(sender) // Make config

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	transaction, err := NewTransaction(0, nil, &sender, &sender, 0, []byte("test")) // Initialize transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created transaction: %s", transaction.Hash.String()) // Log issued tx

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("signed transaction: %s", transaction.Signature.String()) // Log signed

	chain, err := NewChain(sender) // Initialize chain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	abs, _ := filepath.Abs(filepath.FromSlash(fmt.Sprintf("../%s", common.DataDir)))

	err = common.CreateDirIfDoesNotExit(fmt.Sprintf("%s/db/chain", abs)) // Create dir if necessary

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = common.WriteGob(fmt.Sprintf("%s/db/chain/chain_%s.gob", abs, chain.Account.String()), *chain) // Write gob

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

	t.Log("success: " + chain.String()) // Log success
}

// TestBytesChain - test chain to bytes conversion
func TestBytesChain(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := common.NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(address) // Initialize chain

	if err != nil { // Check for errors
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

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := common.NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(address) // Initialize chain

	if err != nil { // Check for errors
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

	alloc[address.String()]["balance"] = "5000000000000" // Set balance

	genesis := genesis{NetworkID: 0, Alloc: alloc} // Init genesis

	json, err := json.MarshalIndent(genesis, "", "  ") // Marshal genesis

	if err != nil { // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile("genesis.json", json, 0644) // Write genesis to JSON

	if err != nil { // Check for errors
		return err // Return error
	}

	config, err := config.NewChainConfig("genesis.json") // Generate config

	if err != nil { // Check for errors
		return err // Return error
	}

	err = config.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		return err // Return error
	}

	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

	if err != nil { // Check for errors
		return err // Return error
	}

	return coordinationChain.WriteToMemory() // Write to memory
}

/* END INTERNAL METHODS */
