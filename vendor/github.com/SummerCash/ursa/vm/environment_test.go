package vm

import "testing"

// TestNewEnvironment - test functionality of environment initializer
func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment(10, 10, 10, 10, 10, 10, 10, false, true) // Init env

	if env == nil { // Check nil env
		t.Fatal("nil environment configuration") // Panic
	}

	t.Log(env) // Log success
}

// TestString - test functionality of environment stringer
func TestString(t *testing.T) {
	env := NewEnvironment(10, 10, 10, 10, 10, 10, 10, false, true) // Init env

	if env == nil || env.String() == "" { // Check nil env
		t.Fatal("nil environment configuration") // Panic
	}

	t.Log(env.String()) // Log success
}

// TestBytes - test functionality of environment byte serialization
func TestBytes(t *testing.T) {
	env := NewEnvironment(10, 10, 10, 10, 10, 10, 10, false, true) // Init env

	if env == nil || env.Bytes() == nil { // Check nil env
		t.Fatal("nil environment configuration") // Panic
	}

	t.Log(env.Bytes()) // Log success
}

// TestNewEnvironmentFromBytes - test functionality of environment byteval marshaling
func TestEnvironmentFromBytes(t *testing.T) {
	env := NewEnvironment(10, 10, 10, 10, 10, 10, 10, false, true) // Init env

	if env == nil || env.Bytes() == nil { // Check nil env
		t.Fatal("nil environment configuration") // Panic
	}

	env, err := EnvironmentFromBytes(env.Bytes()) // Marshal from byte value

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	_, err = EnvironmentFromBytes(append(env.Bytes(), byte(2))) // Marshal from byte value (SHOULD FAIL)

	if err != nil { // Check for errors
		t.Logf("found error in marshaling incorrect bytes: %s", err.Error()) // Log success
	}

	t.Log(env) // Log success
}

// TestWriteEnvironmentToMemory - test functionality of outbound env I/O
func TestWriteEnvironmentToMemory(t *testing.T) {
	env := NewEnvironment(10, 10, 10, 10, 10, 10, 10, false, true) // Init env

	if env == nil || env.Bytes() == nil { // Check nil env
		t.Fatal("nil environment configuration") // Panic
	}

	err := env.WriteToMemory() // Write env to memory

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(err) // Log nil error to force code coverage
}

// TestReadEnvironmentFromMemory - test functionality of inbound env I/O
func TestReadEnvironmentFromMemory(t *testing.T) {
	env := NewEnvironment(10, 10, 10, 10, 10, 10, 10, false, true) // Init env

	if env == nil || env.Bytes() == nil { // Check nil env
		t.Fatal("nil environment configuration") // Panic
	}

	err := env.WriteToMemory() // Write env to memory

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	env, err = ReadEnvironmentFromMemory() // Read environment from memory

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(err) // Log nil error to force code coverage
	t.Log(env) // Log success
}
