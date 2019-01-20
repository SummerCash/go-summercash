package vm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/SummerCash/ursa/common"
)

// Environment - VM config, vars
type Environment struct {
	EnableJIT bool `json:"JIT"` // Should enable JIT

	MaxMemoryPages    int `json:"maxMemPages"`       // Max memory pages at given time
	MaxTableSize      int `json:"maxTableSize"`      // Max mem table size at given time
	MaxValueSlots     int `json:"maxValueSlots"`     // Max value slots at given time
	MaxCallStackDepth int `json:"maxCallStackDepth"` // Max call stack length/depth at given time

	DefaultMemoryPages int `json:"defaultMemPages"`  // Default num mem pages at given time
	DefaultTableSize   int `json:"defaultTableSize"` // Preset table size

	GasLimit uint64 `json:"gasLimit"` // Gas limit

	DisableFloatingPoint     bool `json:"disableFloat"`      // Remove float capacity
	ReturnOnGasLimitExceeded bool `json:"returnOnGasExceed"` // Panic on exceed specified gas limit
}

/* BEGIN EXPORTED METHODS */

// NewEnvironment - initialize virtual machine config with given params
func NewEnvironment(maxMemoryPages int, maxTableSize int, maxValueSlots int, maxCallStackDepth int, defaultMemoryPages int, defaultTableSize int, gasLimit uint64, disableFloatingPoint bool, returnOnGasLimitExceeded bool) *Environment {
	return &Environment{ // Return initialized vm config
		MaxMemoryPages:           maxMemoryPages,
		MaxTableSize:             maxTableSize,
		MaxValueSlots:            maxValueSlots,
		MaxCallStackDepth:        maxCallStackDepth,
		DefaultMemoryPages:       defaultMemoryPages,
		DefaultTableSize:         defaultTableSize,
		GasLimit:                 gasLimit,
		DisableFloatingPoint:     disableFloatingPoint,
		ReturnOnGasLimitExceeded: returnOnGasLimitExceeded,
	}
}

/*
	BEGIN FORMATTING METHODS
*/

// String - get string representation of given env
func (environment *Environment) String() string {
	marshaledVal, _ := json.MarshalIndent(*environment, "", "  ") // Marshal to JSON

	return string(marshaledVal) // Return string value
}

// Bytes - get byte array representation of given env
func (environment *Environment) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*environment, "", "  ") // Marshal to JSON

	return marshaledVal // Return success
}

// EnvironmentFromBytes - marshal byte array into environment struct
func EnvironmentFromBytes(b []byte) (*Environment, error) {
	buffer := &Environment{} // Initialize buffer

	err := json.Unmarshal(b, buffer) // Read json into buffer

	if err != nil { // Check for errors
		return &Environment{}, err // Return error
	}

	return buffer, nil // No error occurred, return read environment
}

/*
	END FORMATTING METHODS
*/

/*
	BEGIN I/O METHODs
*/

// WriteToMemory - write env to memory
func (environment *Environment) WriteToMemory() error {
	json, err := json.MarshalIndent(*environment, "", "  ") // Marshal env

	if err != nil { // Check for errors
		return err // Return error
	}

	err = common.CreateDirIfDoesNotExit(fmt.Sprintf("%s/config", common.DataDir)) // Create dir if necessary

	if err != nil { // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/config/environment.json", common.DataDir)), json, 0644) // Write environment to JSON

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}

// ReadEnvironmentFromMemory - read env from memory
func ReadEnvironmentFromMemory() (*Environment, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/config/environment.json", common.DataDir))) // Read file

	if err != nil { // Check for errors
		return &Environment{}, err // Return error
	}

	buffer := &Environment{} // Initialize buffer

	err = json.Unmarshal(data, buffer) // Read json into buffer

	if err != nil { // Check for errors
		return &Environment{}, err // Return error
	}

	return buffer, nil // No error occurred, return read config
}

/*
	END I/O METHODS
*/

/* END EXPORTED METHODS */
