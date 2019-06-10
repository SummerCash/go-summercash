package config

// HandleReceivedConfigRequest - handle config request
func HandleReceivedConfigRequest() ([]byte, error) {
	config, err := ReadChainConfigFromMemory() // Read config from memory
	if err != nil {                            // Check for errors
		return nil, err // Return found error
	}

	return config.Bytes(), nil // Return byte value
}
