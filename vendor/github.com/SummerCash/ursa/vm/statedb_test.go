package vm

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

// TestNewStateDatabase - test functionality of state db init
func TestNewStateDatabase(t *testing.T) {
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

	t.Log(stateDb) // Log success
}

// TestAddState - test functionality of state appending
func TestAddState(t *testing.T) {
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

	t.Log(stateDb.States[len(stateDb.States)-1].State.ReturnValue) // Log state

	entryID, ok = vm.GetFunctionExport("main2") // Get main func

	if ok != true { // Check for errors
		t.Fatal(entryID) // Panic
	}

	_, err = vm.Run(entryID) // Execute

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	stateEntry2 := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, 1) // Init state entry

	err = stateDb.AddStateEntry(stateEntry2, stateEntry) // Add state entry

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(stateDb.States[len(stateDb.States)-1].State.ReturnValue) // Log state
}

// TestFindMaxDB - test functionality of DB-wide findMax() method
func TestFindMaxDB(t *testing.T) {
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

	t.Log(stateDb.States[len(stateDb.States)-1].State.ReturnValue) // Log state

	entryID, ok = vm.GetFunctionExport("main2") // Get main func

	if ok != true { // Check for errors
		t.Fatal(entryID) // Panic
	}

	_, err = vm.Run(entryID) // Execute

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	stateEntry2 := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, 1) // Init state entry

	err = stateDb.AddStateEntry(stateEntry2, stateEntry) // Add state entry

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(stateDb.States[len(stateDb.States)-1].State.ReturnValue) // Log state

	max, err := stateDb.FindMax() // Find max

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(max) // Log success
}

// TestQueryState - test functionality of state querying in db
func TestQueryState(t *testing.T) {
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

	stateEntry, err = stateDb.QueryState(stateEntry.ID) // Query for state

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(stateEntry) // Log success
}

// TestStringStateDatabase - test string serialization of state database
func TestStringStateDatabase(t *testing.T) {
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

	t.Log(stateDb) // Log success
}

// TestBytesStateDatabase - test functionality of byte array state database serialization
func TestBytesStateDatabase(t *testing.T) {
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

	byteVal := stateDb.Bytes() // Get byte value

	if byteVal == nil { // Check for invalid byte val
		t.Fatal("invalid serialized byte value") // Panic
	}

	t.Log(byteVal) // Log success
}
