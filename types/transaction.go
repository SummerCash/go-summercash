package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/crypto"
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

	Timestamp time.Time `json:"time"` // Transaction timestamp

	DeployedContractAddress *common.Address `json:"contract"` // Contract instance

	ContractCreation bool `json:"is-init-contract"` // Should init contract
	Genesis          bool `json:"genesis"`          // Genesis

	Logs []*Log `json:"logs"` // Logs

	Hash *common.Hash `json:"hash"` // Transaction hash
}

// StringTransaction represents a human-readable transaction.
type StringTransaction struct {
	AccountNonce uint64 `json:"nonce"` // Nonce in set of account transactions

	SenderHex    string `json:"sender"`    // Transaction sender
	RecipientHex string `json:"recipient"` // Transaction recipient

	Amount float64 `json:"amount"` // Amount of coins sent in transaction

	Payload []byte `json:"payload"` // Misc. data transported with transaction

	Signature *Signature `json:"signature"` // Transaction signature meta

	ParentTx *Transaction `json:"-"` // Parent transaction

	Timestamp string `json:"time"` // Transaction timestamp

	DeployedContractAddress *common.Address `json:"contract"` // Contract instance

	ContractCreation bool `json:"is-init-contract"` // Should init contract
	Genesis          bool `json:"genesis"`          // Genesis

	Logs []*Log `json:"logs"` // Logs

	HashHex string `json:"hash"` // Transaction hash
}

/* BEGIN EXPORTED METHODS */

// NewTransaction - attempt to initialize transaction primitive
func NewTransaction(nonce uint64, parentTx *Transaction, sender *common.Address, destination *common.Address, amount float64, payload []byte) (*Transaction, error) {
	transaction := Transaction{ // Init tx
		AccountNonce:     nonce,            // Set nonce
		Sender:           sender,           // Set sender
		Recipient:        destination,      // Set recipient
		Amount:           amount,           // Set amount
		Payload:          payload,          // Set tx payload
		ParentTx:         parentTx,         // Set parent
		Timestamp:        time.Now().UTC(), // Set timestamp
		ContractCreation: false,            // Set should init contract
	}

	hash := common.NewHash(crypto.Sha3(transaction.Bytes())) // Hash transaction

	transaction.Hash = &hash // Set hash

	return &transaction, nil // Return initialized transaction
}

// NewContractCreation - initialize contract designated to an initialized contract, calling contract constructor/provided constructor
func NewContractCreation(nonce uint64, parentTx *Transaction, sender *common.Address, contractInstance *common.Address, amount float64, payload []byte) (*Transaction, error) {
	transaction := Transaction{ // Init tx
		AccountNonce:            nonce,            // Set nonce
		Sender:                  sender,           // Set sender
		Recipient:               contractInstance, // Set dest
		Amount:                  amount,           // Set amount
		Payload:                 payload,          // Set tx payload
		ParentTx:                parentTx,         // Set parent
		Timestamp:               time.Now().UTC(), // Set timestamp
		ContractCreation:        true,             // Set should init contract
		DeployedContractAddress: contractInstance, // Set deployed
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

		((*transaction.Signature).PublicKey) = publicKey // Set public key
	}

	return &transaction, nil // No error occurred, return read value
}

// Publish - publish given transaction
func (transaction *Transaction) Publish() error {
	if transaction.Signature == nil { // Check nil pointer
		return ErrNilSignature // Return error
	}

	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain

	if err != nil { // Check for errors
		return err // Return found error
	}

	node, err := coordinationChain.QueryAddress(*transaction.Recipient) // Get address

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = common.SendBytes(transaction.Bytes(), node.Addresses[0]) // Send transaction

	if err != nil { // Check for errors
		common.Logf("== ERROR == error pushing transaction %s to peer %s %s\n", transaction.Hash.String(), node.Addresses[0], err.Error()) // Log error pushing
	}

	common.Logf("== NETWORK == pushing transaction %s to node %s\n", transaction.Hash.String(), node.Addresses[0]) // Log push

	for x, address := range node.Addresses { // Iterate through addresses
		common.Logf("== NETWORK == pushing transaction %s to node %s\n", transaction.Hash.String(), address) // Log push

		if x != 0 { // Skip first index
			go common.SendBytes(transaction.Bytes(), address) // Send transaction
		}
	}

	return nil // No error occurred, return nil
}

