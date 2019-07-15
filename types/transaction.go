package types

// Transaction defines a generic transaction.
type Transaction struct {
	Nonce uint64 `json:"nonce"` // Index of transaction in account tx set

	GasLimit uint64 `json:"gas_limit"` // Number of finks willing to spend on gas
}
