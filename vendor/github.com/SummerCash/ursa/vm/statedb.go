package vm

import (
	"encoding/json"
	"errors"

	"github.com/SummerCash/ursa/common"
	"github.com/SummerCash/ursa/crypto"
)

var (
	// ErrStateAlreadyExists - describes an error regarding a state addition in a state database already containing the given state
	ErrStateAlreadyExists = errors.New("state already exists in given state DB")

	// ErrInvalidStateNonce - describes an error regarding a state addition with a nonce less than the last provided in the given state db
	ErrInvalidStateNonce = errors.New("invalid state nonce")
)

// StateDatabase - database holding vm states
type StateDatabase struct {
	States    []*StateEntry `json:"states"` // All states (not in tree)
	StateRoot *StateEntry   `json:"root"`   // VM state root

	WorkingRoot *StateEntry `json:"working_root"` // Working VM state root

	MerkleRoot []byte `json:"merkle_root"` // State merkle root

	ID []byte `json:"ID"` // State DB ID
}

/* BEGIN EXPORTED METHODS */

// NewStateDatabase - initialize state DB with given root
func NewStateDatabase(rootState *StateEntry) *StateDatabase {
	stateDB := &StateDatabase{
		States:     []*StateEntry{rootState},       // Set states
		StateRoot:  rootState,                      // Set root
		MerkleRoot: crypto.Sha3(rootState.Bytes()), // Set merkle root
	}

	(*stateDB).ID = crypto.Sha3(stateDB.Bytes()) // Set db id
	(*stateDB).WorkingRoot = rootState           // Set working root

	return stateDB // Return init db
}

// AddStateEntry - add state entry to given state database at given root state
func (stateDB *StateDatabase) AddStateEntry(state *StateEntry, rootState *StateEntry) error {
	if rootState == nil { // Check for nil state
		rootState = stateDB.WorkingRoot // Set root state to working root
	}

	_, err := stateDB.QueryState(state.ID) // Check state already in DB

	if err == nil { // Check for already existent state
		return ErrStateAlreadyExists // Return error
	}

	if state.Nonce <= rootState.Nonce { // Check invalid nonce
		return ErrInvalidStateNonce // Return error
	}

	(*(*rootState).State).StateChildren = append((*(*rootState).State).StateChildren, state) // Append state
	(*stateDB).States = append((*stateDB).States, state)                                     // Append to general states

	(*stateDB).WorkingRoot = state // Set working root

	err = stateDB.WriteToMemory() // Write db to memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// SetWorkingRoot - set current working root (similar to a "git checkout COMMIT_HASH")
func (stateDB *StateDatabase) SetWorkingRoot(rootState *StateEntry) {
	stateDB.WorkingRoot = rootState // Set working root
}

// QueryState - query state in db by identifier
func (stateDB *StateDatabase) QueryState(id []byte) (*StateEntry, error) {
	if len(stateDB.States) == 0 || stateDB.States == nil { // Check no states
		return &StateEntry{}, ErrNilStateEntry // Return error
	}

	for _, state := range stateDB.States { // Iterate through states
		if common.ByteIsEqual(state.ID, id) { // Check for match
			return state, nil // Return found state
		}
	}

	return &StateEntry{}, ErrNilStateEntry // Return error
}

// FindMax - find state entry with max nonce value
func (stateDB *StateDatabase) FindMax() (*StateEntry, error) {
	if stateDB.States == nil || len(stateDB.States) == 0 || stateDB.StateRoot == nil { // Check for nil state db
		return &StateEntry{}, ErrNilStateEntry // Return error
	}

	lastEntry := stateDB.StateRoot // Set last
	currentEntry := &StateEntry{}  // Init buffer

	var err error // Init error buffer

	for lastEntry.Nonce >= currentEntry.Nonce { // Iterate until found max
		if currentEntry.State != nil { // Check is not nil state
			lastEntry = currentEntry // Set last entry
		}

		currentEntry, err = lastEntry.FindMax() // Find max

		if err != nil { // Check for errors
			return &StateEntry{}, err // Return found error
		}
	}

	return currentEntry, nil // Return found state entry
}

/*
	BEGIN TYPE HELPERS
*/

// Bytes - get byte representation of db
func (stateDB *StateDatabase) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*stateDB, "", "  ") // Marshal JSON

	return marshaledVal // Return success
}

// String - get string representation of db
func (stateDB *StateDatabase) String() string {
	marshaledVal, _ := json.MarshalIndent(*stateDB, "", "  ") // Marshal JSON

	return string(marshaledVal) // Return success
}

/*
	END TYPE HELPERS
*/

/* END EXPORTED METHODS */
