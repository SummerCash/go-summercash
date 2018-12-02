package common

import "testing"

/* BEGIN EXPORTED METHODS */

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

// TestEncode - test functionality of Encode() method
func TestEncode(t *testing.T) {
	encoded, err := Encode([]byte("test")) // Encode

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(encoded) // Log success
}

// TestEncodeString - test functionality of EncodeString() method
func TestEncodeString(t *testing.T) {
	encoded, err := EncodeString([]byte("test")) // Encode

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(encoded) // Log success
}

// TestDecode - test functionality of Decode() method
func TestDecode(t *testing.T) {
	encoded, err := Encode([]byte("test")) // Encode

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	decoded, err := Decode(encoded) // Decode

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(string(decoded)) // Log success
}

// TestDecodeString - test functionality of DecodeString() method
func TestDecodeString(t *testing.T) {
	encoded, err := EncodeString([]byte("test")) // Encode

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	decoded, err := DecodeString(encoded) // Decode

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(string(decoded)) // Log success
}

/* END EXPORTED METHODS */
