package config

import "testing"

// TestNewChainConfig - test init method for chainConfigs
func TestNewChainConfig(t *testing.T) {
	chainConfig, err := NewChainConfig("genesis.json") // Initialize chain configuration

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*chainConfig) // Log config
}
