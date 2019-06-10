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
	"math/big"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/crypto"
	"github.com/SummerCash/ursa/compiler"
	"github.com/SummerCash/ursa/vm"
)

var (
	// ErrAlreadySigned - error definition stating transaction has already been signed
	ErrAlreadySigned = errors.New("transaction already signed")

	// ErrNilSignature - error definition describing nil tx signature
	ErrNilSignature = errors.New("nil signature")

	// ErrInvalidSignature - error definition describing invalid tx signature (doesn't match public key)
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrIsNotContractCall - error definition describing tx of non-contract-call type
	ErrIsNotContractCall = errors.New("transaction is not contract call")
)

// Transaction - primitive transaction type
type Transaction struct {
	AccountNonce uint64 `json:"nonce"` // Nonce in set of account transactions

	HashNonce uint64 `json:"hash_nonce"` // Nonce to calculate valid hash

	Sender    *common.Address `json:"sender"`    // Transaction sender
	Recipient *common.Address `json:"recipient"` // Transaction recipient

	Amount *big.Float `json:"amount"` // Amount of coins sent in transaction

	Payload []byte `json:"payload"` // Misc. data transported with transaction

	Signature *Signature `json:"signature"` // Transaction signature meta

	ParentTx *common.Hash `json:"parent_hash"` // Parent transaction

	Timestamp time.Time `json:"time"` // Transaction timestamp

	DeployedContractAddress *common.Address `json:"contract"` // Contract instance

	ContractCreation bool `json:"is-init-contract"` // Should init contract
	Genesis          bool `json:"genesis"`          // Genesis

	State *vm.State `json:"state"` // State

	Logs []*Log `json:"logs"` // Logs

	Hash *common.Hash `json:"hash"` // Transaction hash
}

// StringTransaction represents a human-readable transaction.
type StringTransaction struct {
	AccountNonce uint64 `json:"nonce"` // Nonce in set of account transactions

	HashNonce uint64 `json:"hash_nonce"` // Nonce to calculate valid hash

	SenderHex    string `json:"sender"`    // Transaction sender
	RecipientHex string `json:"recipient"` // Transaction recipient

	Amount float64 `json:"amount"` // Amount of coins sent in transaction

	Payload []byte `json:"payload"` // Misc. data transported with transaction

	Signature *Signature `json:"signature"` // Transaction signature meta

	ParentTx string `json:"parent_hash"` // Parent transaction

	Timestamp string `json:"time"` // Transaction timestamp

	DeployedContractAddress *common.Address `json:"contract"` // Contract instance

	ContractCreation bool `json:"is-init-contract"` // Should init contract
	Genesis          bool `json:"genesis"`          // Genesis

	State *vm.State `json:"state"` // State

	Logs []*Log `json:"logs"` // Logs

	HashHex string `json:"hash"` // Transaction hash
}

/* BEGIN EXPORTED METHODS */

// NewTransaction - attempt to initialize transaction primitive
func NewTransaction(nonce uint64, parentTx *Transaction, sender *common.Address, destination *common.Address, amount *big.Float, payload []byte) (*Transaction, error) {
	parentHash := &common.Hash{} // Init hash buffer

	if parentTx != nil { // Check has parent
		parentHash = parentTx.Hash // Set parent hash
	}

	transaction := Transaction{ // Init tx
		AccountNonce:     nonce,            // Set nonce
		HashNonce:        0,                // Set hash nonce
		Sender:           sender,           // Set sender
		Recipient:        destination,      // Set recipient
		Amount:           amount,           // Set amount
		Payload:          payload,          // Set tx payload
		ParentTx:         parentHash,       // Set parent
		Timestamp:        time.Now().UTC(), // Set timestamp
		ContractCreation: false,            // Set should init contract
	}

	hash := common.NewHash(crypto.Sha3(transaction.Bytes())) // Hash transaction

	for bytes.Contains(hash.Bytes(), []byte{'\r'}) { // Do until does not contain escape character
		transaction.HashNonce++ // Increment hash nonce

		hash = common.NewHash(crypto.Sha3(transaction.Bytes())) // Set hash
	}

	transaction.Hash = &hash // Set hash

	return &transaction, nil // Return initialized transaction
}

