package common

import (
	"os"
	"path/filepath"
)

/* BEGIN EXPORTED METHODS */

/*
	BEGIN I/O METHODS
*/

// GetDefaultDataDir - get default data directory
func GetDefaultDataDir() string {
	abs, _ := filepath.Abs(filepath.FromSlash("./data")) // Get absolute dir

	return filepath.FromSlash(abs) // Match slashes
}

// CreateDirIfDoesNotExit - create given directory if does not exist
func CreateDirIfDoesNotExit(dir string) error {
	dir = filepath.FromSlash(dir) // Just to be safe

	if _, err := os.Stat(dir); os.IsNotExist(err) { // Check dir exists
		err = os.MkdirAll(dir, 0755) // Create directory

		if err != nil { // Check for errors
			return err // Return error
		}
	}

	return nil // No error occurred
}

/*
	END I/O METHODS
*/

/* END EXPORTED METHODS */
