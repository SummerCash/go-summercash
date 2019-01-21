package compiler

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

// TestLoadModule - test functionality of module loader
func TestLoadModule(t *testing.T) {
	abs, err := filepath.Abs(filepath.FromSlash("../examples/main.wasm")) // Get absolute path to test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	testSourceFile, err := ioutil.ReadFile(abs) // Read test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	module, err := LoadModule(testSourceFile) // Load module

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(module) // Log success
}

// TestString - test custom module stringer
func TestString(t *testing.T) {
	abs, err := filepath.Abs(filepath.FromSlash("../examples/main.wasm")) // Get absolute path to test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	testSourceFile, err := ioutil.ReadFile(abs) // Read test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	module, err := LoadModule(testSourceFile) // Load module

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(module.String()) // Log success
}

// TestCompileForInterpreter - test interpreter bytecode compile
func TestCompileForInterpreter(t *testing.T) {
	abs, err := filepath.Abs(filepath.FromSlash("../examples/main.wasm")) // Get absolute path to test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	testSourceFile, err := ioutil.ReadFile(abs) // Read test WASM file

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	module, err := LoadModule(testSourceFile) // Load module

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	simpleGasPolicy := SimpleGasPolicy{GasPerInstruction: 1} // Init simple gas policy

	interpreterCompiled, err := module.CompileForInterpreter(&simpleGasPolicy) // Compile for interpreter

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(interpreterCompiled) // Log success
}
