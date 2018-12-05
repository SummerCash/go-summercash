package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/space55/summertech-blockchain/crypto"
)

/* BEGIN EXPORTED METHODS */

/*
	BEGIN ADDRESS METHODS
*/

// TestNewAddress - test init new address
func TestNewAddress(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(address) // Log address
}

// TestPublicKeyToAddress - test functionality of public key address initializer
func TestPublicKeyToAddress(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address := PublicKeyToAddress(&privateKey.PublicKey) // Generate address

	t.Log(address) // Log address
}

func TestStringToAddress(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := address.String() // Get byte val

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	decoded, err := StringToAddress(stringVal) // Decode string

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(decoded) // Log success
}

// TestBytes - test address to bytes conversion
func TestBytes(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(address.Bytes()) // Log success
}

// TestString - test address to string conversion
func TestString(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(address) // Log success
}

/*
	END ADDRESS METHODS
*/

/*
	BEGIN HASH METHODS
*/

// TestNewHash - test functionality of hash initializer
func TestNewHash(t *testing.T) {
	hash := NewHash(crypto.Sha3([]byte("test"))) // Init hash

	if hash[:] == nil { // Check hash not nil
		t.Errorf("invalid hash") // Log found error
		t.FailNow()              // Panic
	}

	t.Log(hash) // Log success
}

// TestStringToHash - test functionality of StringToHash() method
func TestStringToHash(t *testing.T) {
	hash, err := StringToHash("0x307836f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5") // Decode to hash

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(hash) // Log decoded
}

// TestBytesHash - test functionality of hash bytes() extension method
func TestBytesHash(t *testing.T) {
	hash := NewHash(crypto.Sha3([]byte("test"))) // Init hash

	if hash[:] == nil { // Check hash not nil
		t.Errorf("invalid hash") // Log found error
		t.FailNow()              // Panic
	}

	byteVal := hash.Bytes() // Get byte val

	if byteVal == nil { // Check byte val not nil
		t.Errorf("invalid hash") // Log found error
		t.FailNow()              // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringHash - test functionality of hash string() extension method
func TestStringHash(t *testing.T) {
	hash := NewHash(crypto.Sha3([]byte("test"))) // Init hash

	if hash[:] == nil { // Check hash not nil
		t.Errorf("invalid hash") // Log found error
		t.FailNow()              // Panic
	}

	stringVal := hash.String() // Get string val

	if stringVal == "" { // Check byte val not nil
		t.Errorf("invalid hash") // Log found error
		t.FailNow()              // Panic
	}

	t.Log(stringVal) // Log success
}

/*
	END HASH METHODS
*/

/*
	BEGIN ADDRESS-SPACE METHODS
*/

// TestNewAddressSpace - test functionality of address-space initializer
func TestNewAddressSpace(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	addressSpace, err := NewAddressSpace([]Address{address}) // Init address-space

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(addressSpace) // Log initialized address-space
}

// TestBytesAddressSpace - test functionality of address-space Bytes() extension method
func TestBytesAddressSpace(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	addressSpace, err := NewAddressSpace([]Address{address}) // Init address-space

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := addressSpace.Bytes() // Get string val

	if byteVal == nil { // Check for nil byte value
		t.Errorf("invalid byte value") // Log invalid string val
		t.FailNow()                    // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringAddressSpace - test functionality of address-space String() extension method
func TestStringAddressSpace(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	addressSpace, err := NewAddressSpace([]Address{address}) // Init address-space

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := addressSpace.String() // Get string val

	if stringVal == "" { // Check for nil string value
		t.Errorf("invalid string value") // Log invalid string val
		t.FailNow()                      // Panic
	}

	t.Log(stringVal) // Log success
}

/*
	END ADDRESS-SPACE METHODS
*/

/* END EXPORTED METHODS */
