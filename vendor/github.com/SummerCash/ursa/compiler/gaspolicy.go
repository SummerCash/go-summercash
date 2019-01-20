package compiler

// GasPolicy - gas policy
type GasPolicy interface {
	GetCost(key string) int64 // Get gas cost
}

// SimpleGasPolicy - simple policy setting simple gas per instruction
type SimpleGasPolicy struct {
	GasPerInstruction int64 // Gas per instruction
}

// GetCost - get gas cost
func (p *SimpleGasPolicy) GetCost(key string) int64 {
	return p.GasPerInstruction
}
