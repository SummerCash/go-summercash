package vm

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

// TestNewState - test functionality of state entry init
func TestNewStateEntry(t *testing.T) {
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

	t.Log(stateEntry) // Log state entry
}

// TestFindMax - test functionality of max search
func TestFindMax(t *testing.T) {
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

	stateEntry2 := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, 1) // Init state entry

	stateEntry3 := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, 2) // Init state entry

	stateEntry4 := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, 3) // Init state entry

	(*stateEntry).State.StateChildren = []*StateEntry{stateEntry2, stateEntry3, stateEntry2, stateEntry4} // Set children

	max, err := stateEntry.FindMax() // Find max child

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(max) // Log success
}
