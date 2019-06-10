package types

import (
	"fmt"

	"github.com/SummerCash/ursa/vm"
)

// workingVMTransaction is the working virtual machine transaction.
// Useful for tx state logging.
var workingVMTransaction *Transaction

// TransactionMetaResolver outlines the default go-summercash WASM tx meta resolver.
type TransactionMetaResolver struct {
	tempRet0 int64
}

/* BEGIN EXPORTED METHODS */

// ResolveFunc defines a set of import functions that may be called within a WebAssembly module.
func (r *TransactionMetaResolver) ResolveFunc(module, field string) vm.FunctionImport {
	switch module { // Handle module types
	case "env": // Env module
		switch field { // Handle fields
		case "__ursa_ping":
			return func(vm *vm.VirtualMachine) int64 {
				return vm.GetCurrentFrame().Locals[0] + 1
			}
		case "__log":
			return func(vm *vm.VirtualMachine) int64 {
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
				msg := vm.Memory[ptr : ptr+msgLen]

				(*workingVMTransaction).Logs = append((*workingVMTransaction).Logs, NewLog("message", msg, Custom)) // Append log

				return 0
			}
		case "__log_err":
			return func(vm *vm.VirtualMachine) int64 {
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
				msg := vm.Memory[ptr : ptr+msgLen]

				(*workingVMTransaction).Logs = append((*workingVMTransaction).Logs, NewLog("error", msg, Error)) // Append log

				return 0
			}
		case "__log_return":
			return func(vm *vm.VirtualMachine) int64 {
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
				msg := vm.Memory[ptr : ptr+msgLen]

				(*workingVMTransaction).Logs = append((*workingVMTransaction).Logs, NewLog("return", msg, Return)) // Append log

				return 0
			}
		case "__transaction_get_nonce":
			return func(vm *vm.VirtualMachine) int64 {
				return int64(workingVMTransaction.AccountNonce)
			}
		case "__transaction_get_hash_nonce":
			return func(vm *vm.VirtualMachine) int64 {
				return int64(workingVMTransaction.HashNonce)
			}
		default:
			panic(fmt.Errorf("unknown field: %s", field)) // Panic
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module)) // Panic
	}
}

// ResolveGlobal defines a set of global variables for use within a WebAssembly module.
func (r *TransactionMetaResolver) ResolveGlobal(module, field string) int64 {
	switch module { // Handle module types
	case "env": // Env module
		switch field { // Handle fields
		case "__ursa_magic":
			return 424 // Return magic
		case "__transaction_nonce":
			return int64(workingVMTransaction.AccountNonce) // Return nonce
		case "__transaction_hash_nonce":
			return int64(workingVMTransaction.HashNonce) // Return nonce
		case "__transaction_amount":
			floatVal, _ := workingVMTransaction.Amount.Float64() // Get float val

			return int64(floatVal) // Return amount
		case "__transaction_timestamp":
			return int64(workingVMTransaction.Timestamp.Unix()) // Return timestamp
		default:
			panic(fmt.Errorf("unknown field: %s", field)) // Panic
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module)) // Panic
	}
}

/* END EXPORTED METHODS */
