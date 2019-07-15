// Package crypto implements a set of abstractions enabling the creation,
// verification, and validation of ECDSA signatures. Additionally, the
// aforementioned package implements several helper types and methods for
// hashing via Blake2s.
package crypto

import "crypto/elliptic"

// PreferredCurve defines a globally preferred ECDSA curve.
var PreferredCurve = elliptic.P521
