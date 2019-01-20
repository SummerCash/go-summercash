package compiler

// InsertGasCounters - insert gas counter to given wasm
func (c *SSAFunctionCompiler) InsertGasCounters(gp GasPolicy) {
	cfg := c.NewCFGraph() // Init cf graph

	for _, block := range cfg.Blocks { // Iterate through blocks
		blk := &block // Get current block pointer

		totalCost := int64(0) // Init gas buffer

		for _, ins := range blk.Code { // Iterate through instructions
			totalCost += gp.GetCost(ins.Op) // Get cost of instruction

			if totalCost < 0 { // Check cost will cause overflow
				panic("total cost overflow") // Panic with err
			}
		}

		if totalCost != 0 { // Check total cost is not nil
			blk.Code = append([]Instr{ // Append add_gas instruction
				buildInstr(0, "add_gas", []int64{totalCost}, []TyValueID{}),
			}, blk.Code...)
		}
	}

	c.Code = cfg.ToInsSeq() // Set code with added gas instruction
}
