package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/crypto"
	"github.com/SummerCash/ursa/vm"
)

// Chain - account transactions chain
type Chain struct {
	Account common.Address `json:"account"` // Chain account

	Transactions []*Transaction `json:"transactions"` // Transactions in chain

	Genesis common.Hash `json:"genesis"` // Genesis block hash

	ContractSource []byte `json:"contract"` // Contract

	NetworkID uint        `json:"network"` // Network ID (mainnet: 0, testnet: 1, etc...)
	ID        common.Hash `json:"ID"`      // Chain ID
}

var (
	// ErrChainAlreadyExists - error definition describing a given chain that has already been registered in the coordinationChain
	ErrChainAlreadyExists = errors.New("chain already exists for given account")

	// ErrGenesisAlreadyExists - error definition describing a given chain with an existing genesis block
	ErrGenesisAlreadyExists = errors.New("chain already has existing genesis")

	// ErrIrrelevantTransaction - error definition describing a transaction outside the scope of the given chain
	ErrIrrelevantTransaction = errors.New("irrelevant transaction")

	// ErrInsufficientBalance - error definition describing a transaction worth more than the sender's balance
	ErrInsufficientBalance = errors.New("insufficient transaction sender balance")

	// ErrNilTransaction - error definition describing a query for a non-existent transaction
	ErrNilTransaction = errors.New("couldn't find transaction with given hash")

	// ErrDuplicateTransaction - error definition describing a transaction that already exists in a given chain
	ErrDuplicateTransaction = errors.New("duplicate transaction")

	// ErrInvalidPayload - error definition describing an invalid transaction payload
	ErrInvalidPayload = errors.New("invalid payload")
)

/* BEGIN EXPORTED METHODS */

// NewChain - initialize new chain
func NewChain(account common.Address) (*Chain, error) {
	config, err := config.ReadChainConfigFromMemory() // Read config from memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	_, err = ReadChainFromMemory(account) // Check chain doesn't already exist

	if err == nil { // Check already exists
		return &Chain{}, ErrChainAlreadyExists // Return error
	}

	chain := &Chain{ // Init chain
		Account:      account,
		Transactions: []*Transaction{},
		NetworkID:    config.NetworkID,
	}

	(*chain).ID = common.NewHash(crypto.Sha3(chain.Bytes())) // Set ID

	common.Logf("== ACCOUNT == initialized account chain with account address %s\n", chain.Account.String()) // Log init

	err = chain.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	return chain, nil // Return initialized chain
}

// NewContractChain - initialize new contract chain
func NewContractChain(account common.Address, contractSource []byte, deploymentTransaction *Transaction) (*Chain, error) {
	config, err := config.ReadChainConfigFromMemory() // Read config from memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	_, err = ReadChainFromMemory(account) // Read chain

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	chain := &Chain{ // Init chain
		Account:        account,
		Transactions:   []*Transaction{deploymentTransaction},
		NetworkID:      config.NetworkID,
		ContractSource: contractSource,
	}

	(*chain).ID = common.NewHash(crypto.Sha3(chain.Bytes())) // Set ID

	common.Logf("== ACCOUNT == initialized account chain with account address %s\n", chain.Account.String()) // Log init

	err = chain.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	return chain, nil // Return initialized chain
}

// AddTransaction - append given transaction to chain
func (chain *Chain) AddTransaction(transaction *Transaction) error {
	if chain.ContractSource != nil && transaction.Payload != nil { // Check is contract call
		chain.handleContractCall(transaction) // Handle contract call
	}

	chainConfig, err := config.ReadChainConfigFromMemory() // Read chain config

	if err != nil { // Check for errors
		return err // Return found error
	}

	genesisChain, err := ReadGenesisChainFromMemory(chainConfig) // Read genesis with config

	if err != nil { // Check for errors
		return err // Return found error
	}

	if transaction.Signature == nil && err == nil { // Check for nil signature
		return ErrNilSignature // Return error
	} else if *transaction.Recipient != chain.Account && transaction.Sender != nil && *transaction.Sender != chain.Account { // Check irrelevant
		return ErrIrrelevantTransaction // Return error
	}

	if err == nil { // Check not genesis
		signatureValid, err := VerifyTransactionSignature(transaction) // Verify signature

		if err != nil { // Check for errors
			return err // Return found error
		} else if signatureValid != true { // Check bad signature
			return ErrInvalidSignature // Return error
		}
	}

	if transaction.Recipient == nil || transaction.Recipient.String() == "0x0000000000000000000000000000000000000000" { // Check nil recipient
		return ErrNilAddress // Return error
	} else if transaction.Sender == nil && genesisChain != nil { // Check genesis already exists
		return ErrGenesisAlreadyExists // Return error
	}

	if transaction.Sender != nil { // Check not nil sender
		senderChain, err := ReadChainFromMemory(*transaction.Sender) // Read sender chain

		if err != nil { // Check for errors
			return err // Return found error
		}

		balance := senderChain.CalculateBalance() // Calculate sender balance

		if err != nil { // Check for errors
			common.Logf("== ERROR == error fetching balance for address %s %s\n", transaction.Sender.String(), err.Error()) // Log error

			return err // Return found error
		}

		if balance.Cmp(transaction.Amount) == -1 && genesisChain != nil { // Check balance insufficient
			return ErrInsufficientBalance // Return error
		}
	}

	for _, currentTransaction := range chain.Transactions { // Check for duplicate transaction
		if currentTransaction.Hash == transaction.Hash { // Check for matching hash
			return ErrDuplicateTransaction // Return error
		}
	}

	if len(chain.Transactions) == 0 && transaction.Sender == nil { // Check is genesis
		(*transaction).Genesis = true                    // Set is genesis
		chain.Genesis = *transaction.Hash                // Set genesis
		chain.Transactions = []*Transaction{transaction} // Set transaction
	} else if len(chain.Transactions) == 0 { // Check first index
		chain.Transactions = []*Transaction{transaction} // Set transaction
	} else {
		chain.Transactions = append(chain.Transactions, transaction) // Append transaction
	}

	return chain.WriteToMemory() // No error occurred, return nil
}

