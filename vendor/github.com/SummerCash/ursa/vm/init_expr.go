package vm

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/SummerCash/wagon/wasm/leb128"
	ops "github.com/SummerCash/wagon/wasm/operators"
)

// readU32 reads an unsigned 32-bit integer from a reader.
func readU32(r io.Reader) (uint32, error) {
	var buf [4]byte // Init buffer

	_, err := io.ReadFull(r, buf[:]) // Read into buffer

	if err != nil { // Check for errors
		return 0, err // Return error
	}

	return binary.LittleEndian.Uint32(buf[:]), nil // Return success
}

// readU64 reads an unsigned 64-bit integer from a reader.
func readU64(r io.Reader) (uint64, error) {
	var buf [8]byte // Init buffer

	_, err := io.ReadFull(r, buf[:]) // Read into buffer

	if err != nil { // Check for errors
		return 0, err // Return error
	}

	return binary.LittleEndian.Uint64(buf[:]), nil // Return success
}

// executeInitExpr - execute and return the result of a WebAssembly init expression
func execInitExpr(expr []byte, globals []int64) int64 {
	var stack []int64 // Init stack buffer

	r := bytes.NewReader(expr) // Read expressions

	for { // Iterate
		b, err := r.ReadByte() // Read bytes

		if err == io.EOF { // Check for EOF
			break // Break
		} else if err != nil { // Check for errors
			panic(err) // Panic
		}

		switch b { // Handle different operation types
		case ops.I32Const: // Check is I32Const opcode
			i, err := leb128.ReadVarint32(r) // Read

			if err != nil { // Check for errors
				panic(err) // Panic
			}
			stack = append(stack, int64(i)) // Append to stack
		case ops.I64Const:
			i, err := leb128.ReadVarint64(r) // Read

			if err != nil { // Check for errors
				panic(err) // Panic
			}
			stack = append(stack, int64(i)) // Append to stack
		case ops.F32Const:
			i, err := readU32(r) // Read

			if err != nil { // Check for errors
				panic(err) // Panic
			}
			stack = append(stack, int64(i)) // Append to stack
		case ops.F64Const:
			i, err := readU64(r) // Read

			if err != nil { // Check for errors
				panic(err) // Panic
			}
			stack = append(stack, int64(i)) // Append to stack
		case ops.GetGlobal:
			index, err := leb128.ReadVarUint32(r) // Read

			if err != nil { // Check for errors
				panic(err) // Panic
			}
			stack = append(stack, globals[int(index)]) // Append to stack
		case ops.End: // Check is end opcode
			break // Break
		default:
			panic("invalid opcode in init expr") // Panic
		}
	}

	return stack[len(stack)-1] // Return stack
}
