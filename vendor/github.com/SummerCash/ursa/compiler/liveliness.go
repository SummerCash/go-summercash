package compiler

// FIXME: The current RegAlloc is based on wasm stack info and we probably
// want a real one (in addition to this) with liveness analysis.

// RegAlloc - returns the total number of registers used
func (c *SSAFunctionCompiler) RegAlloc() int {
	regID := TyValueID(1) // Init reg ID

	valueRelocs := make(map[TyValueID]TyValueID) // Init reloc buffer

	for _, values := range c.StackValueSets { // Iterate through stack values
		for _, v := range values { // Iterate through values
			valueRelocs[v] = regID // Set reloc value
		}

		regID++ // Increment iterator
	}

	for i := range c.Code { // Iterate through instructions
		ins := &c.Code[i] // Get current instruction

		if ins.Target != 0 { // Check has instruction target
			if reg, ok := valueRelocs[ins.Target]; ok { // Check label is "ok"
				ins.Target = reg // Set target
			} else {
				panic("Register not found for target") // Panic
			}
		}

		for j, v := range ins.Values { // Iterate through instruction values
			if v != 0 { // Check has value
				if reg, ok := valueRelocs[v]; ok { // Check is "ok"
					ins.Values[j] = reg // Set value
				} else {
					panic("Register not found for value") // Panic
				}
			}
		}
	}

	return int(regID) // Return read regalloc
}
