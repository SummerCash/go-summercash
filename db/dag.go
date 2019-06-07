// Package db defines the standard go-summercash transaction database.
package db

// Dag implements the standard directed acyclic
// graph global chain.
type Dag struct {
	Root *Leaf `json:"root"` // Root leaf
}
