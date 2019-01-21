package opcodes

import "testing"

func TestOpcode(t *testing.T) {
	opcode := Opcode(byte(0)) // Init opcode

	t.Log(opcode) // Log success
}
