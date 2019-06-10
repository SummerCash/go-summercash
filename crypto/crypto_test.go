package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"
)

// TestBytesSignature - test functionality of signature to bytes extension method
func TestBytesSignature(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, []byte("test")) // Sign
	if err != nil {                                                  // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	signature := Signature{ // Init signature
		PublicKey: &privateKey.PublicKey,
		Time:      time.Now().UTC(),
		V:         []byte("test"),
		R:         r,
		S:         s,
	}

	byteVal := signature.Bytes() // Get byte val

	if byteVal == nil { // Check for nil byte val
		t.Errorf("invalid byte value") // Log found error
		t.FailNow()                    // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringSignature - test functionality of signature to string extension method
func TestStringSignature(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, []byte("test")) // Sign
	if err != nil {                                                  // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	signature := Signature{ // Init signature
		PublicKey: &privateKey.PublicKey,
		Time:      time.Now().UTC(),
		V:         []byte("test"),
		R:         r,
		S:         s,
	}

	stringVal := signature.String() // Get string val

	if stringVal == "" { // Check for nil string val
		t.Errorf("invalid string value") // Log found error
		t.FailNow()                      // Panic
	}

	t.Log(stringVal) // Log success
}

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

// TestSha3n - test functionality of sha3n hashing function
func TestSha3n(t *testing.T) {
	hashed := Sha3n([]byte("test"), 10) // Hash

	if hashed == nil { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}

// TestSha3nString - test functionality of sha3n string hashing function
func TestSha3nString(t *testing.T) {
	hashed := Sha3nString([]byte("test"), 10) // Hash

	if hashed == "" { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}

// TestSha3d - test functionality of sha3d hashing function
func TestSha3d(t *testing.T) {
	hashed := Sha3d([]byte("test")) // Hash

	if hashed == nil { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}

// TestSha3dString - test functionality of sha3d hashing string function
func TestSha3dString(t *testing.T) {
	hashed := Sha3dString([]byte("test")) // Hash

	if hashed == "" { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}
