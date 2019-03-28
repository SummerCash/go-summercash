// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"context"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestStartServingStreams tests the functionality of the StartServingStreams helper method.
func TestStartServingStreams(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	host, err := NewHost(ctx, 1234) // Initialize host

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	client := NewClient(host) // Initialize client

	err = client.StartServingStreams("test_network") // Start serving streams

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
