// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"context"
	"testing"

	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/validator"
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

	config := &config.ChainConfig{} // Init empty config

	standardValidator := validator.NewStandardValidator(config) // Initialize validator

	validator := validator.Validator(standardValidator) // Get interface value

	client := NewClient(host, &validator, "test_network") // Initialize client

	err = client.StartServingStreams() // Start serving streams

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
