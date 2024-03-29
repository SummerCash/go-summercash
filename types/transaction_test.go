package types

import (
	"bytes"
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
)

// TestNewTransaction - test functionality of tx initializer
func TestNewTransaction(t *testing.T) {
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

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	marshaledVal, err := json.MarshalIndent(*transaction, "", "  ") // Marshal tx
	if err != nil {                                                 // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(string(marshaledVal)) // Log success
}

// TestNewContractCreation - test functionality of contract initializer
func TestNewContractCreation(t *testing.T) {
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

	contractSource, err := ioutil.ReadFile("main.wasm") // Read test smart contract
	if err != nil {                                     // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	transaction, err := NewContractCreation(0, nil, &sender, nil, big.NewFloat(0), contractSource) // Initialize transaction
	if err != nil {                                                                                // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(string(transaction.Bytes())) // Log success
}

// TestPublishTransaction - test functionality of transaction.Publish() method
func TestPublishTransaction(t *testing.T) {
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

	t.Logf("signed transaction: %s", transaction.Signature.String()) // Log signed

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

	err = common.WriteGob(fmt.Sprintf("%s/db/chain/chain_%s.gob", common.DataDir, chain.Account.String()), *chain) // Write gob

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created chain: %s", chain.ID.String()) // Log init

	err = transaction.Publish() // Publish transaction

	if err != nil && !strings.Contains(err.Error(), "timed out") { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("published transaction: " + transaction.String()) // Log success
}

// TestMakeEncodingSafe - test functionality of transaction.MakeEncodingSafe() method
func TestMakeEncodingSafe(t *testing.T) {
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

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created transaction: %s", transaction.Hash.String()) // Log issued tx

	err = transaction.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(transaction.String()) // Log success
}

// TestRecoverSafeEncoding - test functionality of tx recovery from safe encoding
func TestRecoverSafeEncoding(t *testing.T) {
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

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created transaction: %s", transaction.Hash.String()) // Log issued tx

	err = transaction.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = transaction.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(transaction.String()) // Log success
}

// TestTransactionFromBytes - test transaction serialization from byte array
func TestTransactionFromBytes(t *testing.T) {
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

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	transaction, err = TransactionFromBytes(transaction.Bytes()) // Get transaction from bytes

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(transaction.String()) // Log transaction
}

// TestBytes - test transaction to bytes conversion
func TestBytes(t *testing.T) {
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

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = SignTransaction(transaction, privateKey) // Sign
	if err != nil {                                // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := transaction.Bytes() // Get byte val

	if byteVal == nil || bytes.Contains(byteVal, []byte{'\r'}) { // Check for nil byteVal
		t.Errorf("invalid byteval") // Log found error
		t.FailNow()                 // Panic
	}

	t.Log(byteVal) // Log success
}

// TestString - test transaction to string conversion
func TestString(t *testing.T) {
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

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := transaction.String() // Convert to string

	if stringVal == "" { // Check for nil val
		t.Errorf("invalid stringval") // Return found error
		t.FailNow()                   // Panic
	}

	t.Log(stringVal) // Log success
}

// TestWriteTransactionToMemory - test functionality of transaction outbound I/O
func TestWriteTransactionToMemory(t *testing.T) {
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

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = transaction.WriteToMemory() // Write transaction to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log("success") // Log success
}

// TestReadTransactionFromMemory - test functionality of transaction inbound I/O
func TestReadTransactionFromMemory(t *testing.T) {
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

	transaction, err := NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                               // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = transaction.WriteToMemory() // Write transaction to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	transaction, err = ReadTransactionFromMemory(*transaction.Hash) // Read TX

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(transaction.String()) // Log success
}
