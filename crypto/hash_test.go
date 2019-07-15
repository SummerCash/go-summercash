// Package crypto implements a set of abstractions enabling the creation,
// verification, and validation of ECDSA signatures. Additionally, the
// aforementioned package implements several helper types and methods for
// hashing via Blake2s.
package crypto

import (
	"bytes"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestBlake2s tests the functionality of the Blake2s() helper method.
func TestBlake2s(t *testing.T) {
	hash := Blake2s([]byte("test")) // Hash test contents

	if hash.Bytes() == nil || !bytes.Equal(hash.Bytes(), []byte{243, 8, 252, 2, 206, 145, 114, 173, 2, 167, 215, 88, 0, 236, 252, 2, 113, 9, 188, 103, 152, 126, 163, 42, 186, 155, 141, 204, 123, 16, 21, 14}) { // Check not equal to determined correct byte slice value
		t.Fatal("invalid hash") // Panic
	}
}

// TestHashBytes tests the functionality of the hash Bytes() helper method
// (literally the exact same as the Blake2s test).
func TestHashBytes(t *testing.T) {
	hash := Blake2s([]byte("test")) // Hash test contents

	if hash.Bytes() == nil || !bytes.Equal(hash.Bytes(), []byte{243, 8, 252, 2, 206, 145, 114, 173, 2, 167, 215, 88, 0, 236, 252, 2, 113, 9, 188, 103, 152, 126, 163, 42, 186, 155, 141, 204, 123, 16, 21, 14}) { // Check not equal to determined correct byte slice value
		t.Fatal("invalid hash") // Panic
	}
}

// TestHashString tests the functionality of the hash String() helper method.
func TestHashString(t *testing.T) {
	hash := Blake2s([]byte("test")) // Hash test contents

	if hash.String() != "f308fc02ce9172ad02a7d75800ecfc027109bc67987ea32aba9b8dcc7b10150e" { // Check not equal to determined correct value
		t.Fatal("invalid hash string value") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
