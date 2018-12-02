package crypto

import "testing"

// TestSha3 - test functionality of sha3 hashing function
func TestSha3(t *testing.T) {
	hashed := Sha3([]byte("test")) // Hash

	if hashed == nil { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}

// TestSha3String - test functionality of sha3 hashing string function
func TestSha3String(t *testing.T) {
	hashed := Sha3String([]byte("test")) // Hash

	if hashed == "" { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}
