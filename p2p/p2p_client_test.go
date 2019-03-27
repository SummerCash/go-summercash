// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"context"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewClient tests the functionality of the NewClient helper method.
func TestNewClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	host, err := NewHost(ctx, 1234) // Initialize host

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	client := NewClient(host) // Initialize client

	if client == nil { // Check for nil client
		t.Fatal("nil client") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
