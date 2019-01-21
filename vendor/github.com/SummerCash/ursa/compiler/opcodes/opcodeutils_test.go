package opcodes

import "testing"

// TestString - test opcode stringer
func TestString(t *testing.T) {
	opcode := Opcode(byte(0)) // Init test opcode

	if opcode.String() != "Nop" { // Check invalid opcode
		t.Fatal("invalid opcode string") // Panic
	}

	opcode2 := Opcode(byte(1)) // Init 2nd test opcode

	if opcode2.String() != "Unreachable" { // Check invalid opcode
		t.Fatal("invalid opcode string") // Panic
	}

	opcode3 := Opcode(byte(161)) // Init 3rd test opcode

	if opcode3.String() != "Unknown" { // Check invalid opcode
		t.Fatal("invalid opcode string") // Panic
	}

	t.Logf("opcode strings: %s, %s, %s", opcode.String(), opcode2.String(), opcode3.String()) // Log success
}
