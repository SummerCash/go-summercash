package vm

import (
	"errors"
	"fmt"

	"github.com/SummerCash/ursa/common"
)

var _ ImportResolver = (*NopResolver)(nil)

// NopResolver - nil WebAssembly module import resolver.
type NopResolver struct{}

// Resolver - define imports for WebAssembly modules
type Resolver struct {
	tempRet0 int64
}

// ResolveFunc - panic
func (r *NopResolver) ResolveFunc(module, field string) FunctionImport {
	panic("func import not allowed") // Panic
}

// ResolveGlobal - panic
func (r *NopResolver) ResolveGlobal(module, field string) int64 {
	panic("global import not allowed") // Panic
}

// RunWithGasLimit - run a WebAssembly modules function denoted by its ID with a specified set
// of parameters for a specified amount of instructions (also known as gas) denoted by `limit`
func (vm *VirtualMachine) RunWithGasLimit(entryID, limit int, params ...int64) (int64, error) {
	count := 0 // Init count buffer

	vm.Ignite(entryID, params...) // Ignite vm

	for !vm.Exited { // Do until vm exits
		vm.Execute() // Execute

		if vm.Delegate != nil { // Check for delegate call
			vm.Delegate()     // Run delegate call
			vm.Delegate = nil // Delgate run, set to nil
		}

		count++ // Iterate

		if count >= limit { // Check gas limit exceeded
			return -1, errors.New("gas limit exceeded") // Return error
		}
	}

	if vm.ExitError != nil { // Check for errors
		return -1, common.UnifyError(vm.ExitError) // Return error
	}

	return vm.ReturnValue, nil // Return success
}

// Run runs a WebAssembly modules function denoted by its ID with a specified set
// of parameters.
// Panics on logical errors.
func (vm *VirtualMachine) Run(entryID int, params ...int64) (int64, error) {
	vm.Ignite(entryID, params...) // Ignite VM

	for !vm.Exited { // Check not already exited
		vm.Execute() // Execute

		if vm.Delegate != nil { // Check for errors
			vm.Delegate()     // Run delegate call
			vm.Delegate = nil // Delegate run, set to nil
		}
	}

	if vm.ExitError != nil { // Check for exit error
		return -1, common.UnifyError(vm.ExitError) // Panic/return err
	}

	return vm.ReturnValue, nil // Return success
}

// ResolveFunc - define a set of import functions that may be called within a WebAssembly module
func (r *Resolver) ResolveFunc(module, field string) FunctionImport {
	//fmt.Printf("Resolve func: %s %s\n", module, field) // Log resolve

	switch module { // Handle module types
	case "env": // Env module
		switch field { // Handle fields
		case "__life_ping":
			return func(vm *VirtualMachine) int64 {
				return vm.GetCurrentFrame().Locals[0] + 1
			}
		case "__life_log":
			return func(vm *VirtualMachine) int64 {
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
				msg := vm.Memory[ptr : ptr+msgLen]
				fmt.Printf("[app] %s\n", string(msg))
				return 0
			}

		default:
			panic(fmt.Errorf("unknown field: %s", field)) // Panic
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module)) // Panic
	}
}

// ResolveGlobal - define a set of global variables for use within a WebAssembly module
func (r *Resolver) ResolveGlobal(module, field string) int64 {
	fmt.Printf("Resolve global: %s %s\n", module, field) // Log resolve global

	switch module { // Handle module types
	case "env": // Env module
		switch field { // Handle fields
		case "__ursa_magic":
			return 424 // Return magic
		default:
			panic(fmt.Errorf("unknown field: %s", field)) // Panic
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module)) // Panic
	}
}