// QueryTransaction - attempt to fetch transaction metadata in chain by hash
func (chain *Chain) QueryTransaction(hash common.Hash) (*Transaction, error) {
	for _, transaction := range chain.Transactions { // Iterate through transactions
		if bytes.Equal(transaction.Hash.Bytes(), hash.Bytes()) { // Check for match
			return transaction, nil // Return found transaction
		}
	}

	return &Transaction{}, ErrNilTransaction
}

// QueryTransactionByParent - attempt to fetch transaction with given parent
func (chain *Chain) QueryTransactionByParent(parentHash common.Hash) (*Transaction, error) {
	for _, transaction := range chain.Transactions { // Iterate through transactions
		if bytes.Equal(transaction.ParentTx.Hash.Bytes(), parentHash.Bytes()) { // Check for match
			return transaction, nil // Return found transaction
		}
	}

	return &Transaction{}, ErrNilTransaction
}

// CalculateTargetNonce - calculate the next target nonce for the given chain.
func (chain *Chain) CalculateTargetNonce() uint64 {
	lastNonce := uint64(0) // Init nonce buffer

	for _, currentTransaction := range chain.Transactions { // Iterate through sender txs
		if currentTransaction.AccountNonce > lastNonce && bytes.Equal(currentTransaction.Sender.Bytes(), chain.Account.Bytes()) { // Check greater than last nonce
			lastNonce = currentTransaction.AccountNonce + 1 // Set last nonce
		}
	}

	return lastNonce // Return nonce
}

// CalculateBalance - iterate through tx set, return balance
func (chain *Chain) CalculateBalance() *big.Float {
	balance := big.NewFloat(0) // Init buffer

	for _, transaction := range chain.Transactions { // Iterate through transactions
		if chain.Genesis != *transaction.Hash { // Check is not genesis
			if *transaction.Sender == chain.Account { // Check is sender
				balance.Sub(balance, transaction.Amount) // Subtract value
			} else if *transaction.Recipient == chain.Account { // Check is recipient
				balance.Add(balance, transaction.Amount) // Add value
			}
		} else if chain.Genesis == *transaction.Hash { // Check is genesis
			balance.Add(balance, transaction.Amount) // Add value
		}
	}

	return balance // Return balance
}

// MakeEncodingSafe - make all transactions in chain encoding safe
func (chain *Chain) MakeEncodingSafe() error {
	for _, transaction := range chain.Transactions { // Iterate through transactions
		err := transaction.MakeEncodingSafe() // Make encoding safe

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	return nil // No error occurred, return nil
}

// RecoverSafeEncoding - recover chain from safely encoded
func (chain *Chain) RecoverSafeEncoding() error {
	for _, transaction := range chain.Transactions { // Iterate through transactions
		err := transaction.RecoverSafeEncoding() // Recover

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	return nil // No error occurred, return nil
}

// FromBytes - decode given byte array to chain
func FromBytes(b []byte) (*Chain, error) {
	chain := Chain{} // Init buffer

	err := json.NewDecoder(bytes.NewReader(b)).Decode(&chain) // Decode into buffer

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &chain, nil // No error occurred, return read value
}

// Bytes - convert given chain to byte array
func (chain *Chain) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*chain) // Serialize chain

	return buffer.Bytes() // Return serialized
}

