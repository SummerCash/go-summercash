package common

import (
	"encoding/hex"
	"errors"

	"github.com/SummerCash/go-summercash/crypto"
)

var (
	// MemPrefix - "0x" memory address  prefixbyte array representation
	MemPrefix = []byte{48, 120}

	// ErrNoMem - error definition describing input < 2 chars long (< len(0x))
	ErrNoMem = errors.New("insufficient length for memory address char")
)

/* BEGIN EXPORTED METHODS */

// Encode - encode given byte array to hex-formatted, MemPrefix-compliant byte array
func Encode(b []byte) ([]byte, error) {
	if len(b) == 0 { // Check for nil input
		return []byte{}, ErrNilInput // Return error
	} else if len(b) < 2 { // Check no 0x
		return []byte{}, ErrNoMem // Return error
	}

	if string(crypto.Sha3(b[0:2])) == string(crypto.Sha3(MemPrefix)) { // Check already encoded
		b = b[2:] // Trim 0x
	}

	enc := make([]byte, len(b)*2+2) // Init encoder buffer

	copy(enc, "0x") // Copy prefix

	hex.Encode(enc[2:], b) // Encode given byte array

	return enc, nil // Return encoded
}

// EncodeString - encode given byte array to hex-formatted, MemPrefix-compliant string
func EncodeString(b []byte) (string, error) {
	if len(b) == 0 { // Check for nil input
		return "", ErrNilInput // Return error
	} else if len(b) < 2 { // Check no 0x
		return "", ErrNoMem // Return error
	}

	enc := make([]byte, len(b)*2+2) // Init encoder buffer

	copy(enc, "0x") // Copy prefix

	hex.Encode(enc[2:], b) // Encode given byte array

	return string(enc), nil // Return encoded string
}

// Decode - decode given hex-formatted, MemPrefix-compliant byte array to standard byte array
func Decode(b []byte) ([]byte, error) {
	if len(b) == 0 { // Check for nil input
		return []byte{}, ErrNilInput // Return error
	} else if len(b) < 2 { // Check no 0x
		return []byte{}, ErrNoMem // Return error
	}

	b = b[2:] // Trim 0x

	n, err := hex.Decode(b, b) // Decode
	if err != nil {            // Check for errors
		return []byte{}, err // Return found error
	}

	return b[:n], nil // Return decoded
}

// DecodeString - decode given hex-formatted, MemPrefix-compliant string to standard byte array
func DecodeString(s string) ([]byte, error) {
	if len(s) == 0 { // Check for nil input
		return []byte{}, ErrNilInput // Return error
	} else if len(s) < 2 { // Check no 0x
		return []byte{}, ErrNoMem // Return error
	}

	b, err := hex.DecodeString(s[2:]) // Decode input (trimming 0x)
	if err != nil {                   // Check for errors
		return []byte{}, err // Return found error
	}

	return b, err // Return decoded
}

/* END EXPORTED METHODS */
