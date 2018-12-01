package common

import "testing"

// TestMemPrefix - test MemPrefix definition
func TestMemPrefix(t *testing.T) {
	if string(MemPrefix) != "0x" || len(MemPrefix) != 2 { // Check invalid memPrefix
		t.Errorf("invalid memory prefix %s", string(MemPrefix)) // Log found error
		t.FailNow()                                             // Panic
	}

	t.Log(string(MemPrefix)) // Log success
}

// TestErrNilInput - test ErrNilInput error definition
func TestErrNilInput(t *testing.T) {
	_, err := EncodeString([]byte{}) // Encode string

	if err != nil && err != ErrNilInput { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("success: %s", err.Error()) // Log success
}

// TestErrNoMem - test ErrNoMem error definition
func TestErrNoMem(t *testing.T) {
	_, err := EncodeString([]byte("0")) // Encode string

	if err != nil && err != ErrNoMem { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("success: %s", err.Error()) // Log success
}
