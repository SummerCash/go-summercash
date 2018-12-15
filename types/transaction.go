package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/crypto"
)

var (
	// ErrAlreadySigned - error definition stating transaction has already been signed
	ErrAlreadySigned = errors.New("transaction already signed")

	// ErrNilSignature - error definition describing nil tx signature
	ErrNilSignature = errors.New("nil signature")

	// ErrInvalidSignature - error definition describing invalid tx signature (doesn't match public key)
	ErrInvalidSignature = errors.New("invalid signature")
)

// Transaction - primitive transaction type
type Transaction struct {
	AccountNonce uint64 `json:"nonce"` // Nonce in set of account transactions

	Sender    *common.Address `json:"sender"`    // Transaction sender
	Recipient *common.Address `json:"recipient"` // Transaction recipient

	Amount float64 `json:"amount"` // Amount of coins sent in transaction

	Payload []byte `json:"payload"` // Misc. data transported with transaction

	Signature *Signature `json:"signature"` // Transaction signature meta

	ParentTx *Transaction `json:"-"` // Parent transaction

	Genesis bool `json:"genesis"` // Genesis

	Hash *common.Hash `json:"hash"` // Transaction hash
}

/* BEGIN EXPORTED METHODS */

// NewTransaction - attempt to initialize transaction primitive
func NewTransaction(nonce uint64, parentTx *Transaction, sender *common.Address, destination *common.Address, amount float64, payload []byte) (*Transaction, error) {
	transaction := Transaction{ // Init tx
		AccountNonce: nonce,       // Set nonce
		Sender:       sender,      // Set sender
		Recipient:    destination, // Set recipient
		Amount:       amount,      // Set amount
		Payload:      payload,     // Set tx payload
		ParentTx:     parentTx,    // Set parent
	}

	hash := common.NewHash(crypto.Sha3(transaction.Bytes())) // Hash transaction

	transaction.Hash = &hash // Set hash

	return &transaction, nil // Return initialized transaction
}

// TransactionFromBytes - serialize transaction from byte array
func TransactionFromBytes(b []byte) (*Transaction, error) {
	transaction := Transaction{} // Init buffer

	err := json.NewDecoder(bytes.NewReader(b)).Decode(&transaction) // Decode into buffer

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	if transaction.Signature != nil { // Check signature
		blockPub, _ := pem.Decode([]byte(transaction.Signature.SerializedPublicKey)) // Decode

		x509EncodedPub := blockPub.Bytes // Get x509 byte val

		genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub) // Parse public  key

		publicKey := genericPublicKey.(*ecdsa.PublicKey) // Get public key value

		(*(*transaction.Signature).PublicKey) = *publicKey // Set public key
	}

	return &transaction, nil // No error occurred, return read value
}

// Bytes - convert given transaction to byte array
func (tx *Transaction) Bytes() []byte {
	publicKey := ecdsa.PublicKey{} // Init buffer

	if tx.Signature != nil {
		publicKey = *(*(*tx).Signature).PublicKey // Set public key

		encoded, _ := x509.MarshalPKIXPublicKey(tx.Signature.PublicKey) // Encode

		pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded}) // Encode PEM

		(*(*tx).Signature).SerializedPublicKey = pemEncodedPub // Write encoded

		*(*tx).Signature.PublicKey = ecdsa.PublicKey{} // Set nil
	}

	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*tx) // Serialize tx

	if tx.Signature != nil {
		*(*(*tx).Signature).PublicKey = publicKey // Reset public key
	}

	return buffer.Bytes() // Return serialized
}

// String - convert given transaction to string
func (tx *Transaction) String() string {
	marshaled, _ := json.MarshalIndent(*tx, "", "  ") // Marshal tx

	return string(marshaled) // Return marshaled
}

/* END EXPORTED METHODS */
