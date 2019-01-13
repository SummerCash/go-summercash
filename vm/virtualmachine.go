package vm

// VirtualMachine - container holding wasm env metadata
type VirtualMachine struct {
	Environment *Environment `json:"environment"` // Environment settings
}

// Environment - container holding VM configuration variables
type Environment struct {
}
