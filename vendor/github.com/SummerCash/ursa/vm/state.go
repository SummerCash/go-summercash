package vm

import (
	"encoding/json"
	"errors"

	"github.com/SummerCash/ursa/common"
	"github.com/SummerCash/ursa/crypto"
)

var (
	// ErrNilStateEntry - describes an error regarding a stateDB of 0 entry length
	ErrNilStateEntry = errors.New("no state entries found")
)

// StateEntry - state entry header
type StateEntry struct {
	State *State `json:"state"` // State

	Nonce uint64 `json:"nonce"` // Entry nonce

	ID []byte `json:"ID"` // Entry ID
}

// State - state node
type State struct {
	CallStack    []Frame `json:"call_stack"`    // VM call stack
	CurrentFrame int     `json:"current_frame"` // Current callstack frame

	Table []uint32 `json:"table"` // VM mem/runtime table

	Globals []int64 `json:"globals"` // Global vrs

	Memory []byte `json:"memory"` // Virtual machine memory

	NumValueSlots int `json:"num_val_slots"` // Num of used value slots

	Yielded int64 `json:"yielded"` // Did yield

	InsideExecute bool `json:"inside_execute"` // Inside execute

	Exited    bool        `json:"has_exited"` // Did exit
	ExitError interface{} `json:"exit_err"`   // Error on exit

	ReturnValue int64 `json:"return"` // Return value

	Gas              uint64 `json:"gas"`                // Gas usage
	GasLimitExceeded bool   `json:"gas_limit_exceeded"` // Has exceeded given gas limit

	StateChildren []*StateEntry `json:"children"` // State children

	ID []byte `json:"ID"` // State ID
}

/* BEGIN EXPORTED METHODS */

// NewStateEntry - initialize new state entry
func NewStateEntry(callStack []Frame, currentFrame int, table []uint32, globals []int64, memory []byte, numValueSlots int, yielded int64, insideExecute bool, exited bool, exitError interface{}, returnValue int64, gas uint64, gasLimitExceeded bool, nonce uint64) *StateEntry {
	state := &State{
		CallStack:        callStack,        // Set call stack
		CurrentFrame:     currentFrame,     // Set current frame
		Table:            table,            // Set table
		Globals:          globals,          // Set globals
		Memory:           memory,           // Set memory
		NumValueSlots:    numValueSlots,    // Set value slots
		Yielded:          yielded,          // Set yielded
		InsideExecute:    insideExecute,    // Set inside execute
		Exited:           exited,           // Set has exited
		ExitError:        exitError,        // Set exit error
		ReturnValue:      returnValue,      // Set return value
		Gas:              gas,              // Set gas
		GasLimitExceeded: gasLimitExceeded, // Set gas limit exceeded
	} // Init state

	(*state).ID = crypto.Sha3(state.Bytes()) // Hash

	entry := &StateEntry{
		State: state, // Set state
		Nonce: nonce, // Set nonce
	} // Init state db entry

	(*entry).ID = crypto.Sha3(entry.Bytes()) // Hash

	return entry // Return success
}

// FindMax - find state child of greatest nonce value
func (stateEntry *StateEntry) FindMax() (*StateEntry, error) {
	if stateEntry.State == nil || stateEntry.State.StateChildren == nil || len(stateEntry.State.StateChildren) == 0 { // Check for errors
		return &StateEntry{}, ErrNilStateEntry // Return error
	}

	nilState := &StateEntry{Nonce: 0} // Init nil state

	lastState := &StateEntry{Nonce: 0} // Init state buffer

	for _, stateChild := range stateEntry.State.StateChildren { // Iterate through children
		if stateChild.Nonce > lastState.Nonce { // Check is max
			lastState = stateChild // Set last state
		}
	}

	if common.ByteIsEqual(lastState.Bytes(), nilState.Bytes()) { // Check no results
		return &StateEntry{}, ErrNilStateEntry // Return error
	}

	return lastState, nil // Return found max
}

/*
	BEGIN TYPE HELPERS
*/

// Bytes - get byte representation of state
func (state *State) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*state, "", "  ") // Marshal JSON

	return marshaledVal // Return success
}

// String - get string representation of state
func (state *State) String() string {
	marshaledVal, _ := json.MarshalIndent(*state, "", "  ") // Marshal JSON

	return string(marshaledVal) // Return success
}

// Bytes - get byte representation of entry
func (stateEntry *StateEntry) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*stateEntry, "", "  ") // Marshal JSON

	return marshaledVal // Return success
}

// String - get string representation of entry
func (stateEntry *StateEntry) String() string {
	marshaledVal, _ := json.MarshalIndent(*stateEntry, "", "  ") // Marshal JSON

	return string(marshaledVal) // Return success
}

/*
	END TYPE HELPERS
*/

/* END EXPORTED METHODS */
