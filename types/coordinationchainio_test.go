package types

import "testing"

// TestWriteCoordinationChainToMemory - test functionality of WriteToMemory() method
func TestWriteCoordinationChainToMemory(t *testing.T) {
	coordinationChain := NewCoordinationChain(0, &CoordinationNode{}) // Init coordinationChain

	if coordinationChain == nil { // Check for nil coordination chain
		t.Errorf("invalid coordination chain") // Log found error
		t.FailNow()                            // Panic
	}

	err := coordinationChain.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*coordinationChain) // Log success
}

// TestReadCoordinationChainFromMemory - test functionality of ReadCoordinationChainFromMemory() method
func TestReadCoordinationChainFromMemory(t *testing.T) {
	coordinationChain := NewCoordinationChain(0, &CoordinationNode{}) // Init coordinationChain

	if coordinationChain == nil { // Check for nil coordination chain
		t.Errorf("invalid coordination chain") // Log found error
		t.FailNow()                            // Panic
	}

	err := coordinationChain.WriteToMemory() // Write to memory

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
