package compiler

import "testing"

// TestGetCost - test functionality of gas cost calculator
func TestGetCost(t *testing.T) {
	simpleGasPolicy := SimpleGasPolicy{GasPerInstruction: 1} // Init simple gas policy

	gasCost := simpleGasPolicy.GetCost("") // Get cost

	t.Log(gasCost) // Log success
}