// NewContractCreation - initialize contract designated to an initialized contract, calling contract constructor/provided constructor
func NewContractCreation(nonce uint64, parentTx *Transaction, sender *common.Address, contractInstance *common.Address, amount *big.Float, payload []byte) (*Transaction, error) {
	transaction := Transaction{ // Init tx
		AccountNonce:            nonce,            // Set nonce
		Sender:                  sender,           // Set sender
		Recipient:               contractInstance, // Set dest
		Amount:                  amount,           // Set amount
		Payload:                 payload,          // Set tx payload
		ParentTx:                parentTx.Hash,    // Set parent
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
	if err != nil {                                                 // Check for errors
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

// EvaluateNewState evaluates the new state for a given transaction.
// Does not set the transactions' state, but does return pointer to new state.
func (transaction *Transaction) EvaluateNewState(gasPolicy *compiler.GasPolicy) (*vm.State, error) {
	if transaction.Payload == nil && !bytes.Contains(transaction.Payload, []byte("(")) { // Check is not contract call
		return &vm.State{}, ErrIsNotContractCall // Return error
	}

	recipientChain, err := ReadChainFromMemory(*transaction.Recipient) // Read recipient chain
	if err != nil {                                                    // Check for errors
		return &vm.State{}, err // Return found error
	}

	workingVM, err := vm.NewVirtualMachine(recipientChain.ContractSource, common.VMConfig, new(TransactionMetaResolver), common.GasPolicy) // Init vm
	if err != nil {                                                                                                                        // Check for errors
		return nil, err // Return found error
	}

	workingVMTransaction = transaction // Set working vm tx

	parentTx, err := recipientChain.QueryTransaction(*transaction.ParentTx) // Query parent
	if err != nil {                                                         // Check for errors
		return &vm.State{}, err // Return found error
	}

	if transaction.ParentTx != nil { // Check has parent
		workingVM.CallStack = parentTx.State.CallStack               // Set call stack
		workingVM.CurrentFrame = parentTx.State.CurrentFrame         // Set current frame
		workingVM.Table = parentTx.State.Table                       // Set table
		workingVM.Globals = parentTx.State.Globals                   // Set globals
		workingVM.Memory = parentTx.State.Memory                     // Set memory
		workingVM.NumValueSlots = parentTx.State.NumValueSlots       // Set num value slots
		workingVM.Yielded = parentTx.State.Yielded                   // Set yielded
		workingVM.InsideExecute = parentTx.State.InsideExecute       // Set inside execute
		workingVM.Exited = parentTx.State.Exited                     // Set has exited
		workingVM.ExitError = parentTx.State.ExitError               // Set exit error
		workingVM.ReturnValue = parentTx.State.ReturnValue           // Set return value
		workingVM.Gas = parentTx.State.Gas                           // Set gas
		workingVM.GasLimitExceeded = parentTx.State.GasLimitExceeded // Set gas limit exceeded
	}

	callMethod, callParams, err := common.ParseStringMethodCallNoReceiver(string(transaction.Payload)) // Parse payload method call
	if err != nil {                                                                                    // Check for errors
		return &vm.State{}, err // Return found error
	}

	entryID, valid := workingVM.GetFunctionExport(callMethod) // Get function ID from payload

	if !valid { // Check for errors
		return &vm.State{}, ErrInvalidPayload // Return found error
	}

	var parsedCallParams []int64 // Init params buffer

	for _, param := range callParams { // Iterate through params
		intVal, err := strconv.ParseInt(param, 10, 64) // Parse int
		if err != nil {                                // Check for errors
			return nil, err // Return found error
		}

		parsedCallParams = append(parsedCallParams, intVal) // Append parse param
	}

	result, err := workingVM.Run(entryID, parsedCallParams...) // Run
	if err != nil {                                            // Check for errors
		common.Logf("== VM == Contract call exited with code %d and error %s", result, err.Error()) // Log err

		return &vm.State{}, err // Return found error
	}

	common.Logf("== VM == Contract call exited with code %d", result) // Log finish

	return &vm.State{
		CallStack:        workingVM.CallStack,        // Set call stack
		CurrentFrame:     workingVM.CurrentFrame,     // Set current frame
		Table:            workingVM.Table,            // Set table
		Globals:          workingVM.Globals,          // Set globals
		Memory:           workingVM.Memory,           // Set memory
		NumValueSlots:    workingVM.NumValueSlots,    // Set num value slots
		Yielded:          workingVM.Yielded,          // Set yielded
		InsideExecute:    workingVM.InsideExecute,    // Set inside execute
		Exited:           workingVM.Exited,           // Set has exited
		ExitError:        workingVM.ExitError,        // Set exit error
		ReturnValue:      workingVM.ReturnValue,      // Set return value
		Gas:              workingVM.Gas,              // Set gas
		GasLimitExceeded: workingVM.GasLimitExceeded, // Set gas limit exceeded
	}, nil // Return state
}

// Publish - publish given transaction
func (transaction *Transaction) Publish() error {
	if transaction.Signature == nil { // Check nil pointer
		return ErrNilSignature // Return error
	}

	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain
	if err != nil {                                             // Check for errors
		return err // Return found error
	}

	node, err := coordinationChain.QueryAddress(*transaction.Recipient) // Get address
	if err != nil {                                                     // Check for errors
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
		if err != nil {                                                            // Check for errors
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
		if err != nil {                                                  // Check for errors
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
	floatVal, _ := transaction.Amount.Float64() // Get float value

	var senderHex, recipientHex string // Init hex buffer

	parent := "" // Init parent buffer

	if transaction.Sender != nil { // Check has sender
		senderHex = transaction.Sender.String() // Set string
	}

	if transaction.Recipient != nil { // Check has recipient
		recipientHex = transaction.Recipient.String() // Set string
	}

	if transaction.ParentTx != nil { // Check has parent
		parent = transaction.ParentTx.String() // Set parent
	}

	stringTransaction := &StringTransaction{
		AccountNonce:            transaction.AccountNonce,                           // Set account nonce
		SenderHex:               senderHex,                                          // Set sender hex
		RecipientHex:            recipientHex,                                       // Set recipient hex
		Amount:                  floatVal,                                           // Set amount
		Payload:                 transaction.Payload,                                // Set payload
		Signature:               transaction.Signature,                              // Set signature
		ParentTx:                parent,                                             // Set parent
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
	err := common.CreateDirIfDoesNotExist(fmt.Sprintf("%s/mem/pending_tx", common.DataDir)) // Create dir if necessary
	if err != nil {                                                                         // Check for errors
		return err // Return error
	}

	err = transaction.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		return err // Return error
	}

	json, err := json.MarshalIndent(*transaction, "", "  ") // Marshal JSOn
	if err != nil {                                         // Check for errors
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
	if err != nil {                                                                                                             // Check for errors
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
