// Package common implements a set of commonly-used types and helper methods.
package common

import "crypto/elliptic"

// PreferredCurve defines a globally preferred ECDSA curve.
var PreferredCurve = elliptic.P521
