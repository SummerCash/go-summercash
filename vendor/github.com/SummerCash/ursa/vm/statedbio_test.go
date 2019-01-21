package vm

import (
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"testing"
)

// TestStateDatabaseFromBytes - test functionality of bytes to state database decoder
func TestStateDatabaseFromBytes(t *testing.T) {
	env := &Environment{EnableJIT: false, DefaultMemoryPages: 128, DefaultTableSize: 65536} // Init env

	if env == nil { // Check nil env
		t.Fatal("nil environment configuration") // Panic
	}

	abs, err := filepath.Abs(filepath.FromSlash("../examples/main.wasm")) // Get absolute path to test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	testSourceFile, err := ioutil.ReadFile(abs) // Read test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	vm, err := NewVirtualMachine(testSourceFile, *env, new(NopResolver), nil) // Init vm

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	entryID, ok := vm.GetFunctionExport("main") // Get main func

	if ok != true { // Check for errors
		t.Fatal(entryID) // Panic
	}

	_, err = vm.Run(entryID) // Execute

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	stateEntry := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, 0) // Init state entry

	stateDb := NewStateDatabase(stateEntry) // Init state db

	if stateDb == nil { // Check for nil state db
		t.Fatal("nil state database") // Panic
	}

	stateDb, err = StateDatabaseFromBytes(stateDb.Bytes()) // Read from bytes

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(stateDb) // Log success
}

// TestWriteToMemory - test functionality of outbound stateDB I/O
func TestWriteToMemory(t *testing.T) {
	env := &Environment{EnableJIT: false, DefaultMemoryPages: 128, DefaultTableSize: 65536} // Init env

	if env == nil { // Check nil env
		t.Fatal("nil environment configuration") // Panic
	}

	abs, err := filepath.Abs(filepath.FromSlash("../examples/main.wasm")) // Get absolute path to test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	testSourceFile, err := ioutil.ReadFile(abs) // Read test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	vm, err := NewVirtualMachine(testSourceFile, *env, new(NopResolver), nil) // Init vm

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	entryID, ok := vm.GetFunctionExport("main") // Get main func

	if ok != true { // Check for errors
		t.Fatal(entryID) // Panic
	}

	_, err = vm.Run(entryID) // Execute

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	stateEntry := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, 0) // Init state entry

	stateDb := NewStateDatabase(stateEntry) // Init state db

	if stateDb == nil { // Check for nil state db
		t.Fatal("nil state database") // Panic
	}

	err = stateDb.WriteToMemory() // Write to persistent storage

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log("success") // Log success
}

// TestReadFromMemory - test functionality of inbound stateDB I/O
func TestReadFromMemory(t *testing.T) {
	env := &Environment{EnableJIT: false, DefaultMemoryPages: 128, DefaultTableSize: 65536} // Init env

	if env == nil { // Check nil env
		t.Fatal("nil environment configuration") // Panic
	}

	abs, err := filepath.Abs(filepath.FromSlash("../examples/main.wasm")) // Get absolute path to test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	testSourceFile, err := ioutil.ReadFile(abs) // Read test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	vm, err := NewVirtualMachine(testSourceFile, *env, new(NopResolver), nil) // Init vm

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	entryID, ok := vm.GetFunctionExport("main") // Get main func

	if ok != true { // Check for errors
		t.Fatal(entryID) // Panic
	}

	_, err = vm.Run(entryID) // Execute

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	stateEntry := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, 0) // Init state entry

	stateDb := NewStateDatabase(stateEntry) // Init state db

	if stateDb == nil { // Check for nil state db
		t.Fatal("nil state database") // Panic
	}

	err = stateDb.WriteToMemory() // Write to persistent storage

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	stateDb, err = ReadStateDBFromMemory(hex.EncodeToString(stateDb.ID)) // Read state DB

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(stateDb) // Log success
}
