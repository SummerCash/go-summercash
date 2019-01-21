package compiler

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/SummerCash/ursa/common"
	"github.com/SummerCash/ursa/compiler/opcodes"
	"github.com/SummerCash/ursa/crypto"
	"github.com/SummerCash/wagon/disasm"
	"github.com/SummerCash/wagon/wasm"
	"github.com/SummerCash/wagon/wasm/leb128"
)

// Module - wasm module
type Module struct {
	Base                 *wasm.Module   `json:"-"` // Base parsed module
	FunctionNames        map[int]string // Module functions
	DisableFloatingPoint bool           // Config to disable float ops
	Identifier           []byte         `json:"ID"` // Unique module identifier
}

// InterpreterCode - interpreter metadata
type InterpreterCode struct {
	NumRegs    int // Reg count
	NumParams  int // Param count
	NumLocals  int // Local vars count
	NumReturns int // Returns count

	Bytes   []byte      // Byte val
	JITInfo interface{} // Just-in-time meta
	JITDone bool        // Finished just-in-time compilation
}

// LoadModule - load WASM raw bytes module
func LoadModule(moduleBytes []byte) (*Module, error) {
	reader := bytes.NewReader(moduleBytes) // Generate reader for inputted raw WASM module

	module, err := wasm.ReadModule(reader, nil) // Read module via wagon

	if err != nil { // Check for errors
		return &Module{}, err // Return error
	}

	/*
		err = validate.VerifyModule(module) // Verify module

		if err != nil { // Check for errors
			return &Module{}, err // Return error
		}
	*/

	if err != nil { // Check for errors
		return &Module{}, err // Return error
	}

	functionNames := make(map[int]string) // Init names buffer

	for _, sec := range module.Customs { // Iterate through customs
		if sec.Name == "name" { // Check should be analyzed
			r := bytes.NewReader(sec.RawSection.Bytes) // Get section bytes as byte reader

			for { // Iterate
				ty, err := leb128.ReadVarUint32(r) // Read

				if err != nil || ty != 1 { // Check for errors
					break // Break
				}

				payloadLen, err := leb128.ReadVarUint32(r) // Get payload length

				if err != nil { // Check for errors
					return &Module{}, err // Return errors
				}

				data := make([]byte, int(payloadLen)) // Init payload buffer

				n, err := r.Read(data) // Read data into buffer

				if err != nil { // Check for errors
					return &Module{}, err // Return errors
				}

				if n != len(data) { // Check for invalid length
					return &Module{}, errors.New("len mismatch") // Return errors
				}
				{
					r := bytes.NewReader(data) // Init reader

					for { // Iterate
						count, err := leb128.ReadVarUint32(r) // Read payload count

						if err != nil { // Check for errors
							break // Break
						}

						for i := 0; i < int(count); i++ { // Iterate through
							index, err := leb128.ReadVarUint32(r) // Read index

							if err != nil { // Check for errors
								return &Module{}, err // Return errors
							}

							nameLen, err := leb128.ReadVarUint32(r) // Read name length

							if err != nil { // Check for errors
								return &Module{}, err // Return errors
							}

							name := make([]byte, int(nameLen)) // Init name buffer

							n, err := r.Read(name) // Read into buffer

							if err != nil { // Check for errors
								return &Module{}, err // Return errors
							}

							if n != len(name) { // Check for length mismatch
								return &Module{}, errors.New("len mismatch") // Return errors
							}

							functionNames[int(index)] = string(name) // Append function name
							//fmt.Printf("%d -> %s\n", int(index), string(name))
						}
					}
				}
			}
		}
		//fmt.Printf("%d function names written\n", len(functionNames))
	}

	return &Module{ // Return initialized module
		Base:          module,                   // Set base module
		FunctionNames: functionNames,            // Set function names
		Identifier:    crypto.Sha3(moduleBytes), // Gen, set identifier
	}, nil
}

// String - module stringer
func (module *Module) String() string {
	marshaledVal, _ := json.MarshalIndent(*module, "", "  ") // Marshal to JSON

	return string(marshaledVal) // Return string value
}

