// Package db defines the standard go-summercash transaction database.
package db

import (
	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
)

// Leaf defines a leaf in the go-summercash DAG.
type Leaf struct {
	Transaction *types.Transaction `json:"transaction"` // Transaction at leaf

	Parents  []*Leaf `json:"parents"`  // Transaction parents
	Children []*Leaf `json:"children"` // Leaf children

	Hash common.Hash `json:"hash"` // Transaction hash at leaf
}
