package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/crypto"
)

// Signature - struct containing signature values
type Signature struct {
	PublicKey           *ecdsa.PublicKey `json:"-"` // Public key
	SerializedPublicKey []byte           // Serialized public key

	V []byte   // Hash signature value
	R *big.Int // Signature R
	S *big.Int // Signature S
}

// Witness - struct containing transaction witness values
type Witness struct {
	Signature *Signature // Witness signature

	Address common.Address // Witnessing address

	Weight float64 // Witness weight (stake)
}

var (
	// ErrInvalidPublicKey - error definition describing a public key input not equal to transaction sender address
	ErrInvalidPublicKey = errors.New("signing public key does not match transaction public key")

	// ErrCannotWitnessSelf - error definition describing a transaction that has already been signed by an attempted witness
	ErrCannotWitnessSelf = errors.New("cannot witness self-signed transaction")
)

/* BEGIN EXPORTED METHODS */

// SignTransaction - sign given transaction
func SignTransaction(transaction *Transaction, privateKey *ecdsa.PrivateKey) error {
	return selfSignTransaction(transaction, privateKey) // Sign tx
}

// VerifyTransactionSignature - verify given transaction signature, returning false if signature invalid
func VerifyTransactionSignature(transaction *Transaction) (bool, error) {
	if transaction.Signature == nil { // Check nil signature
		return false, ErrNilSignature // Return nil signature error
	} else if common.PublicKeyToAddress(transaction.Signature.PublicKey) != *transaction.Sender { // Check for invalid public key
		return false, ErrInvalidSignature // Return invalid signature error
	}

	return ecdsa.Verify(transaction.Signature.PublicKey, transaction.Signature.V, transaction.Signature.R, transaction.Signature.S), nil // Check signature valid
}

// Bytes - convert given signature to byte array
func (signature *Signature) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*signature) // Serialize signature

	return buffer.Bytes() // Return serialized
}

// String - convert given signature to string
func (signature *Signature) String() string {
	marshaled, _ := json.MarshalIndent(*signature, "", "  ") // Marshal signature

	return string(marshaled) // Return marshaled
}

// Bytes - convert given witness to byte array
func (witness *Witness) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*witness) // Serialize witness

	return buffer.Bytes() // Return serialized
}

// String - convert given witness to string
func (witness *Witness) String() string {
	marshaled, _ := json.MarshalIndent(*witness, "", "  ") // Marshal witness

	return string(marshaled) // Return marshaled
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// selfSignTransaction - handle signing of transaction by sender
func selfSignTransaction(transaction *Transaction, privateKey *ecdsa.PrivateKey) error {
	if transaction.Signature != nil { // Check not already signed
		return ErrAlreadySigned // Return already signed error
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, crypto.Sha3(transaction.Bytes())) // Sign tx

	if err != nil { // Check for errors
		return err // Return found error
	}

	txSignature := Signature{ // Initialize signature
		PublicKey: &privateKey.PublicKey,            // Set public key
		V:         crypto.Sha3(transaction.Bytes()), // Set val
		R:         r,                                // Set R
		S:         s,                                // Set S
	}

	(*transaction).Signature = &txSignature // Set signature

	return nil // No error occurred, return nil
}

/* END INTERNAL METHODS */
