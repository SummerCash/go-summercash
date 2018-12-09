package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/space55/summertech-blockchain/common"
)

// TestHandleReceivedCoordinationNode - test functionality of HandleReceivedCoordinationNode() method
func TestHandleReceivedCoordinationNode(t *testing.T) {
	coordinationChain, err := NewCoordinationChain(0) // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.WriteToMemory() // Write coordination chain to memory

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

	byteVal := coordinationNode.Bytes() // Get byte val

	err = HandleReceivedCoordinationNode(byteVal) // Handle node

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log("success") // Log success
}
