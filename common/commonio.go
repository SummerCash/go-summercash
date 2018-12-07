package common

import "os"

/* BEGIN EXPORTED METHODS */

// CreateDirIfDoesNotExit - create given directory if does not exist
func CreateDirIfDoesNotExit(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) { // Check dir exists
		err = os.MkdirAll(dir, 0755) // Create directory

		if err != nil { // Check for errors
			return err // Return error
		}
	}

	return nil // No error occurred
}

/* END EXPORTED METHODS */