// CompileForInterpreter - compile given WASM bytecode for interpreter
func (module *Module) CompileForInterpreter(gp GasPolicy) (_retCode []InterpreterCode, retErr error) {
	defer common.CatchPanic(&retErr) // Catch panic

	ret := make([]InterpreterCode, 0) // Init interpreter code buffer
	importTypeIDs := make([]int, 0)   // Init imports buffer

	if module.Base.Import != nil { // Check has imports
		for i := 0; i < len(module.Base.Import.Entries); i++ { // Iterate through imports
			e := &module.Base.Import.Entries[i] // Get import entry

			if e.Type.Kind() != wasm.ExternalFunction { // Check is extern
				continue // Continue to next import
			}

			tyID := e.Type.(wasm.FuncImport).Type       // Get import type
			ty := &module.Base.Types.Entries[int(tyID)] // Get import entry ty

			buf := &bytes.Buffer{} // Init buffer

			binary.Write(buf, binary.LittleEndian, uint32(1))            // value ID
			binary.Write(buf, binary.LittleEndian, opcodes.InvokeImport) // Write invoked import
			binary.Write(buf, binary.LittleEndian, uint32(i))            // Write to buffer

			binary.Write(buf, binary.LittleEndian, uint32(0)) // Write to buffer

			if len(ty.ReturnTypes) != 0 { // Check has return types
				binary.Write(buf, binary.LittleEndian, opcodes.ReturnValue) // Write return value
				binary.Write(buf, binary.LittleEndian, uint32(1))           // Write buffer
			} else { // No return types
				binary.Write(buf, binary.LittleEndian, opcodes.ReturnVoid) // Write is void
			}

			code := buf.Bytes() // Get final bytes

			ret = append(ret, InterpreterCode{ // Append to interpreter code
				NumRegs:    2,
				NumParams:  len(ty.ParamTypes),
				NumLocals:  0,
				NumReturns: len(ty.ReturnTypes),
				Bytes:      code,
			})

			importTypeIDs = append(importTypeIDs, int(tyID)) // Append to import types
		}
	}

	numFuncImports := len(ret)                                                         // Get # of func imports
	ret = append(ret, make([]InterpreterCode, len(module.Base.FunctionIndexSpace))...) // Append function index space to parsed interpreter source

	for i, f := range module.Base.FunctionIndexSpace { // Iterate thorugh function index space
		//fmt.Printf("Compiling function %d (%+v) with %d locals\n", i, f.Sig, len(f.Body.Locals))
		d, err := disasm.NewDisassembly(f, module.Base) // Disassemble function

		if err != nil { // Check for errors
			panic(err) // Panic to catch
		}

		compiler := NewSSAFunctionCompiler(module.Base, d) // Init compiler
		compiler.CallIndexOffset = numFuncImports          // Set index offset
		compiler.Compile(importTypeIDs)                    // Compile

		if module.DisableFloatingPoint { // Check should disable floats
			compiler.FilterFloatingPoint() // Set filter floating
		}

		if gp != nil { // Check has gas policy
			compiler.InsertGasCounters(gp) // Set gas policy/counter
		}
		//fmt.Println(compiler.Code)
		//fmt.Printf("%+v\n", compiler.NewCFGraph())
		numRegs := compiler.RegAlloc() // Alloc reg
		//fmt.Println(compiler.Code)
		numLocals := 0 // Init local vrs buffer

		for _, v := range f.Body.Locals { // Iterate through locals
			numLocals += int(v.Count) // Increment counter
		}

		ret[numFuncImports+i] = InterpreterCode{ // Set interpreter code
			NumRegs:    numRegs,
			NumParams:  len(f.Sig.ParamTypes),
			NumLocals:  numLocals,
			NumReturns: len(f.Sig.ReturnTypes),
			Bytes:      compiler.Serialize(),
		}
	}

	return ret, nil // Return read/compiled interpreter code
}
