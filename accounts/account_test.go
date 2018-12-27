package accounts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

// TestNewAccount - test functionality of account generation
func TestNewAccount(t *testing.T) {
	account, err := NewAccount() // Generate account

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(account.String()) // Log success
}

// TestAccountFromKey - test functionality of account generation given a privateKey x
func TestAccountFromKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(account.String()) // Log success
}

// TestGetAllAccounts - test functionality of keystore walk method
func TestGetAllAccounts(t *testing.T) {
	account, err := NewAccount() // Generate account

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.WriteToMemory() // Make sure we have at least one account to walk

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	addresses, err := GetAllAccounts() // Walk

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(addresses) // Log success
}

// TestMakeEncodingSafe - test functionality of safe account encoding
func TestMakeEncodingSafe(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(account.String()) // Log success
}

// TestRecoverSafeEncoding - test functionality of safe account encoding recovery
func TestRecoverSafeEncoding(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.RecoverSafeEncoding() // Recover from safe encoding

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(account.String()) // Log success
}

// TestString - test functionality of string account serialization
func TestString(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := account.String() // Get string val

	if stringVal == "" { // Check nil
		t.Error("invalid string val") // Log error
		t.FailNow()                   // Panic
	}

	t.Log(stringVal) // Log success
}

// TestBytes - test byte array serialization of account
func TestBytes(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := account.Bytes() // Get byte val

	if byteVal == nil { // Check is nil
		t.Error("invalid byte val") // Log error
		t.FailNow()                 // Panic
	}

	t.Log(byteVal) // Log success
}
