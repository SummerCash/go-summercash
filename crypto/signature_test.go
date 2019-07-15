// Package crypto implements a set of abstractions enabling the creation,
// verification, and validation of ECDSA signatures. Additionally, the
// aforementioned package implements several helper types and methods for
// hashing via Blake2s.
package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestSign tests the functionality of the Sign() helper method.
func TestSign(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(PreferredCurve(), rand.Reader) // Generate private key
	if err != nil {                                                     // Check for errors
		t.Fatal(err) // Panic
	}

	signature, err := Sign(privateKey, []byte("test")) // Sign test message
	if err != nil {                                    // Check for errors
		t.Fatal(err) // Panic
	}

	if signature == nil { // Check nil signature
		t.Fatal("nil signature") // Panic
	}
}

// TestSignatureVerify tests the functionality of the Signature Verify() helper method.
func TestSignatureVerify(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(PreferredCurve(), rand.Reader) // Generate private key
	if err != nil {                                                     // Check for errors
		t.Fatal(err) // Panic
	}

	signature, err := Sign(privateKey, []byte("test")) // Sign test message
	if err != nil {                                    // Check for errors
		t.Fatal(err) // Panic
	}

	if !signature.Verify() { // Check is not valid
		t.Fatal("signature should be valid") // Panic
	}
}

// TestSignatureString tests the functionality of the Signature String() helper method.
func TestSignatureString(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(PreferredCurve(), rand.Reader) // Generate private key
	if err != nil {                                                     // Check for errors
		t.Fatal(err) // Panic
	}

	signature, err := Sign(privateKey, []byte("test")) // Sign test message
	if err != nil {                                    // Check for errors
		t.Fatal(err) // Panic
	}

	stringVal, err := signature.String() // Get signature string value
	if err != nil {                      // Check for errors
		t.Fatal(err) // Panic
	}

	if stringVal == "" { // Check empty string value
		t.Fatal("string value is empty") // Panic
	}
}

// TestSignatureBytes tests the functionality of the Signature Bytes() helper method.
func TestSignatureBytes(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(PreferredCurve(), rand.Reader) // Generate private key
	if err != nil {                                                     // Check for errors
		t.Fatal(err) // Panic
	}

	signature, err := Sign(privateKey, []byte("test")) // Sign test message
	if err != nil {                                    // Check for errors
		t.Fatal(err) // Panic
	}

	bytesVal, err := signature.Bytes() // Get byte slice value
	if err != nil {                    // Check for errors
		t.Fatal(err) // Panic
	}

	if bytesVal == nil { // Check nil bytes value
		t.Fatal("bytes val should not be nil") // Panic
	}
}

// TestSignatureFromBytes tests the functionality of the SignatureFromBytes() helper method.
func TestSignatureFromBytes(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(PreferredCurve(), rand.Reader) // Generate private key
	if err != nil {                                                     // Check for errors
		t.Fatal(err) // Panic
	}

	signature, err := Sign(privateKey, []byte("test")) // Sign test message
	if err != nil {                                    // Check for errors
		t.Fatal(err) // Panic
	}

	bytesVal, err := signature.Bytes() // Get byte slice value
	if err != nil {                    // Check for errors
		t.Fatal(err) // Panic
	}

	decodedSignature, err := SignatureFromBytes(bytesVal) // Decode bytes val
	if err != nil {                                       // Check for errors
		t.Fatal(err) // Panic
	}

	if !decodedSignature.Verify() { // Check decoded signature also valid
		t.Fatal("decoded signature should also be valid") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
