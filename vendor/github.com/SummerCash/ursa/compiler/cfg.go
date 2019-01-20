package compiler

//CFGraph - cf graph struct
type CFGraph struct {
	Blocks []BasicBlock // Blocks
}

// BasicBlock - cf block struct
type BasicBlock struct {
	Code       []Instr   // Block instruction
	JmpKind    TyJmpKind // Block jmp (if any)
	JmpTargets []int     // Block jmp target (if any)

	JmpCond    TyValueID // Block jmp condition (if any)
	YieldValue TyValueID // Block yield value
}

// TyJmpKind - jmp type/identifier (uint val)
type TyJmpKind uint8

const (
	// JmpUndef - jmp instr undef
	JmpUndef TyJmpKind = iota

	// JmpUncond - jmp uncond
	JmpUncond

	// JmpEither - jmp either
	JmpEither

	// JmpTable - jmp table
	JmpTable

	// JmpReturn - jmp return
	JmpReturn
)

// ToInsSeq - convert given cfGraph to instruction set
func (g *CFGraph) ToInsSeq() []Instr {
	out := make([]Instr, 0) // Init instruction buffer

	blockRelocs := make([]int, len(g.Blocks)) // Init block re-allocs buffer
	blockEnds := make([]int, len(g.Blocks))   // Init block end buffer

	for i, bb := range g.Blocks { // Iterate through graph blocks
		blockRelocs[i] = len(out)    // Get burrent block reloc
		for _, op := range bb.Code { // Iterate through operations
			out = append(out, op) // Append found operation
		}

		out = append(out, Instr{}) // jmp placeholder
		blockEnds[i] = len(out)    // Set block ends
	}

	for i, bb := range g.Blocks { // Iterate through blocks
		jmpIns := &out[blockEnds[i]-1] // Get jmp ins

		jmpIns.Immediates = make([]int64, len(bb.JmpTargets)) // Init jmp immediates buffer

		for j, target := range bb.JmpTargets { // Iterate through block jmp targets
			jmpIns.Immediates[j] = int64(blockRelocs[target]) // Set jump immediate
		}

		switch bb.JmpKind { // Handle different jmp types
		case JmpUndef: // Check undef jmp
			panic("got JmpUndef")
		case JmpUncond: // Check uncond jmp
			jmpIns.Op = "jmp"
			jmpIns.Values = []TyValueID{bb.YieldValue}
		case JmpEither: // Check either jmp
			jmpIns.Op = "jmp_either"
			jmpIns.Values = []TyValueID{bb.JmpCond, bb.YieldValue}
		case JmpTable: // Check table jmp
			jmpIns.Op = "jmp_table"
			jmpIns.Values = []TyValueID{bb.JmpCond, bb.YieldValue}
		case JmpReturn: // Check return jmp
			jmpIns.Op = "return"

			if bb.YieldValue != 0 {
				jmpIns.Values = []TyValueID{bb.YieldValue}
			}
		default:
			panic("unreachable") // Panic
		}
	}

	return out // Return read instructions
}

// NewCFGraph - init new cf graph from function compiler
func (c *SSAFunctionCompiler) NewCFGraph() *CFGraph {
	g := &CFGraph{} // Init empty graph

	insLabels := make(map[int]int) // Init instruction label buffer

	insLabels[0] = 0 // Set placeholder
	nextLabel := 1   // Init iterator

	for i, ins := range c.Code { // Iterate through instructions
		switch ins.Op { // Check operation type
		case "jmp", "jmp_if", "jmp_either", "jmp_table":
			for _, target := range ins.Immediates { // Iterate through instruction immediates
				if _, ok := insLabels[int(target)]; !ok { // Check label is "ok"
					insLabels[int(target)] = nextLabel // Set label
					nextLabel++                        // Increment iterator
				}
			}

			if _, ok := insLabels[i+1]; !ok { // Check label is "ok"
				insLabels[i+1] = nextLabel // Set label
				nextLabel++                // Increment iterator
			}
		case "return":
			if _, ok := insLabels[i+1]; !ok { // Check label is "ok"
				insLabels[i+1] = nextLabel // Set label
				nextLabel++                // Increment iterator
			}
		}
	}

	g.Blocks = make([]BasicBlock, nextLabel) // Init block buffer
	var currentBlock *BasicBlock             // Init current block buffer

	for i, ins := range c.Code { // Iterate through instructions
		if label, ok := insLabels[i]; ok { // Check label is "ok"
			if currentBlock != nil { // Check if block is not nil
				currentBlock.JmpKind = JmpUncond       // Set jmp type
				currentBlock.JmpTargets = []int{label} // Set jmp target
			}

			currentBlock = &g.Blocks[label] // Set current block
		}

		switch ins.Op { // Handle different instruction types
		case "jmp": // Check is jump type
			currentBlock.JmpKind = JmpUncond                                   // Set jmp kind
			currentBlock.JmpTargets = []int{insLabels[int(ins.Immediates[0])]} // Set jmp target
			currentBlock.YieldValue = ins.Values[0]                            // Set yield
			currentBlock = nil                                                 // Set current block to nil
		case "jmp_if": // Check is jmp_if type
			currentBlock.JmpKind = JmpEither                                                        // Set jmp kind
			currentBlock.JmpTargets = []int{insLabels[int(ins.Immediates[0])], insLabels[int(i+1)]} // Set jmp target
			currentBlock.JmpCond = ins.Values[0]                                                    // Set jmp condition
			currentBlock.YieldValue = ins.Values[1]                                                 // Set yield
			currentBlock = nil                                                                      // Set current block to nil
		case "jmp_either": // Check is jmp_either type
			currentBlock.JmpKind = JmpEither                                                                      // Set jmp kind
			currentBlock.JmpTargets = []int{insLabels[int(ins.Immediates[0])], insLabels[int(ins.Immediates[1])]} // Set jmp target
			currentBlock.JmpCond = ins.Values[0]                                                                  // Set jmp condition
			currentBlock.YieldValue = ins.Values[1]                                                               // Set yield
			currentBlock = nil                                                                                    // Set current block to nil
		case "jmp_table": // Check is jmp_table type
			currentBlock.JmpKind = JmpTable                            // Set jmp kind
			currentBlock.JmpTargets = make([]int, len(ins.Immediates)) // Set jmp target

			for j, imm := range ins.Immediates { // Iterate through jmp immediates
				currentBlock.JmpTargets[j] = insLabels[int(imm)] // Set jmp target
			}

			currentBlock.JmpCond = ins.Values[0]    // Set jmp condition
			currentBlock.YieldValue = ins.Values[1] // Set yield value
			currentBlock = nil                      // Set current block to nil
		case "return": // Check is return
			currentBlock.JmpKind = JmpReturn // Set jmp kind

			if len(ins.Values) > 0 { // Check returns more than nil
				currentBlock.YieldValue = ins.Values[0] // Set yield
			}

			currentBlock = nil // Set current block to nil
		default:
			currentBlock.Code = append(currentBlock.Code, ins) // Append instruction
		}
	}

	if label, ok := insLabels[len(c.Code)]; ok { // Check label is "ok"
		lastBlock := &g.Blocks[label] // Init last block

		if lastBlock.JmpKind != JmpUndef { // Check jmp kind is undefined
			panic("last block should always have an undefined jump target") // Panic
		}

		lastBlock.JmpKind = JmpReturn // Set jmp kind
	}

	return g // Return graph
}
