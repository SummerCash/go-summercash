package accounts

import "testing"

// TestWriteToMemory - test functionality of persistent memory account I/O writing
func TestWriteToMemory(t *testing.T) {
	account, err := NewAccount() // Generate account
	if err != nil {              // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log("success") // Log success
}

// TestReadFromMemory - test functionality of account reader
func TestReadFromMemory(t *testing.T) {
	account, err := NewAccount() // Generate account
	if err != nil {              // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err = ReadAccountFromMemory(account.Address) // Read account from persistent memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(account.String()) // Log success
}
