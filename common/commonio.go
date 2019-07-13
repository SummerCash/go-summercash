package common

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrTimedOut defines an error representing an IO timeout.
var ErrTimedOut = errors.New("the io operation timed out")

/* BEGIN EXPORTED METHODS */

// ReadAll reads the entire contents from a given reader.
func ReadAll(reader *bufio.Reader) ([]byte, error) {
	data := make(chan []byte) // Initialize data buffer
	err := make(chan error)   // Initialize error buffer

	go func() {
		scanner := bufio.NewScanner(reader) // Initialize reader

		var buffer []byte // Initialize buffer

		for scanner.Scan() { // Scan
			if bytes.Contains(scanner.Bytes(), []byte("\n\r")) || len(scanner.Bytes()) == 0 { // Check for EOF'
				break // Break
			}

			buffer = append(buffer, scanner.Bytes()...) // Append read line
		}

		if scanErr := scanner.Err(); scanErr != nil && scanErr.Error() != "stream reset" { // Check for errors
			err <- scanErr // Write error to parent routine

			return // Return
		}

		data <- buffer // Write contents of buffer to data chan var
	}() // Run with timeout

	deadline := time.Tick(2 * time.Second) // Wait 2 seconds to declare dead

	select {
	case pickedUpErr := <-err: // Wait for errors
		return nil, pickedUpErr // Return found errors
	case <-deadline:
		return nil, ErrTimedOut // Return timeout error
	case readData := <-data:
		return readData, nil // Return read data
	}
}

// CreateDirIfDoesNotExist - create given directory if does not exist
func CreateDirIfDoesNotExist(dir string) error {
	dir = filepath.FromSlash(dir) // Just to be safe

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
	if err != nil {                  // Check for errors
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
	if err != nil {                // Check for errors
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
