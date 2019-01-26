package vm

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/bits"

	"github.com/SummerCash/ursa/common"
	"github.com/SummerCash/ursa/compiler"
	"github.com/SummerCash/ursa/compiler/opcodes"
	"github.com/SummerCash/wagon/wasm"
)

const (
	// DefaultCallStackSize - default call stack size.
	DefaultCallStackSize = 512

	// DefaultPageSize - linear memory page size
	DefaultPageSize = 65536
)

// FunctionImport represents the function import type. If len(sig.ReturnTypes) == 0, the return value will be ignored.
type FunctionImport func(vm *VirtualMachine) int64

// VirtualMachine - container holding VM config, metadata
type VirtualMachine struct {
	Environment Environment `json:"environment"` // VM config

	Module *compiler.Module // Loaded web-assembly module

	FunctionCode    []compiler.InterpreterCode // Func bytecode
	FunctionImports []FunctionImport           // Import funcs

	CallStack    []Frame // VM call stack
	CurrentFrame int     // Current callstack frame

	Table []uint32 // VM mem/runtime table

	Globals []int64 // Global vrs

	Memory []byte // Virtual machine memory

	NumValueSlots int // Num of used value slots

	Yielded int64 // Did yield

	InsideExecute bool // Inside execute

	Delegate func() // Delegate call

	Exited    bool        // Did exit
	ExitError interface{} // Error on exit

	ReturnValue int64 // Return value

	Gas              uint64 // Gas usage
	GasLimitExceeded bool   // Has exceeded given gas limit

	StateDB *StateDatabase // State database
}

// Frame - call stack frame
type Frame struct {
	FunctionID   int     // Module function identifier
	Code         []byte  // Frame bytecode
	Regs         []int64 // Regs
	Locals       []int64 // Local vrs
	IP           int     // IP
	ReturnReg    int     // Returning reg
	Continuation int32   // Continuation
}

// ImportResolver - interface for allowing one to define imports to WebAssembly modules
type ImportResolver interface {
	ResolveFunc(module, field string) FunctionImport // Func resolver method
	ResolveGlobal(module, field string) int64        // Global resolver method
}

// NewVirtualMachine - instantiate a virtual machine for a given WebAssembly module, with
// specific execution options specified under a VMConfig, and a WebAssembly module import
// resolver
func NewVirtualMachine(code []byte, config Environment, impResolver ImportResolver, gasPolicy compiler.GasPolicy) (_retVM *VirtualMachine, retErr error) {
	if config.EnableJIT { // Check needs JIT
		fmt.Println("Warning: JIT support is removed.") // Log removed JIT support
	}

	m, err := compiler.LoadModule(code) // Load module

	if err != nil { // Check for errors
		return nil, err // Return error
	}

	m.DisableFloatingPoint = config.DisableFloatingPoint // Set floating point disabled

	functionCode, err := m.CompileForInterpreter(gasPolicy) // Compile function code for interpreter

	if err != nil { // Check for errors
		return nil, err // Return error
	}

	defer common.CatchPanic(&retErr) // Catch panic

	table := make([]uint32, 0)               // Init buffer
	globals := make([]int64, 0)              // Init buffer
	funcImports := make([]FunctionImport, 0) // Init buffer

	if m.Base.Import != nil && impResolver != nil { // Check has import/import resolver
		for _, imp := range m.Base.Import.Entries { // Iterate through imports
			switch imp.Type.Kind() { // Handle import types
			case wasm.ExternalFunction: // Check is extern func import
				funcImports = append(funcImports, impResolver.ResolveFunc(imp.ModuleName, imp.FieldName)) // Append to func imports
			case wasm.ExternalGlobal: // Check is extern global import
				globals = append(globals, impResolver.ResolveGlobal(imp.ModuleName, imp.FieldName)) // Handle
			case wasm.ExternalMemory:
				// TODO: Do we want a real import?
				if m.Base.Memory != nil && len(m.Base.Memory.Entries) > 0 { // Check mem not nil
					panic("cannot import another memory while we already have one") // Panic
				}

				m.Base.Memory = &wasm.SectionMemories{ // Set memory
					Entries: []wasm.Memory{
						wasm.Memory{
							Limits: wasm.ResizableLimits{
								Initial: uint32(config.DefaultMemoryPages),
							},
						},
					},
				}
			case wasm.ExternalTable:
				// TODO: Do we want a real import?
				if m.Base.Table != nil && len(m.Base.Table.Entries) > 0 { // Check table not empty
					panic("cannot import another table while we already have one")
				}

				m.Base.Table = &wasm.SectionTables{ // Set table
					Entries: []wasm.Table{
						wasm.Table{
							Limits: wasm.ResizableLimits{
								Initial: uint32(config.DefaultTableSize),
							},
						},
					},
				}
			default:
				panic(fmt.Errorf("import kind not supported: %d", imp.Type.Kind())) // Panic
			}
		}
	}

	for _, entry := range m.Base.GlobalIndexSpace { // Load global entries
		globals = append(globals, execInitExpr(entry.Init, globals)) // Append global entry
	}

	if m.Base.Table != nil && len(m.Base.Table.Entries) > 0 { // Populate table elements
		t := &m.Base.Table.Entries[0] // Set table entry buffer

		if config.MaxTableSize != 0 && int(t.Limits.Initial) > config.MaxTableSize { // Check table size exceeded
			panic("max table size exceeded") // Panic
		}

		table = make([]uint32, int(t.Limits.Initial)) // Init table buffer

		for i := 0; i < int(t.Limits.Initial); i++ { // Iterate
			table[i] = 0xffffffff // Set table entry
		}

		if m.Base.Elements != nil && len(m.Base.Elements.Entries) > 0 { // Check elements not nil
			for _, e := range m.Base.Elements.Entries { // Iterate through entries
				offset := int(execInitExpr(e.Offset, globals)) // Get offset

				copy(table[offset:], e.Elems) // Copy
			}
		}
	}

	memory := make([]byte, 0)                                   // Init memory buffer
	if m.Base.Memory != nil && len(m.Base.Memory.Entries) > 0 { // Check base memory not nil
		initialLimit := int(m.Base.Memory.Entries[0].Limits.Initial) // Init initial limit

		if config.MaxMemoryPages != 0 && initialLimit > config.MaxMemoryPages { // Check max memory exceeded
			panic("max memory exceeded") // Panic
		}

		capacity := initialLimit * DefaultPageSize // Get capacity

		memory = make([]byte, capacity) // Init empty memory

		for i := 0; i < capacity; i++ { // Iterate until capacity
			memory[i] = 0 // Set empty mem
		}

		if m.Base.Data != nil && len(m.Base.Data.Entries) > 0 { // Iterate through entries
			for _, e := range m.Base.Data.Entries { // Iterate through entires
				offset := int(execInitExpr(e.Offset, globals)) // Get offset

				copy(memory[int(offset):], e.Data) // Copy
			}
		}
	}

	vm := &VirtualMachine{
		Module:          m,
		Environment:     config,
		FunctionCode:    functionCode,
		FunctionImports: funcImports,
		CallStack:       make([]Frame, DefaultCallStackSize),
		CurrentFrame:    -1,
		Table:           table,
		Globals:         globals,
		Memory:          memory,
		Exited:          true,
	} // Init VM

	rootState := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, 0) // Init state entry

	stateDB := NewStateDatabase(rootState) // Init state database

	err = stateDB.WriteToMemory() // Write state db

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	(*vm).StateDB = stateDB // Set state DB

	return vm, nil // Return init vm
}

