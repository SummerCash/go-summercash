package types

import "testing"

/* BEGIN EXTERNAL METHDOS */

/*
	BEGIN COORDINATIONCHAIN METHODS
*/

// TestNewCoordinationChain - test coordinationChain initializer
func TestNewCoordinationChain(t *testing.T) {
	coordinationChain := NewCoordinationChain(0, &CoordinationNode{}) // Init coordinationChain

	if coordinationChain == nil { // Check for nil coordination chain
		t.Errorf("invalid coordination chain") // Log found error
		t.FailNow()                            // Panic
	}

	t.Log(*coordinationChain) // Log success
}

// TestBytesCoordinationChain - test functionality of coordinationChain Bytes() extension method
func TestBytesCoordinationChain(t *testing.T) {
	coordinationChain := NewCoordinationChain(0, &CoordinationNode{}) // Init coordinationChain

	if coordinationChain == nil { // Check for nil coordination chain
		t.Errorf("invalid coordination chain") // Log found error
		t.FailNow()                            // Panic
	}

	byteVal := coordinationChain.Bytes() // Get byte val

	if byteVal == nil { // Check for nil byte val
		t.Errorf("invalid byte val") // Log found error
		t.FailNow()                  // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringCoordinationChain - test functionality of coordinationChain String() extension method
func TestStringCoordinationChain(t *testing.T) {
	coordinationChain := NewCoordinationChain(0, &CoordinationNode{}) // Init coordinationChain

	if coordinationChain == nil { // Check for nil coordination chain
		t.Errorf("invalid coordination chain") // Log found error
		t.FailNow()                            // Panic
	}

	stringVal := coordinationChain.String() // Get string val

	if stringVal == "" { // Check for nil string val
		t.Errorf("invalid string val") // Log found error
		t.FailNow()                    // Panic
	}

	t.Log(stringVal) // Log success
}

/*
	END COORDINATIONCHAIN METHODS
*/

/* END EXTERNAL METHODS */
