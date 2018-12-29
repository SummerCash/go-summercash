package types

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/config"
	"github.com/space55/summertech-blockchain/crypto"
)

// Chain - account transactions chain
type Chain struct {
	Account common.Address `json:"account"` // Chain account

	Transactions []*Transaction `json:"transactions"` // Transactions in chain

	Genesis common.Hash `json:"genesis"` // Genesis block hash

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

	// ErrInsufficientBalance - error definition describing a transaction worth less than the sender's balance
	ErrInsufficientBalance = errors.New("insufficient transaction sender balance")
)

/* BEGIN EXPORTED METHODS */

// NewChain - initialize new chain
func NewChain(account common.Address) (*Chain, error) {
	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain from memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	config, err := config.ReadChainConfigFromMemory() // Read config from memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	_, err = coordinationChain.QueryAddress(account) // Query address

	if err != nil && err != ErrNilNode { // Check chain with address does not exist
		return &Chain{}, err // Return error
	} else if err == nil { // Check exists
		return &Chain{}, ErrAlreadySigned // Return error
	}

	chain := &Chain{ // Init chain
		Account:      account,
		Transactions: []*Transaction{},
		NetworkID:    config.NetworkID,
	}

	localIP, err := common.GetExtIPAddrWithoutUPnP() // Get IP addr

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	if coordinationChain.Nodes == nil || len(coordinationChain.Nodes) == 0 { // Check genesis
		_, err := chain.makeGenesis(config) // Make genesis

		if err != nil { // Check for errors
			return &Chain{}, err // Return found error
		}
	}

	node, err := NewCoordinationNode(account, []string{localIP}) // Initialize node

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	foundNode, err := coordinationChain.QueryAddress(node.Address) // Check node already exists

	if err == nil { // Check already exists
		(*foundNode).Addresses = append((*foundNode).Addresses, node.Addresses[len(node.Addresses)-1]) // Append node
	} else {
		err = coordinationChain.AddNode(node, false) // Add node

		if err != nil { // Check for errors
			return &Chain{}, err // Return found error
		}
	}

	(*chain).ID = common.NewHash(crypto.Sha3(chain.Bytes())) // Set ID

	err = chain.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	return chain, nil // Return initialized chain
}

// AddTransaction - append given transaction to chain
func (chain *Chain) AddTransaction(transaction *Transaction) error {
	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain from memory

	if err != nil { // Check for errors
		return err // Return error
	}

	genesis, err := coordinationChain.GetGenesis() // Get genesis block

	if err != nil && err != ErrNilNode { // Check for errors
		return err // Return error
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

	nilCoordinationNode := CoordinationNode{} // Init nil buffer

	if transaction.Recipient == nil || transaction.Recipient.String() == "0x0000000000000000000000000000000000000000" { // Check nil recipient
		return ErrNilAddress // Return error
	} else if transaction.Sender == nil && genesis.String() != nilCoordinationNode.String() { // Check genesis already exists
		return ErrGenesisAlreadyExists // Return error
	}

	_, err = coordinationChain.GetGenesis() // Get genesis

	if err == nil && transaction.Sender == nil { // Check genesis already exists
		return ErrGenesisAlreadyExists // Return error
	}

	if transaction.Sender != nil { // Check not nil sender
		balance, err := coordinationChain.GetBalance(*transaction.Sender) // Get sender balance

		if err != nil { // Check for errors
			return err // Return found error
		}

		if balance < transaction.Amount && genesis.String() != nilCoordinationNode.String() { // Check balance insufficient
			return ErrInsufficientBalance // Return error
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

// CalculateBalance - iterate through tx set, return balance
func (chain *Chain) CalculateBalance() float64 {
	balance := float64(0) // Init buffer

	for _, transaction := range chain.Transactions { // Iterate through transactions
		if chain.Genesis != *transaction.Hash { // Check is not genesis
			if *transaction.Sender == chain.Account { // Check is sender
				balance -= transaction.Amount // Subtract value
			} else if *transaction.Recipient == chain.Account { // Check is recipient
				balance += transaction.Amount // Add value
			}
		} else if chain.Genesis == *transaction.Hash { // Check is genesis
			balance += transaction.Amount // Add value
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

// makeGenesis - generate genesis blocks from genesis file
func (chain *Chain) makeGenesis(genesis *config.ChainConfig) (common.Hash, error) {
	coordinationChain, err := ReadCoordinationChainFromMemory() // Read coordination chain from memory

	if err != nil { // Check for errors
		return common.Hash{}, err // Return error
	}

	genesisNode, _ := coordinationChain.GetGenesis() // Get genesis

	emptyNode := CoordinationNode{} // Init empty buffer

	if genesisNode.String() != emptyNode.String() { // Check genesis already exists
		return common.Hash{}, ErrGenesisAlreadyExists // Return error
	}

	genesisTx, err := NewTransaction(0, nil, nil, &genesis.AllocAddresses[0], genesis.Alloc[genesis.AllocAddresses[0].String()], []byte("genesis")) // Init transaction

	if err != nil { // Check for errors
		return common.Hash{}, err // Return error
	}

	err = chain.AddTransaction(genesisTx) // Add genesis

	if err != nil { // Check for errors
		return common.Hash{}, err // Return found error
	}

	lastTx := genesisTx // Set initial

	for x := 1; x != len(genesis.AllocAddresses); x++ { // Iterate through allocations
		lastTx, err = NewTransaction(uint64(x+1), lastTx, nil, &genesis.AllocAddresses[x], genesis.Alloc[genesis.AllocAddresses[x].String()], []byte("genesisChild")) // Init transaction

		if err != nil { // Check for errors
			return common.Hash{}, err // Return error
		}
	}

	return *genesisTx.Hash, nil // Return genesis
}

/* END INTERNAL METHODS */