// String - convert given chain to string
func (chain *Chain) String() string {
	marshaled, _ := json.MarshalIndent(*chain, "", "  ") // Marshal chain

	return string(marshaled) // Return marshaled
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// handleContractCall - handle given contract call
func (chain *Chain) handleContractCall(transaction *Transaction) error {
	env, err := vm.ReadEnvironmentFromMemory() // Read environment from memory

	if err != nil { // Check for errors
		env = &common.VMConfig // Get VM config

		err = env.WriteToMemory() // Write to persistent memory

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	vm, err := vm.NewVirtualMachine(chain.ContractSource, *env, new(vm.Resolver), common.GasPolicy) // Init vm

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = vm.LoadStateDB(hex.EncodeToString(vm.StateDB.ID)) // Load state db

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = vm.LoadWorkingRoot() // Load last valid working root

	if err != nil { // Check for errors
		return err // Return found error
	}

	callMethod, callParams, err := common.ParseStringMethodCallNoReceiver(string(transaction.Payload)) // Parse payload method call

	if err != nil { // Check for errors
		return err // Return found error
	}

	entryID, valid := vm.GetFunctionExport(callMethod) // Get function ID from payload

	if !valid { // Check for errors
		return ErrInvalidPayload // Return error
	}

	var parsedCallParams []int64 // Init params buffer

	for _, param := range callParams { // Iterate through params
		intVal, err := strconv.ParseInt(param, 10, 64) // Parse int

		if err != nil { // Check for errors
			return err // Return found error
		}

		parsedCallParams = append(parsedCallParams, intVal) // Append parse param
	}

	result, err := vm.Run(entryID, parsedCallParams...) // Run

	if err != nil { // Check for errors
		errLog := NewLog("error", []byte(err.Error()), Error) // Init log

		(*transaction).Logs = append((*transaction).Logs, errLog) // Append error log

		err = chain.WriteToMemory() // Write chain to persistent memory

		if err != nil { // Check for errors
			return err // Return found error
		}

		common.Logf("== ERROR == call stopped with error: %d\n", err.Error()) // Log result

		return nil // Break
	}

	resultBuffer := make([]byte, 8) // Init result buffer

	binary.LittleEndian.PutUint64(resultBuffer, uint64(result)) // Encode to []byte

	returnLog := NewLog("return", resultBuffer, Return) // Init log

	(*transaction).Logs = append((*transaction).Logs, returnLog) // Append return log TODO: custom logs, gas

	err = chain.WriteToMemory() // Write chain to persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	common.Logf("== CONTRACT == call executed successfully: %d, using %d gas\n", result, vm.Gas) // Log result

	common.Logf("== STATE == attempting to save state for contract: %s\n", chain.Account.String()) // Log save state

	err = vm.SaveState() // Save state

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// MakeGenesis - generate genesis blocks from genesis file
func (chain *Chain) MakeGenesis(genesis *config.ChainConfig) (common.Hash, error) {
	genesisTx, err := NewTransaction(0, nil, nil, &genesis.AllocAddresses[0], genesis.Alloc[genesis.AllocAddresses[0].String()], []byte("genesis")) // Init transaction

	if err != nil { // Check for errors
		return common.Hash{}, err // Return error
	}

	common.Logf("== NETWORK == initialized genesis transaction %s\n", genesisTx.Hash.String())     // Log genesis TX
	common.Logf("== NETWORK == adding genesis transaction %s to chain\n", genesisTx.Hash.String()) // Log add

	err = chain.AddTransaction(genesisTx) // Add genesis tx

	if err != nil { // Check for errors
		common.Logf("== ERROR == error making genesis block: %s\n", err.Error()) // Log error

		return common.Hash{}, err // Return found error
	}

	common.Logf("== SUCCESS == added genesis tx %s to chain %s\n", genesisTx.Hash.String(), chain.ID.String()) // Log success

	lastTx := genesisTx // Set initial

	if len(genesis.AllocAddresses) > 1 { // Check needs genesis children
		common.Log("== CHAIN == initializing gensis children") // Log genesis children

		for x := 1; x != len(genesis.AllocAddresses); x++ { // Iterate through allocations
			lastTx, err = NewTransaction(uint64(x+1), lastTx, nil, &genesis.AllocAddresses[x], genesis.Alloc[genesis.AllocAddresses[x].String()], []byte("genesisChild")) // Init transaction

			if err != nil { // Check for errors
				return common.Hash{}, err // Return error
			}

			common.Logf("== CHAIN == initialized genesis child transaction %s for alloc address %s\n", lastTx.Hash.String(), genesis.AllocAddresses[x]) // Log init

			err = chain.AddTransaction(lastTx) // Add tx

			if err != nil { // Check for errors
				return common.Hash{}, err // Return error
			}

			common.Logf("== SUCCESS == added genesis child tx %s to chain %s\n", lastTx.Hash.String(), chain.ID.String()) // Log success
		}
	}

	return *genesisTx.Hash, nil // Return genesis
}

/* END INTERNAL METHODS */
