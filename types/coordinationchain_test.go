package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/space55/summertech-blockchain/common"
)

/* BEGIN EXTERNAL METHDOS */

/*
	BEGIN COORDINATIONCHAIN METHODS
*/

// TestNewCoordinationChain - test coordinationChain initializer
func TestNewCoordinationChain(t *testing.T) {
	coordinationChain, err := NewCoordinationChain(0) // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationChain) // Log success
}

// TestAddNode - test addNode() method
func TestAddNode(t *testing.T) {
	coordinationChain, err := NewCoordinationChain(0) // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1"}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.AddNode(coordinationNode, false) // Add node

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // panic
	}

	t.Log(*coordinationNode) // Log success
}

// TestQueryAddress - test QueryAddress() method
func TestQueryAddress(t *testing.T) {
	coordinationChain, err := NewCoordinationChain(0) // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1"}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.AddNode(coordinationNode, false) // Add node

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // panic
	}

	coordinationNode, err = coordinationChain.QueryAddress(address) // Query address

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationNode) // Log success
}

// TestPushNode - test PushNode() method
func TestPushNode(t *testing.T) {
	coordinationChain, err := NewCoordinationChain(0) // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1"}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.AddNode(coordinationNode, false) // Add node

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // panic
	}

	err = coordinationChain.PushNode(coordinationNode) // Push changes

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationChain) // Log success
}

// TestBytesCoordinationChain - test functionality of coordinationChain Bytes() extension method
func TestBytesCoordinationChain(t *testing.T) {
	coordinationChain, err := NewCoordinationChain(0) // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := coordinationChain.Bytes() // Get byte val

	if byteVal == nil { // Check for nil byte val
		t.Errorf("invalid byte val") // Log found error
		t.FailNow()                  // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringCoordinationChain - test functionality of coordinationChain String() extension method
func TestStringCoordinationChain(t *testing.T) {
	coordinationChain, err := NewCoordinationChain(0) // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := coordinationChain.String() // Get string val

	if stringVal == "" { // Check for nil string val
		t.Errorf("invalid string val") // Log found error
		t.FailNow()                    // Panic
	}

	t.Log(stringVal) // Log success
}

/*
	END COORDINATIONCHAIN METHODS
*/

/*
	BEGIN COORDINATIONNODE METHODS
*/

// TestNewCoordinationNode - test functionality of coordinationNode initializer
func TestNewCoordinationNode(t *testing.T) {
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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1"}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationNode) // Log success
}

// TestCoordinationNodeFromBytes - test conversion from byte array to coordination node
func TestCoordinationNodeFromByte(t *testing.T) {
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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1"}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := coordinationNode.Bytes() // Get byteVal

	if byteVal == nil { // Check for nil byte val
		t.Errorf("invalid byte val") // Log found error
		t.FailNow()                  // Panic
	}

	coordinationNode, err = CoordinationNodeFromBytes(byteVal) // Get coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationNode) // Log success
}

// TestBytesCoordinationNode - test functionality of coordinationNode Byte() extension method
func TestBytesCoordinationNode(t *testing.T) {
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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1"}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := coordinationNode.Bytes() // Get byteVal

	if byteVal == nil { // Check for nil byte val
		t.Errorf("invalid byte val") // Log found error
		t.FailNow()                  // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringCoordinationNode - test functionality of coordinationNode String() extension method
func TestStringCoordinationNode(t *testing.T) {
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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1"}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := coordinationNode.String() // Get stringVal

	if stringVal == "" { // Check for nil string val
		t.Errorf("invalid string val") // Log found error
		t.FailNow()                    // Panic
	}

	t.Log(stringVal) // Log success
}

/*
	END COORDINATIONNODE METHODS
*/

/* END EXTERNAL METHODS */
