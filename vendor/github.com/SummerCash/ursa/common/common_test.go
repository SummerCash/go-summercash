package common

import "testing"

// TestByteIsEqual - test functionality of byte cmp method
func TestByteIsEqual(t *testing.T) {
	a := []byte("test")  // Init test byte array
	b := []byte("test2") // Init test byte array

	t.Log(ByteIsEqual(a, b)) // Log is equal
}
