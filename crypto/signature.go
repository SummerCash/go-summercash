// Package crypto implements a set of abstractions enabling the creation,
// verification, and validation of ECDSA signatures. Additionally, the
// aforementioned package implements several helper types and methods for
// hashing via Blake2s.
package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"math/big"

	"github.com/stellar/go/xdr"
)

// Signature represents a standard ECDSA signature.
type Signature struct {
	PublicKey []byte `json:"public_key"` // Signature public key

	Data signatureData `json:"signature_data"` // Signature data
}

// signatureData represents the body data of a signature.
type signatureData struct {
	V []byte `json:"v"` // V
	R []byte `json:"r"` // R
	S []byte `json:"s"` // S
}

/* BEGIN EXPORTED METHODS */

// Sign signs the given message, message via ecdsa.
func Sign(privateKey *ecdsa.PrivateKey, message []byte) (*Signature, error) {
	hash := Blake2s(message) // Hash message

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash.Bytes()) // Calculate signature r, s values
	if err != nil {                                                // Check for errors
		return &Signature{}, err // Return found error
	}

	return &Signature{
		PublicKey: MarshalPublicKey(&privateKey.PublicKey), // Set marshalled public key
		Data: signatureData{
			V: hash.Bytes(), // Set message hash
			R: r.Bytes(),    // Set r
			S: s.Bytes(),    // Set s
		}, // Set signature data
	}, nil // Return initialized signature
}

// Verify verifies the contents of a given signature.
func (signature *Signature) Verify() bool {
	r, s := new(big.Int), new(big.Int) // Initialize big int buffer

	r.SetBytes(signature.Data.R) // Set bytes
	s.SetBytes(signature.Data.S) // Set bytes

	return ecdsa.Verify(UnmarshalPublicKey(signature.PublicKey), signature.Data.V, r, s) // Return is valid
}

// String converts the given signature to a JSON-formatted string.
func (signature *Signature) String() (string, error) {
	marshalled, err := json.Marshal(signature) // Marshal signature
	if err != nil {                            // Check for errors
		return "", err // Return found error
	}

	return string(marshalled), nil // Return JSON string val
}

// Bytes serializes the given signature to a byte slice via rlp encoding.
func (signature *Signature) Bytes() ([]byte, error) {
	var buffer bytes.Buffer // Initialize serialized data buffer

	_, err := xdr.Marshal(&buffer, *signature) // Marshal signature via XDR
	if err != nil {                            // Check for errors
		return nil, err // Return found error
	}

	return buffer.Bytes(), nil // Return marshalled data
}

// SignatureFromBytes deserializes a signature from a given byte slice, b.
func SignatureFromBytes(b []byte) (*Signature, error) {
	r := bytes.NewReader(b) // Initialize reader from provided data

	var signature Signature // Initialize signature buffer

	_, err := xdr.Unmarshal(r, &signature) // Unmarshal
	if err != nil {
		return &Signature{}, err // Return found error
	}

	return &signature, nil // Return unmarshalled signature
}

/* END EXPORTED METHODS */
