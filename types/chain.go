package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/crypto"
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

	// ErrInsufficientBalance - error definition describing a transaction worth more than the sender's balance
	ErrInsufficientBalance = errors.New("insufficient transaction sender balance")

	// ErrNilTransaction - error definition describing a query for a non-existent transaction
	ErrNilTransaction = errors.New("couldn't find transaction with given hash")

	// ErrDuplicateTransaction - error definition describing a transaction that already exists in a given chain
	ErrDuplicateTransaction = errors.New("duplicate transaction")
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

	(*chain).ID = common.NewHash(crypto.Sha3(chain.Bytes())) // Set ID

	common.Logf("== ACCOUNT == initialized account chain with account address %s\n", chain.Account.String()) // Log init

	localIP, err := common.GetExtIPAddrWithoutUPnP() // Get IP addr

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	if coordinationChain.Nodes == nil || len(coordinationChain.Nodes) == 0 { // Check genesis
		common.Log("== NETWORK == making genesis block") // Log genesis block

		_, err := chain.makeGenesis(config) // Make genesis

		if err != nil { // Check for errors
			return &Chain{}, err // Return found error
		}
	}

	common.Log("== CHAIN == initializing account chain coordination node") // Log coordination node init

	node, err := NewCoordinationNode(account, []string{localIP + ":" + strconv.Itoa(common.NodePort)}) // Initialize node

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	common.Logf("== SUCCESS == initialized account chain coordination node %s\n", node.ID.String()) // Log coordination node init

	foundNode, err := coordinationChain.QueryAddress(node.Address) // Check node already exists

	if err == nil { // Check already exists
		(*foundNode).Addresses = append((*foundNode).Addresses, node.Addresses[len(node.Addresses)-1]) // Append node
	} else {
		if coordinationChain.Nodes == nil || len(coordinationChain.Nodes) == 0 { // Check is genesis
			common.Log("== CHAIN == appending genesis coordination node to local chain\n") // Log coordination node init

			err = coordinationChain.AddNode(node, false) // Add node

			if err != nil { // Check for errors
				return &Chain{}, err // Return found error
			}

			common.Logf("== SUCCESS == appended node %s to local chain\n", node.ID.String()) // Log coordination node init
		} else {
			common.Log("== CHAIN == appending new coordination node to remote chains") // Log coordination node append

			err = coordinationChain.AddNode(node, true) // Add node

			if err != nil { // Check for errors
				return &Chain{}, err // Return found error
			}

			common.Logf("== SUCCESS == appended node %s to remote chain\n", node.ID.String()) // Log coordination node append
		}
	}

	if coordinationChain.Nodes != nil || len(coordinationChain.Nodes) > 0 { // Check not genesis
		common.Logf("== CHAIN == pushing chain %s to network\n", chain.ID.String()) // Log push

		nodes, err := coordinationChain.QueryAllArchivalNodes() // Query all archival nodes

		if err != nil { // Check for errors
			common.Logf("== ERROR == error pushing chain to network %s\n", err.Error()) // Log error

			return nil, err // Return found error
		}

		common.Logf("== NETWORK == found %d peers to push chain to\n", len(nodes)) // Log push

		if len(nodes) > 0 {
			if nodes[0] != localIP+":"+strconv.Itoa(common.NodePort) { // Check not current node
				common.Logf("== NETWORK == pushing chain %s to peer %s\n", chain.ID.String(), nodes[0]) // Log push

				common.SendBytes(chain.Bytes(), nodes[0]) // Send chain
			}

			for x, address := range nodes { // Iterate through addresses
				if x != 0 && address != localIP+":"+strconv.Itoa(common.NodePort) { // Check not first index and not current node
					common.Logf("== NETWORK == pushing chain %s to peer %s\n", chain.ID.String(), address) // Log push

					go common.SendBytes(chain.Bytes(), address) // Send chain
				}
			}
		}
	}

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
			common.Logf("== ERROR == error fetching balance for address %s %s", transaction.Sender.String(), err.Error()) // Log error

			return err // Return found error
		}

		if balance < transaction.Amount && genesis.String() != nilCoordinationNode.String() { // Check balance insufficient
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
		if *transaction.Hash == hash { // Check for match
			return transaction, nil // Return found transaction
		}
	}

	return &Transaction{}, ErrNilTransaction
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
		common.Log("== ERROR == genesis block already exists") // Log already exists

		return common.Hash{}, ErrGenesisAlreadyExists // Return error
	}

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
