package config

import "testing"

// TestHandleReceivedConfigRequest - handle config request
func TestHandleReceivedConfigRequest(t *testing.T) {
	chainConfig, err := NewChainConfig("genesis.json") // Initialize chain configuration

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = chainConfig.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	configBytes, err := HandleReceivedConfigRequest() // Handle request

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chainConfig, err = FromBytes(configBytes) // Decode bytes

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*chainConfig) // Log success
}
