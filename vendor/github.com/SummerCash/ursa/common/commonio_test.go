package common

import "testing"

// TestGetDefaultDataDir - test functionality of data dir getter
func TestGetDefaultDataDir(t *testing.T) {
	dataDir := GetDefaultDataDir() // Get default data dir

	if dataDir == "" { // Check for nil data dir
		t.Fatal("nil data dir") // Panic
	}

	t.Log(dataDir) // Log success
}

// TestCreateDirIfDoesNotExit - test functionality of dir creation /not exist
func TestCreateDirIfDoesNotExist(t *testing.T) {
	err := CreateDirIfDoesNotExit(DataDir) // Create data dir if not exist

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = CreateDirIfDoesNotExit(".../test") // Test should fail

	if err != nil { // Check for errors
		t.Logf("expected error: %s", err) // Log success
	}
}
