package vm

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/SummerCash/ursa/common"
)

// StateDatabaseFromBytes - marshal given byte array into state database
func StateDatabaseFromBytes(b []byte) (*StateDatabase, error) {
	var stateBytes bytes.Buffer // Init state bytes buffer

	_, err := stateBytes.Write(b) // Write to buffer

	if err != nil { // Check for errors
		return &StateDatabase{}, err // Return found error
	}

	stateBuffer := &StateDatabase{} // Init state buffer

	decoder := gob.NewDecoder(&stateBytes) // Init gob decoder

	err = decoder.Decode(stateBuffer) // Decode state bytes

	if err != nil { // Check for errors
		return &StateDatabase{}, err // Return found error
	}

	return stateBuffer, nil // Return read state
}

// WriteToMemory - write state database to persistent storage with given data path
func (stateDB *StateDatabase) WriteToMemory() error {
	var stateBuffer bytes.Buffer // Init encoding buffer

	encoder := gob.NewEncoder(&stateBuffer) // Write to state buffer

	err := encoder.Encode(*stateDB) // Encode virtual machine state

	if err != nil { // Check for errors
		return err // Return found error
	}

	abs, err := filepath.Abs(filepath.FromSlash(fmt.Sprintf("%s/state", common.DataDir))) // Get absolute dir

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = common.CreateDirIfDoesNotExit(abs) // Create data dir if not exist

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/%s.db", abs, hex.EncodeToString(stateDB.ID))), stateBuffer.Bytes(), 0644) // Write state buffer to persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// ReadStateDBFromMemory - read state database from persistent memory
func ReadStateDBFromMemory(id string) (*StateDatabase, error) {
	abs, err := filepath.Abs(filepath.FromSlash(fmt.Sprintf("%s/state", common.DataDir))) // Get absolute dir

	if err != nil { // Check for errors
		return &StateDatabase{}, err // Return found error
	}

	stateBuffer := &StateDatabase{} // Init state buffer

	stateFile, err := os.Open(filepath.FromSlash(fmt.Sprintf("%s/%s.db", abs, id))) // Read state bytes

	if err != nil { // Check for errors
		return &StateDatabase{}, err // Return found error
	}

	decoder := gob.NewDecoder(stateFile) // Init gob decoder

	err = decoder.Decode(stateBuffer) // Decode state bytes

	if err != nil { // Check for errors
		return &StateDatabase{}, err // Return found error
	}

	stateFile.Close() // Close state file

	return stateBuffer, nil // Return read state
}