// MakeEncodingSafe - encode transaction to safe format
func (transaction *Transaction) MakeEncodingSafe() error {
	if transaction.Signature != nil && transaction.Signature.PublicKey != nil { // Check has signature
		encoded, err := x509.MarshalPKIXPublicKey(transaction.Signature.PublicKey) // Encode

		if err != nil { // Check for errors
			return err // Return error
		}

		pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded}) // Encode PEM

		(*(*transaction).Signature).SerializedPublicKey = pemEncodedPub // Write encoded

		*(*transaction).Signature.PublicKey = ecdsa.PublicKey{} // Set nil
	}

	return nil // No error occurred, return nil
}

// RecoverSafeEncoding - recover transaction from safe encoding
func (transaction *Transaction) RecoverSafeEncoding() error {
	if transaction.Signature != nil { // Check has signature
		blockPub, _ := pem.Decode([]byte(transaction.Signature.SerializedPublicKey)) // Decode

		x509EncodedPub := blockPub.Bytes // Get x509 byte val

		genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub) // Parse public  key

		if err != nil { // Check for errors
			return err // Return found error
		}

		publicKey := genericPublicKey.(*ecdsa.PublicKey) // Get public key value

		((*transaction.Signature).PublicKey) = publicKey // Set public key
	}

	return nil // No error occurred, return nil
}

// Bytes - convert given transaction to byte array
func (transaction *Transaction) Bytes() []byte {
	publicKey := ecdsa.PublicKey{} // Init buffer

	if transaction.Signature != nil {
		publicKey = *(*(*transaction).Signature).PublicKey // Set public key

		encoded, _ := x509.MarshalPKIXPublicKey(transaction.Signature.PublicKey) // Encode

		pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded}) // Encode PEM

		(*(*transaction).Signature).SerializedPublicKey = pemEncodedPub // Write encoded

		*(*transaction).Signature.PublicKey = ecdsa.PublicKey{} // Set nil
	}

	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*transaction) // Serialize tx

	if transaction.Signature != nil {
		*(*(*transaction).Signature).PublicKey = publicKey // Reset public key
	}

	return buffer.Bytes() // Return serialized
}

// String - convert given transaction to string
func (transaction *Transaction) String() string {
	stringTransaction := &StringTransaction{
		AccountNonce:            transaction.AccountNonce,                           // Set account nonce
		SenderHex:               transaction.Sender.String(),                        // Set sender hex
		RecipientHex:            transaction.Recipient.String(),                     // Set recipient hex
		Amount:                  transaction.Amount,                                 // Set amount
		Payload:                 transaction.Payload,                                // Set payload
		Signature:               transaction.Signature,                              // Set signature
		ParentTx:                transaction.ParentTx,                               // Set parent
		Timestamp:               transaction.Timestamp.Format("01/02/2006 3:04 PM"), // Set timestamp
		DeployedContractAddress: transaction.DeployedContractAddress,                // Set deployed contract address
		ContractCreation:        transaction.ContractCreation,                       // Set is contract creation
		Genesis:                 transaction.Genesis,                                // Set is genesis
		Logs:                    transaction.Logs,                                   // Set logs
		HashHex:                 transaction.Hash.String(),                          // Set hash hex
	}

	marshaled, _ := json.MarshalIndent(*stringTransaction, "", "  ") // Marshal tx

	return string(marshaled) // Return marshaled
}

// WriteToMemory - write given transaction to memory
func (transaction *Transaction) WriteToMemory() error {
	err := common.CreateDirIfDoesNotExit(fmt.Sprintf("%s/mem/pending_tx", common.DataDir)) // Create dir if necessary

	if err != nil { // Check for errors
		return err // Return error
	}

	err = transaction.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		return err // Return error
	}

	json, err := json.MarshalIndent(*transaction, "", "  ") // Marshal JSOn

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/mem/pending_tx/tx_%s.gob", common.DataDir, transaction.Hash.String())), json, 0644) // Write chainConfig to JSON

	if err != nil { // Check for errors
		return err // Return error
	}

	err = transaction.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// ReadTransactionFromMemory - read transaction from memory
func ReadTransactionFromMemory(hash common.Hash) (*Transaction, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/mem/pending_tx/tx_%s.gob", common.DataDir, hash.String()))) // Read file

	if err != nil { // Check for errors
		return &Transaction{}, err // Return error
	}

	buffer := &Transaction{} // Initialize buffer

	err = json.Unmarshal(data, buffer) // Read json into buffer

	if err != nil { // Check for errors
		return &Transaction{}, err // Return error
	}

	err = buffer.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return &Transaction{}, err // Return error
	}

	return buffer, nil // No error occurred, return read tx
}

/* END EXPORTED METHODS */
