package types

import "testing"

// TestWriteCoordinationChainToMemory - test functionality of WriteToMemory() method
func TestWriteCoordinationChainToMemory(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationChain) // Log success
}

// TestReadCoordinationChainFromMemory - test functionality of ReadCoordinationChainFromMemory() method
func TestReadCoordinationChainFromMemory(t *testing.T) {
	coordinationChain, err := NewCoordinationChain() // Init coordinationChain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = coordinationChain.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	coordinationChain, err = ReadCoordinationChainFromMemory() // Read coordination chain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationChain) // Log success
}