// SaveState - save state
func (vm *VirtualMachine) SaveState() error {
	workingRoot := vm.StateDB.WorkingRoot // Get working root

	nonce := workingRoot.Nonce + 1 // Init nonce buffer

	maxChild, err := workingRoot.FindMax() // Find max

	if err == nil { // Check no errors
		nonce = maxChild.Nonce + 1 // Set nonce
	}

	state := NewStateEntry(vm.CallStack, vm.CurrentFrame, vm.Table, vm.Globals, vm.Memory, vm.NumValueSlots, vm.Yielded, vm.InsideExecute, vm.Exited, vm.ExitError, vm.ReturnValue, vm.Gas, vm.GasLimitExceeded, nonce) // Init state entry

	err = vm.StateDB.AddStateEntry(state, workingRoot) // Add state entry

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// ResetToState - revert vm state to given state with ID
func (vm *VirtualMachine) ResetToState(id []byte) error {
	err := vm.LoadStateDB(hex.EncodeToString(vm.StateDB.ID)) // Load state db

	if err != nil { // Check for errors
		return err // Return found error
	}

	state, err := vm.StateDB.QueryState(id) // Query state

	if err != nil { // Check for errors
		return err // Return found error
	}

	vm.StateDB.SetWorkingRoot(state) // Set working root

	(*vm).CallStack = state.State.CallStack               // Set call stack
	(*vm).CurrentFrame = state.State.CurrentFrame         // Set current frame
	(*vm).Table = state.State.Table                       // Set table
	(*vm).Globals = state.State.Globals                   // Set globals
	(*vm).Memory = state.State.Memory                     // Set memory
	(*vm).NumValueSlots = state.State.NumValueSlots       // Set # value slots
	(*vm).Yielded = state.State.Yielded                   // Set yielded
	(*vm).InsideExecute = state.State.InsideExecute       // Set inside execute
	(*vm).Exited = state.State.Exited                     // Set has exited
	(*vm).ExitError = state.State.ExitError               // Set exit error
	(*vm).ReturnValue = state.State.ReturnValue           // Set return val
	(*vm).Gas = state.State.Gas                           // Set gas
	(*vm).GasLimitExceeded = state.State.GasLimitExceeded // Set has exceeded gas limit

	return nil // No error occurred, return nil
}

// LoadWorkingRoot - load last saved state
func (vm *VirtualMachine) LoadWorkingRoot() error {
	err := vm.LoadStateDB(hex.EncodeToString(vm.StateDB.ID)) // Load state db

	if err != nil { // Check for errors
		return err // Return found error
	}

	(*vm).CallStack = vm.StateDB.WorkingRoot.State.CallStack               // Set call stack
	(*vm).CurrentFrame = vm.StateDB.WorkingRoot.State.CurrentFrame         // Set current frame
	(*vm).Table = vm.StateDB.WorkingRoot.State.Table                       // Set table
	(*vm).Globals = vm.StateDB.WorkingRoot.State.Globals                   // Set globals
	(*vm).Memory = vm.StateDB.WorkingRoot.State.Memory                     // Set memory
	(*vm).NumValueSlots = vm.StateDB.WorkingRoot.State.NumValueSlots       // Set # value slots
	(*vm).Yielded = vm.StateDB.WorkingRoot.State.Yielded                   // Set yielded
	(*vm).InsideExecute = vm.StateDB.WorkingRoot.State.InsideExecute       // Set inside execute
	(*vm).Exited = vm.StateDB.WorkingRoot.State.Exited                     // Set has exited
	(*vm).ExitError = vm.StateDB.WorkingRoot.State.ExitError               // Set exit error
	(*vm).ReturnValue = vm.StateDB.WorkingRoot.State.ReturnValue           // Set return val
	(*vm).Gas = vm.StateDB.WorkingRoot.State.Gas                           // Set gas
	(*vm).GasLimitExceeded = vm.StateDB.WorkingRoot.State.GasLimitExceeded // Set has exceeded gas limit

	return nil // No error occurred, return nil
}

// LoadStateDB - load state database
func (vm *VirtualMachine) LoadStateDB(id string) error {
	stateDB, err := ReadStateDBFromMemory(id) // Read state db

	if err != nil { // Check for errors
		return err // Return found error
	}

	(*vm).StateDB = stateDB // Set state db

	return nil // No error occurred, return nil
}

// Init - initializes a frame; must be called on `call` and `call_indirect`
func (f *Frame) Init(vm *VirtualMachine, functionID int, code compiler.InterpreterCode) {
	numValueSlots := code.NumRegs + code.NumParams + code.NumLocals // Get num slots

	if vm.Environment.MaxValueSlots != 0 && vm.NumValueSlots+numValueSlots > vm.Environment.MaxValueSlots { // Check for max count exceeded
		panic("max value slot count exceeded") // Panic
	}

	vm.NumValueSlots += numValueSlots // Set num value slots

	values := make([]int64, numValueSlots) // Init value buffer

	f.FunctionID = functionID        // Set frame function ID
	f.Regs = values[:code.NumRegs]   // Set frame regs
	f.Locals = values[code.NumRegs:] // Set frame locals
	f.Code = code.Bytes              // Set frame bytecode
	f.IP = 0                         // Set frame IP
	f.Continuation = 0               // Set frame continuation

	//fmt.Printf("Enter function %d (%s)\n", functionID, vm.Module.FunctionNames[functionID])
}

// Destroy - destroy a frame; must be called on return
func (f *Frame) Destroy(vm *VirtualMachine) {
	numValueSlots := len(f.Regs) + len(f.Locals) // Get value slots
	vm.NumValueSlots -= numValueSlots            // Destroy frame

	//fmt.Printf("Leave function %d (%s)\n", f.FunctionID, vm.Module.FunctionNames[f.FunctionID])
}

// GetCurrentFrame - return the current frame
func (vm *VirtualMachine) GetCurrentFrame() *Frame {
	if vm.Environment.MaxCallStackDepth != 0 && vm.CurrentFrame >= vm.Environment.MaxCallStackDepth { // Check for stack limit exceeded
		panic("max call stack depth exceeded") // Panic
	}

	if vm.CurrentFrame >= len(vm.CallStack) { // Check for stack overflow ( ͡° ͜ʖ ͡°)
		panic("call stack overflow")
		//vm.CallStack = append(vm.CallStack, make([]Frame, DefaultCallStackSize / 2)...)
	}

	return &vm.CallStack[vm.CurrentFrame] // Return frame
}

// getExport - get export with given key
func (vm *VirtualMachine) getExport(key string, kind wasm.External) (int, bool) {
	if vm.Module.Base.Export == nil { // Check exports nil
		return -1, false // Return does not exist
	}

	entry, ok := vm.Module.Base.Export.Entries[key] // Get entry

	if !ok { // Not "ok"
		return -1, false // Return does not exist
	}

	if entry.Kind != kind { // Check kind does not match
		return -1, false // Return does not exist
	}

	return int(entry.Index), true // Return index of export entry
}

// GetGlobalExport - return the global export with the given name
func (vm *VirtualMachine) GetGlobalExport(key string) (int, bool) {
	return vm.getExport(key, wasm.ExternalGlobal) // Return export
}

// GetFunctionExport - return the function export with the given name
func (vm *VirtualMachine) GetFunctionExport(key string) (int, bool) {
	return vm.getExport(key, wasm.ExternalFunction) // Return export
}

// PrintStackTrace - print the entire VM stack trace for debugging
func (vm *VirtualMachine) PrintStackTrace() {
	fmt.Println("--- Begin stack trace ---") // Log begin stack trace

	for i := vm.CurrentFrame; i >= 0; i-- { // Iterate through frames
		functionID := vm.CallStack[i].FunctionID                                         // Get function ID
		fmt.Printf("<%d> [%d] %s\n", i, functionID, vm.Module.FunctionNames[functionID]) // Log
	}

	fmt.Println("--- End stack trace ---") // Log end stack trace
}

// Ignite - initialize the first call frame
func (vm *VirtualMachine) Ignite(functionID int, params ...int64) {
	if vm.ExitError != nil { // Check for exit error
		panic("last execution exited with error; cannot ignite.") // Panic
	}

	if vm.CurrentFrame != -1 { // Check call stack not empty
		panic("call stack not empty; cannot ignite.") // Panic
	}

	code := vm.FunctionCode[functionID] // Get bytecode

	if code.NumParams != len(params) { // Check param mismatch
		panic("param count mismatch") // Panic
	}

	vm.Exited = false // Set exited

	vm.CurrentFrame++ // Increment current frame

	frame := vm.GetCurrentFrame() // Get current frame
	frame.Init(                   // Initialize frame
		vm,
		functionID,
		code,
	)

	copy(frame.Locals, params) // Copy params to frame locals
}

// AddAndCheckGas - add and check gas
func (vm *VirtualMachine) AddAndCheckGas(delta uint64) bool {
	newGas := vm.Gas + delta // Calculate gas

	if newGas < vm.Gas { // Check for gas overflow
		panic("gas overflow") // Panic
	}

	if vm.Environment.GasLimit != 0 && newGas > vm.Environment.GasLimit { // Check gas limit exceeded
		if vm.Environment.ReturnOnGasLimitExceeded { // Check should return
			return false // Return
		}

		panic("gas limit exceeded") // Panic
	}

	vm.Gas = newGas // Set gas

	return true // Return success
}

// Execute - start the virtual machines main instruction processing loop.
// May return at any point and is guaranteed to return
// at least once every 10000 instructions. Caller is responsible for
// detecting VM status in a loop
func (vm *VirtualMachine) Execute() {
	if vm.Exited == true { // Check already exited
		panic("attempting to execute an exited vm") // Panic
	}

	if vm.Delegate != nil { // Check delegate not cleared
		panic("delegate not cleared") // Panic
	}

	if vm.InsideExecute { // Check is not re-entrant
		panic("vm execution is not re-entrant") // Panic
	}

	vm.InsideExecute = true     // Set inside execute
	vm.GasLimitExceeded = false // Set gas limit exceeded

	defer func() {
		vm.InsideExecute = false // Set inside execute

		if err := recover(); err != nil { // Check for errors
			vm.Exited = true   // Set exited
			vm.ExitError = err // Set exit error
		}
	}()

	frame := vm.GetCurrentFrame() // Get current frame

	for { // Iterate
		valueID := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])) // Init valueID
		ins := opcodes.Opcode(frame.Code[frame.IP+4])                                 // Init instruction
		frame.IP += 5                                                                 // Set frame IP

		//fmt.Printf("INS: [%d] %s\n", valueID, ins.String())

		switch ins { // Handle different opcodes
		case opcodes.Nop: // Handle Nop
		case opcodes.Unreachable: // Handle Unreachable
			panic("wasm: unreachable executed")
		case opcodes.Select: // Handle Select
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			c := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])
			frame.IP += 12
			if c != 0 {
				frame.Regs[valueID] = a
			} else {
				frame.Regs[valueID] = b
			}
		case opcodes.I32Const: // Handle I32Const
			val := binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			frame.IP += 4
			frame.Regs[valueID] = int64(val)
		case opcodes.I32Add: // Handle I32Add
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			frame.Regs[valueID] = int64(a + b)
		case opcodes.I32Sub: // Handle I32Sub
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			frame.Regs[valueID] = int64(a - b)
		case opcodes.I32Mul: // Handle I32Mul
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			frame.Regs[valueID] = int64(a * b)
		case opcodes.I32DivS: // Handle I32DivS
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			if b == 0 {
				panic("integer division by zero")
			}

			if a == math.MinInt32 && b == -1 {
				panic("signed integer overflow")
			}

			frame.IP += 8
			frame.Regs[valueID] = int64(a / b)
		case opcodes.I32DivU: // Handle I32DivU
			a := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			if b == 0 {
				panic("integer division by zero")
			}

			frame.IP += 8
			frame.Regs[valueID] = int64(a / b)
		case opcodes.I32RemS: // Handle I32RemS
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			if b == 0 {
				panic("integer division by zero")
			}

			frame.IP += 8
			frame.Regs[valueID] = int64(a % b)
		case opcodes.I32RemU: // Handle I32RemU
			a := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			if b == 0 {
				panic("integer division by zero")
			}

			frame.IP += 8
			frame.Regs[valueID] = int64(a % b)
		case opcodes.I32And: // Handle I32And
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(a & b)
		case opcodes.I32Or: // Handle I32Or
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(a | b)
		case opcodes.I32Xor: // Handle I32Xor
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(a ^ b)
		case opcodes.I32Shl: // Handle I32Shl
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(a << (b % 32))
		case opcodes.I32ShrS: // Handle I32ShrS
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(a >> (b % 32))
		case opcodes.I32ShrU: // Handle I32ShrU
			a := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(a >> (b % 32))
		case opcodes.I32Rotl: // Handle I32Rotl
			a := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(bits.RotateLeft32(a, int(b)))
		case opcodes.I32Rotr: // Handle I32Rotr
			a := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(bits.RotateLeft32(a, -int(b)))
		case opcodes.I32Clz: // Handle I32Clz
			val := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])

			frame.IP += 4
			frame.Regs[valueID] = int64(bits.LeadingZeros32(val))
		case opcodes.I32Ctz: // Handle I32Ctz
			val := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])

			frame.IP += 4
			frame.Regs[valueID] = int64(bits.TrailingZeros32(val))
		case opcodes.I32PopCnt: // Handle I32PopCnt
			val := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])

			frame.IP += 4
			frame.Regs[valueID] = int64(bits.OnesCount32(val))
		case opcodes.I32EqZ: // Handle I32EqZ
			val := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])

			frame.IP += 4
			if val == 0 {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32Eq: // Handle I32Eq
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a == b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32Ne: // Handle I32Ne
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a != b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32LtS: // Handle I32LtS
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a < b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32LtU: // Handle I32LtU
			a := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a < b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32LeS: // Handle I32LeS
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a <= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32LeU: // Handle I32LeU
			a := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a <= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32GtS: // Handle I32GtS
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a > b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32GtU: // Handle I32GtU
			a := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a > b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32GeS: // Handle I32GeS
			a := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a >= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I32GeU: // Handle I32GeU
			a := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a >= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64Const: // Handle I64Const
			val := binary.LittleEndian.Uint64(frame.Code[frame.IP : frame.IP+8])
			frame.IP += 8
			frame.Regs[valueID] = int64(val)
		case opcodes.I64Add: // Handle I64Add
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			frame.Regs[valueID] = a + b
		case opcodes.I64Sub: // Handle I64Sub
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			frame.Regs[valueID] = a - b
		case opcodes.I64Mul: // Handle I64Mul
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			frame.Regs[valueID] = a * b
		case opcodes.I64DivS: // Handle I64DivS
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]

			if b == 0 {
				panic("integer division by zero")
			}

			if a == math.MinInt64 && b == -1 {
				panic("signed integer overflow")
			}

			frame.IP += 8
			frame.Regs[valueID] = a / b
		case opcodes.I64DivU: // Handle I64DivU
			a := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			if b == 0 {
				panic("integer division by zero")
			}

			frame.IP += 8
			frame.Regs[valueID] = int64(a / b)
		case opcodes.I64RemS: // Handle I64RemS
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]

			if b == 0 {
				panic("integer division by zero")
			}

			frame.IP += 8
			frame.Regs[valueID] = a % b
		case opcodes.I64RemU: // Handle I64RemU
			a := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			if b == 0 {
				panic("integer division by zero")
			}

			frame.IP += 8
			frame.Regs[valueID] = int64(a % b)
		case opcodes.I64And: // Handle I64And
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]

			frame.IP += 8
			frame.Regs[valueID] = a & b
		case opcodes.I64Or: // Handle I64Or
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]

			frame.IP += 8
			frame.Regs[valueID] = a | b
		case opcodes.I64Xor: // Handle I64Xor
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]

			frame.IP += 8
			frame.Regs[valueID] = a ^ b
		case opcodes.I64Shl: // Handle I64Shl
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = a << (b % 64)
		case opcodes.I64ShrS: // Handle I64ShrS
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = a >> (b % 64)
		case opcodes.I64ShrU: // Handle I64ShrU
			a := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(a >> (b % 64))
		case opcodes.I64Rotl: // Handle I64Rotl
			a := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(bits.RotateLeft64(a, int(b)))
		case opcodes.I64Rotr: // Handle I64Rotr
			a := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])

			frame.IP += 8
			frame.Regs[valueID] = int64(bits.RotateLeft64(a, -int(b)))
		case opcodes.I64Clz: // Handle I64Clz
			val := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])

			frame.IP += 4
			frame.Regs[valueID] = int64(bits.LeadingZeros64(val))
		case opcodes.I64Ctz: // Handle I64Ctz
			val := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])

			frame.IP += 4
			frame.Regs[valueID] = int64(bits.TrailingZeros64(val))
		case opcodes.I64PopCnt: // Handle I64PopCnt
			val := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])

			frame.IP += 4
			frame.Regs[valueID] = int64(bits.OnesCount64(val))
		case opcodes.I64EqZ: // Handle I64EqZ
			val := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])

			frame.IP += 4
			if val == 0 {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64Eq: // Handle I64Eq
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			if a == b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64Ne: // Handle I64Ne
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			if a != b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64LtS: // Handle I64LtS
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			if a < b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64LtU: // Handle I64LtU
			a := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a < b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64LeS: // Handle I64LeS
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			if a <= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64LeU: // Handle I64LeU
			a := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a <= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64GtS: // Handle I64GtS
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			if a > b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64GtU: // Handle I64GtU
			a := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a > b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64GeS: // Handle I64GeS
			a := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			b := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			if a >= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.I64GeU: // Handle I64GeU
			a := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			b := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))])
			frame.IP += 8
			if a >= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F32Add: // Handle F32Add
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float32bits(a + b))
		case opcodes.F32Sub: // Handle F32Sub
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float32bits(a - b))
		case opcodes.F32Mul: // Handle F32Mul
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float32bits(a * b))
		case opcodes.F32Div: // Handle F32Div
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float32bits(a / b))
		case opcodes.F32Sqrt: // Handle F32Sqrt
			val := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(math.Sqrt(float64(val)))))
		case opcodes.F32Min: // Handle F32Min
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float32bits(float32(math.Min(float64(a), float64(b)))))
		case opcodes.F32Max: // Handle F32Max
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float32bits(float32(math.Max(float64(a), float64(b)))))
		case opcodes.F32Ceil: // Handle F32Ceil
			val := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(math.Ceil(float64(val)))))
		case opcodes.F32Floor: // Handle F32Floor
			val := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(math.Floor(float64(val)))))
		case opcodes.F32Trunc: // Handle F32Trunc
			val := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(math.Trunc(float64(val)))))
		case opcodes.F32Nearest: // Handle F32Nearest
			val := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(math.RoundToEven(float64(val)))))
		case opcodes.F32Abs: // Handle F32Abs
			val := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(math.Abs(float64(val)))))
		case opcodes.F32Neg: // Handle F32Neg
			val := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(-val))
		case opcodes.F32CopySign: // Handle F32CopySign
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float32bits(float32(math.Copysign(float64(a), float64(b)))))
		case opcodes.F32Eq: // Handle F32Eq
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a == b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F32Ne: // Handle F32Ne
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a != b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F32Lt: // Handle F32Lt
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a < b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F32Le: // Handle F32Le
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a <= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F32Gt: // Handle F32Gt
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a > b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F32Ge: // Handle F32Ge
			a := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a >= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F64Add: // Handle F64Add
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float64bits(a + b))
		case opcodes.F64Sub: // Handle F64Sub
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float64bits(a - b))
		case opcodes.F64Mul: // Handle F64Mul
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float64bits(a * b))
		case opcodes.F64Div: // Handle F64Div
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float64bits(a / b))
		case opcodes.F64Sqrt: // Handle F64Sqrt
			val := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(math.Sqrt(val)))
		case opcodes.F64Min: // Handle F64Min
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float64bits(math.Min(a, b)))
		case opcodes.F64Max: // Handle F64Max
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float64bits(math.Max(a, b)))
		case opcodes.F64Ceil: // Handle F64Ceil
			val := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(math.Ceil(val)))
		case opcodes.F64Floor: // Handle F64Floor
			val := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(math.Floor(val)))
		case opcodes.F64Trunc: // Handle F64Trunc
			val := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(math.Trunc(val)))
		case opcodes.F64Nearest: // Handle F64Nearest
			val := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(math.RoundToEven(val)))
		case opcodes.F64Abs: // Handle F64Abs
			val := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(math.Abs(val)))
		case opcodes.F64Neg: // Handle F64Neg
			val := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(-val))
		case opcodes.F64CopySign: // Handle F64CopySign
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			frame.Regs[valueID] = int64(math.Float64bits(math.Copysign(a, b)))
		case opcodes.F64Eq: // Handle F64Eq
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a == b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F64Ne: // Handle F64Ne
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a != b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F64Lt: // Handle F64Lt
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a < b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F64Le: // Handle F64Le
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a <= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F64Gt: // Handle F64Gt
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a > b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}
		case opcodes.F64Ge: // Handle F64Ge
			a := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			b := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]))
			frame.IP += 8
			if a >= b {
				frame.Regs[valueID] = 1
			} else {
				frame.Regs[valueID] = 0
			}

		case opcodes.I32WrapI64: // Handle I32WrapI64
			v := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(v)

		case opcodes.I32TruncSF32, opcodes.I32TruncUF32: // Handle I32TruncUF32
			v := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(int32(math.Trunc(float64(v))))

		case opcodes.I32TruncSF64, opcodes.I32TruncUF64: // Handle I32TruncUF64
			v := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(int32(math.Trunc(v)))

		case opcodes.I64TruncSF32, opcodes.I64TruncUF32: // Handle I64TruncUF32
			v := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Trunc(float64(v)))

		case opcodes.I64TruncSF64, opcodes.I64TruncUF64: // Handle I64TruncUF64
			v := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Trunc(v))

		case opcodes.F32DemoteF64: // Handle F32DemoteF64
			v := math.Float64frombits(uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(v)))

		case opcodes.F64PromoteF32: // Handle F64PromoteF32
			v := math.Float32frombits(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(float64(v)))

		case opcodes.F32ConvertSI32: // Handle F32ConvertSI32
			v := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(v)))

		case opcodes.F32ConvertUI32: // Handle F32ConvertUI32
			v := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(v)))

		case opcodes.F32ConvertSI64: // Handle F32ConvertSI64
			v := int64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(v)))

		case opcodes.F32ConvertUI64: // Handle F32ConvertUI64
			v := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float32bits(float32(v)))

		case opcodes.F64ConvertSI32: // Handle F64ConvertSI32
			v := int32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(int32(math.Float64bits(float64(v))))

		case opcodes.F64ConvertUI32: // Handle F64ConvertUI32
			v := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(int32(math.Float64bits(float64(v))))

		case opcodes.F64ConvertSI64: // Handle F64ConvertSI64
			v := int64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(float64(v)))

		case opcodes.F64ConvertUI64: // Handle F64ConvertUI64
			v := uint64(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(math.Float64bits(float64(v)))

		case opcodes.I64ExtendUI32: // Handle I64ExtendUI32
			v := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))])
			frame.IP += 4
			frame.Regs[valueID] = int64(v)

		case opcodes.I64ExtendSI32: // Handle I64ExtendSI32
			v := int32(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4
			frame.Regs[valueID] = int64(v)

		case opcodes.I32Load, opcodes.I64Load32U: // Handle I64Load32U
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			frame.IP += 12

			effective := int(uint64(base) + uint64(offset))
			frame.Regs[valueID] = int64(uint32(binary.LittleEndian.Uint32(vm.Memory[effective : effective+4])))
		case opcodes.I64Load32S: // Handle I64Load32S
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			frame.IP += 12

			effective := int(uint64(base) + uint64(offset))
			frame.Regs[valueID] = int64(int32(binary.LittleEndian.Uint32(vm.Memory[effective : effective+4])))
		case opcodes.I64Load: // Handle I64Load
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			frame.IP += 12

			effective := int(uint64(base) + uint64(offset))
			frame.Regs[valueID] = int64(binary.LittleEndian.Uint64(vm.Memory[effective : effective+8]))
		case opcodes.I32Load8S, opcodes.I64Load8S: // Handle I64Load8S
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			frame.IP += 12

			effective := int(uint64(base) + uint64(offset))
			frame.Regs[valueID] = int64(int8(vm.Memory[effective]))
		case opcodes.I32Load8U, opcodes.I64Load8U: // Handle I64Load8U
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			frame.IP += 12

			effective := int(uint64(base) + uint64(offset))
			frame.Regs[valueID] = int64(uint8(vm.Memory[effective]))
		case opcodes.I32Load16S, opcodes.I64Load16S: // Handle I64Load16S
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			frame.IP += 12

			effective := int(uint64(base) + uint64(offset))
			frame.Regs[valueID] = int64(int16(binary.LittleEndian.Uint16(vm.Memory[effective : effective+2])))
		case opcodes.I32Load16U, opcodes.I64Load16U: // Handle I64Load16U
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			frame.IP += 12

			effective := int(uint64(base) + uint64(offset))
			frame.Regs[valueID] = int64(uint16(binary.LittleEndian.Uint16(vm.Memory[effective : effective+2])))
		case opcodes.I32Store, opcodes.I64Store32: // Handle I64Store32
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			value := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+12:frame.IP+16]))]

			frame.IP += 16

			effective := int(uint64(base) + uint64(offset))
			binary.LittleEndian.PutUint32(vm.Memory[effective:effective+4], uint32(value))
		case opcodes.I64Store: // Handle I64Store
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			value := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+12:frame.IP+16]))]

			frame.IP += 16

			effective := int(uint64(base) + uint64(offset))
			binary.LittleEndian.PutUint64(vm.Memory[effective:effective+8], uint64(value))
		case opcodes.I32Store8, opcodes.I64Store8: // Handle I64Store8
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			value := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+12:frame.IP+16]))]

			frame.IP += 16

			effective := int(uint64(base) + uint64(offset))
			vm.Memory[effective] = byte(value)
		case opcodes.I32Store16, opcodes.I64Store16: // Handle I64Store16
			binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4])
			offset := binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8])
			base := uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8:frame.IP+12]))])

			value := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+12:frame.IP+16]))]

			frame.IP += 16

			effective := int(uint64(base) + uint64(offset))
			binary.LittleEndian.PutUint16(vm.Memory[effective:effective+2], uint16(value))

		case opcodes.Jmp: // Handle Jmp
			target := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			vm.Yielded = frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP = target
		case opcodes.JmpEither: // Handle JmpEither
			targetA := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			targetB := int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8]))
			cond := int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8 : frame.IP+12]))
			yieldedReg := int(binary.LittleEndian.Uint32(frame.Code[frame.IP+12 : frame.IP+16]))
			frame.IP += 16

			vm.Yielded = frame.Regs[yieldedReg]
			if frame.Regs[cond] != 0 {
				frame.IP = targetA
			} else {
				frame.IP = targetB
			}
		case opcodes.JmpIf: // Handle JmpIf
			target := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			cond := int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4 : frame.IP+8]))
			yieldedReg := int(binary.LittleEndian.Uint32(frame.Code[frame.IP+8 : frame.IP+12]))
			frame.IP += 12
			if frame.Regs[cond] != 0 {
				vm.Yielded = frame.Regs[yieldedReg]
				frame.IP = target
			}
		case opcodes.JmpTable: // Handle JmpTable
			targetCount := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			frame.IP += 4

			targetsRaw := frame.Code[frame.IP : frame.IP+4*targetCount]
			frame.IP += 4 * targetCount

			defaultTarget := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			frame.IP += 4

			cond := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			frame.IP += 4

			vm.Yielded = frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			frame.IP += 4

			val := int(frame.Regs[cond])
			if val >= 0 && val < targetCount {
				frame.IP = int(binary.LittleEndian.Uint32(targetsRaw[val*4 : val*4+4]))
			} else {
				frame.IP = defaultTarget
			}
		case opcodes.ReturnValue: // Handle ReturnValue
			val := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			frame.Destroy(vm)
			vm.CurrentFrame--
			if vm.CurrentFrame == -1 {
				vm.Exited = true
				vm.ReturnValue = val
				return
			}

			frame = vm.GetCurrentFrame()
			frame.Regs[frame.ReturnReg] = val
		case opcodes.ReturnVoid: // Handle ReturnVoid
			frame.Destroy(vm)
			vm.CurrentFrame--
			if vm.CurrentFrame == -1 {
				vm.Exited = true
				vm.ReturnValue = 0
				return
			}

			frame = vm.GetCurrentFrame()
		case opcodes.GetLocal: // Handle GetLocal
			id := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			val := frame.Locals[id]
			frame.IP += 4
			frame.Regs[valueID] = val
			//fmt.Printf("GetLocal %d = %d\n", id, val)
		case opcodes.SetLocal: // Handle SetLocal
			id := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			val := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8
			frame.Locals[id] = val
			//fmt.Printf("SetLocal %d = %d\n", id, val)
		case opcodes.GetGlobal: // Handle GetGlobal
			frame.Regs[valueID] = vm.Globals[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			frame.IP += 4
		case opcodes.SetGlobal: // Handle SetGlobal
			id := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			val := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP+4:frame.IP+8]))]
			frame.IP += 8

			vm.Globals[id] = val
		case opcodes.Call: // Handle Call
			functionID := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			frame.IP += 4
			argCount := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			frame.IP += 4
			argsRaw := frame.Code[frame.IP : frame.IP+4*argCount]
			frame.IP += 4 * argCount

			oldRegs := frame.Regs
			frame.ReturnReg = valueID

			vm.CurrentFrame++
			frame = vm.GetCurrentFrame()
			frame.Init(vm, functionID, vm.FunctionCode[functionID])
			for i := 0; i < argCount; i++ {
				frame.Locals[i] = oldRegs[int(binary.LittleEndian.Uint32(argsRaw[i*4:i*4+4]))]
			}
			//fmt.Println("Call params =", frame.Locals[:argCount])

		case opcodes.CallIndirect: // Handle CallIndirect
			typeID := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			frame.IP += 4
			argCount := int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4])) - 1
			frame.IP += 4
			argsRaw := frame.Code[frame.IP : frame.IP+4*argCount]
			frame.IP += 4 * argCount
			tableItemID := frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]
			frame.IP += 4

			sig := &vm.Module.Base.Types.Entries[typeID]

			functionID := int(vm.Table[tableItemID])
			code := vm.FunctionCode[functionID]

			// TODO: We are only checking CC here; Do we want strict type-check?
			if code.NumParams != len(sig.ParamTypes) || code.NumReturns != len(sig.ReturnTypes) {
				panic("type mismatch")
			}

			oldRegs := frame.Regs
			frame.ReturnReg = valueID

			vm.CurrentFrame++
			frame = vm.GetCurrentFrame()
			frame.Init(vm, functionID, code)
			for i := 0; i < argCount; i++ {
				frame.Locals[i] = oldRegs[int(binary.LittleEndian.Uint32(argsRaw[i*4:i*4+4]))]
			}

		case opcodes.InvokeImport: // Handle InvokeImport
			importID := int(binary.LittleEndian.Uint32(frame.Code[frame.IP : frame.IP+4]))
			frame.IP += 4
			vm.Delegate = func() {
				frame.Regs[valueID] = vm.FunctionImports[importID](vm)
			}
			return

		case opcodes.CurrentMemory: // Handle CurrentMemory
			frame.Regs[valueID] = int64(len(vm.Memory) / DefaultPageSize)

		case opcodes.GrowMemory: // Handle GrowMemory
			n := int(uint32(frame.Regs[int(binary.LittleEndian.Uint32(frame.Code[frame.IP:frame.IP+4]))]))
			frame.IP += 4

			current := len(vm.Memory) / DefaultPageSize
			if vm.Environment.MaxMemoryPages == 0 || (current+n >= current && current+n <= vm.Environment.MaxMemoryPages) {
				frame.Regs[valueID] = int64(current)
				vm.Memory = append(vm.Memory, make([]byte, n*DefaultPageSize)...)
			} else {
				frame.Regs[valueID] = -1
			}

		case opcodes.Phi: // Handle Phi
			frame.Regs[valueID] = vm.Yielded

		case opcodes.AddGas: // Handle AddGas
			delta := binary.LittleEndian.Uint64(frame.Code[frame.IP : frame.IP+8])
			frame.IP += 8
			if !vm.AddAndCheckGas(delta) {
				vm.GasLimitExceeded = true
				return
			}

		case opcodes.FPDisabledError: // Handle FPDisabledError
			panic("wasm: floating point disabled") // Panic

		default:
			panic("unknown instruction") // Panic
		}
	}
}
