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

// TestStringChainConfig - test conversion from chainConfig to string
func TestStringChainConfig(t *testing.T) {
	chainConfig, err := NewChainConfig("genesis.json") // Initialize chain configuration

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := chainConfig.String() // Get string val

	if stringVal == "" { // Check for nil string value
		t.Errorf("invalid string val") // Log found error
		t.FailNow()                    // Panic
	}

	t.Log(stringVal) // Log string value
}

// TestBytesChainConfig - test conversion from chainConfig to bytes
func TestBytesChainConfig(t *testing.T) {
	chainConfig, err := NewChainConfig("genesis.json") // Initialize chain configuration

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := chainConfig.Bytes() // Get byte val

	if byteVal == nil { // Check for nil byte value
		t.Errorf("invalid byte val") // Log found error
		t.FailNow()                  // Panic
	}

	t.Log(byteVal) // Log string value
}

// TestFromBytes - test decode byte array into chain config
func TestFromBytes(t *testing.T) {
	chainConfig, err := NewChainConfig("genesis.json") // Initialize chain configuration

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := chainConfig.Bytes() // Get byte value

	chainConfig, err = FromBytes(byteVal) // Decode bytes

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*chainConfig) // Log success
}

// TestWriteToMemoryChainConfig - test i/o for chainConfig
func TestWriteToMemoryChainConfig(t *testing.T) {
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

	t.Logf("wrote chain config to memory") // Log success
}

// TestReadChainConfigFromMemory - test read chain config from json file
func TestReadChainConfigFromMemory(t *testing.T) {
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

	chainConfig, err = ReadChainConfigFromMemory() // Read chain config

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*chainConfig) // Log success
}
