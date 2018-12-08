package common

import (
	"encoding/gob"
	"os"
)

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

// WriteGob - create gob from specified object, at filePath
func WriteGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath) // Attempt to create file at path

	if err != nil { // Check for errors
		return err // Return found error
	}

	encoder := gob.NewEncoder(file) // Write to file

	err = encoder.Encode(object) // Encode object

	if err != nil { // Check for errors
		return err // Return found error
	}

	file.Close() // Close file operation

	return err // Return error (might be nil)
}

// ReadGob - read gob specified at path
func ReadGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath) // Attempt to open file at path

	if err != nil { // Check for errors
		return err // Return found error
	}

	decoder := gob.NewDecoder(file) // Attempt to decode gob

	err = decoder.Decode(object) // Assign to error

	if err != nil { // Check for errors
		return err // Return found error
	}

	file.Close() // Close file

	return err // Return error
}

/* END EXPORTED METHODS */
