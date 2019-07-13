package common

import (
	"bufio"
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ErrTimedOut defines an error representing an IO timeout.
var ErrTimedOut = errors.New("the io operation timed out")

/* BEGIN EXPORTED METHODS */

// ReadAll reads the entire contents from a given reader.
func ReadAll(reader *bufio.Reader) ([]byte, error) {
	data := make(chan []byte) // Initialize data buffer
	err := make(chan error)   // Initialize error buffer
	started := false          // Initialize started watch var

	var wg sync.WaitGroup // Initialize wait group

	wg.Add(1) // Only one process to complete

	go func(started *bool) {
		scanner := bufio.NewScanner(reader) // Initialize reader

		for scanner.Scan() { // Scan
			data <- append(<-data, scanner.Bytes()...) // Append read line

			if !*started { // Check hasn't started yet
				*started = true // Set started
			}
		}

		if scanErr := scanner.Err(); scanErr != nil { // Check for errors
			err <- scanErr // Write error to parent routine
		}

		wg.Done() // Done!
	}(&started) // Run with timeout

	startTime := time.Now() // Get start time

	for !started { // Wait until started
		if pickedUpErr := <-err; pickedUpErr != nil { // Check for errors
			return nil, pickedUpErr // Return found errors
		}

		if time.Now().Sub(startTime) > 2*time.Second { // Check timeout
			return nil, ErrTimedOut // Return timeout error
		}
	}

	wg.Wait() // Wait...

	return <-data, nil // Return read data
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
