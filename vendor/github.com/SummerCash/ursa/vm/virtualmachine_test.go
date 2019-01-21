package vm

import (
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/SummerCash/ursa/compiler"
)

// TestNewVirtualMachine - test functionality of vm init
func TestNewVirtualMachine(t *testing.T) {
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

	t.Log(vm) // Log success
}

// TestSaveState - test functionality of outbound state I/O
func TestSaveState(t *testing.T) {
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

	err = vm.SaveState() // Save state

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(vm.ReturnValue) // Log state

	entryID, ok = vm.GetFunctionExport("main2") // Get main func

	if ok != true { // Check for errors
		t.Fatal(entryID) // Panic
	}

	_, err = vm.Run(entryID) // Execute

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = vm.SaveState() // Save state

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestResetToState - test functionality of state reverting
func TestResetToState(t *testing.T) {
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

	err = vm.SaveState() // Save state

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(vm.ReturnValue) // Log state

	entryID, ok = vm.GetFunctionExport("main2") // Get main func

	if ok != true { // Check for errors
		t.Fatal(entryID) // Panic
	}

	_, err = vm.Run(entryID) // Execute

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = vm.SaveState() // Save state

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(vm.ReturnValue) // Log state

	err = vm.ResetToState(vm.StateDB.States[1].ID) // Reset to state

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(vm.ReturnValue) // Log success
}

// TestLoadWorkingRoot - test functionality of inbound state root I/O
func TestLoadWorkingRoot(t *testing.T) {
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

	err = vm.SaveState() // Save state

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(vm.ReturnValue) // Log state

	entryID, ok = vm.GetFunctionExport("main2") // Get main func

	if ok != true { // Check for errors
		t.Fatal(entryID) // Panic
	}

	_, err = vm.Run(entryID) // Execute

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(vm.ReturnValue) // Log state

	err = vm.SaveState() // Save state

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = vm.LoadWorkingRoot() // Load working root

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(vm.ReturnValue) // Log state
}

// TestLoadStateDB - test functionality of inbound state DB I/O
func TestLoadStateDB(t *testing.T) {
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

	err = vm.SaveState() // Save state

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = vm.LoadWorkingRoot() // Load working root

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = vm.LoadStateDB(hex.EncodeToString(vm.StateDB.ID)) // Load state DB

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(vm.StateDB) // Log success
}

// TestRun - test functionality of vm run
func TestRun(t *testing.T) {
	abs, err := filepath.Abs(filepath.FromSlash("../examples/main.wasm")) // Get absolute path to test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	testSourceFile, err := ioutil.ReadFile(abs) // Read test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	simpleGasPolicy := &compiler.SimpleGasPolicy{GasPerInstruction: 1} // Init simple gas policy

	vm, err := NewVirtualMachine(testSourceFile, Environment{}, &NopResolver{}, simpleGasPolicy) // Init vm

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	entryID, ok := vm.GetFunctionExport("main") // Get main func

	if ok != true { // Check for errors
		t.Fatal(entryID) // Panic
	}

	result, err := vm.Run(entryID) // Execute

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(result) // Log success
}
