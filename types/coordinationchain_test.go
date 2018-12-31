package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"strconv"
	"testing"

	"github.com/space55/summertech-blockchain/common"
)

/* BEGIN EXTERNAL METHDOS */

/*
	BEGIN COORDINATIONCHAIN METHODS
*/

// TestNewCoordinationChain - test coordinationChain initializer
func TestNewCoordinationChain(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationChain) // Log success
}

// TestAddNode - test addNode() method
func TestAddNode(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort)}) // Init coordination node

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

// TestJoinNetwork - test functionality of network joining
func TestJoinNetwork(t *testing.T) {
	err := JoinNetwork(common.BootstrapNodes[0], false) // Join network

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log("success") // Log success
}

// TestQueryAddress - test QueryAddress() method
func TestQueryAddress(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort)}) // Init coordination node

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

// TestQueryNode - test QueryNode() method
func TestQueryNode(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort)}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.AddNode(coordinationNode, false) // Add node

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // panic
	}

	coordinationNode, err = coordinationChain.QueryNode("1.1.1.1") // Query address

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationNode) // Log success
}

// TestQueryArchivalNode - test QueryArchivalNode() method
func TestQueryArchivalNode(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort)}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.AddNode(coordinationNode, false) // Add node

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // panic
	}

	coordinationNodes, err := coordinationChain.QueryArchivalNode("1.1.1.1") // Query address

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // Panic
	}

	t.Log(coordinationNodes) // Log success
}

// TestQueryAllArchivalNodes - test QueryAllArchivalNodes() method
func TestQueryAllArchivalNodes(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

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

	privateKey2, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address2, err := common.NewAddress(privateKey2) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort), "2.2.2.2:" + strconv.Itoa(common.NodePort)}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	coordinationNode2, err := NewCoordinationNode(address2, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort), "2.2.2.2:" + strconv.Itoa(common.NodePort)}) // Init coordination node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.AddNode(coordinationNode, false) // Add node

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // panic
	}

	err = coordinationChain.AddNode(coordinationNode2, false) // Add node

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // panic
	}

	coordinationNodes, err := coordinationChain.QueryAllArchivalNodes() // Query addresses

	if err != nil { // Check for errors
		t.Error(err) // Log error
		t.FailNow()  // Panic
	}

	t.Log(coordinationNodes) // Log success
}

// TestPushNode - test PushNode() method
func TestPushNode(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort)}) // Init coordination node

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

// TestCoordinationChainFromBytes - test bytes decoder for coordination c
func TestCoordinationChainFromBytes(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := coordinationChain.Bytes() // Get byte val

	if byteVal == nil { // Check for nil byte val
		t.Errorf("invalid byte val") // Log found error
		t.FailNow()                  // Panic
	}

	coordinationChain, err = CoordinationChainFromBytes(byteVal) // Decode byte value

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(coordinationChain.String()) // Log success
}

// TestBytesCoordinationChain - test functionality of coordinationChain Bytes() extension method
func TestBytesCoordinationChain(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

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
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort)}) // Init coordination node

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort)}) // Init coordination node

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort)}) // Init coordination node

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

	coordinationNode, err := NewCoordinationNode(address, []string{"1.1.1.1:" + strconv.Itoa(common.NodePort)}) // Init coordination node

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
