// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SummerCash/go-summercash/accounts"

	"github.com/SummerCash/go-summercash/crypto"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
	"github.com/SummerCash/go-summercash/validator"
	protocol "github.com/libp2p/go-libp2p-protocol"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"

	commonGoP2P "github.com/dowlandaiello/GoP2P/common"
)

// Client represents a peer on the network with a known routed libp2p host.
type Client struct {
	Host *routed.RoutedHost `json:"host"` // Host

	Validator *validator.Validator `json:"validator"` // Validator

	Network string `json:"network"` // Network
}

/* BEGIN EXPORTED METHODS */

// NewClient initializes a new client with a given host.
func NewClient(host *routed.RoutedHost, validator *validator.Validator, network string) *Client {
	return &Client{
		Host:      host,      // Set host
		Validator: validator, // Set validator
		Network:   network,   // Set network
	} // Return initialized client
}

// StartIntermittentSync syncs the dag with a given context and duration.
func (client *Client) StartIntermittentSync(duration time.Duration) {
	for range time.Tick(duration) { // Sync every duration seconds
		err := client.SyncNetwork() // Sync network

		if err != nil { // Check for errors
			common.Logf("== P2P == intermittent sync errored (if private net, this is expected): %s\n", err.Error()) // Log error
		}
	}
}

// PublishTransaction publishes a given transaction.
func (client *Client) PublishTransaction(ctx context.Context, transaction *types.Transaction) error {
	peers := client.Host.Peerstore().Peers() // Get peers

	for _, peer := range peers { // Iterate through peers
		if peer == (*client.Host).ID() { // Check not same node
			continue // Continue
		}

		common.Logf("== P2P == sending tx to peer with ID %s\n", peer.Pretty()) // Log send

		stream, err := (*client.Host).NewStream(ctx, peer, protocol.ID(GetStreamHeaderProtocolPath(client.Network, PublishTransaction))) // Connect

		if err != nil { // Check for errors
			continue // Continue
		}

		common.Logf("== P2P == stream to peer %s initialized\n", peer.Pretty()) // Log stream init

		writer := bufio.NewWriter(stream) // Initialize writer

		_, err = writer.Write(append(transaction.Bytes(), '\r')) // Write message

		if err != nil { // Check for errors
			continue // Continue
		}

		common.Logf("== P2P == wrote tx to peer %s\n", peer.Pretty()) // Log stream init

		writer.Flush() // Flush
	}

	return nil // No error occurred, return nil
}

// SyncNetwork syncs all available chains and state roots.
func (client *Client) SyncNetwork() error {
	common.Logf("== P2P == starting sync...\n") // Log sync chain

	localChains, err := types.GetAllLocalizedChains() // Get all local chains

	if err != nil { // Check for errors
		return err // Return found error
	}

	common.Logf("== P2P == requesting peers for chains to sync\n") // Log sync chain

	remoteChains, err := client.RequestAllChains(16) // Request remote chains

	if err != nil { // Check for errors
		return err // Return found error
	}

	common.Logf("== P2P == found remote chains: %s\n", strings.Join(remoteChains, ", ")) // Log sync chain

	if len(remoteChains) == 0 && (*client.Validator).GetWorkingConfig() != nil { // Check no remote chains
		localAccounts, err := accounts.GetAllAccounts() // Get all accounts

		if err != nil { // Check for errors
			return err // Return found error
		}

		genesisAddress, err := common.StringToAddress(localAccounts[0]) // Get account addr

		if err != nil { // Check for errors
			return err // Return found error
		}

		chain, err := types.ReadChainFromMemory(genesisAddress) // Read genesis chain

		if err != nil { // Check for errors
			chain, err = types.NewChain(genesisAddress) // Initialize chain

			if err != nil { // Check for errors
				return err // Return found error
			}
		}

		_, err = chain.MakeGenesis((*client.Validator).GetWorkingConfig()) // Make genesis

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	for _, remoteChain := range remoteChains { // Iterate through remote chains
		common.Logf("== P2P == syncing chain %s\n", remoteChain) // Log sync chain

		if remoteChain == "" { // Check nil chain
			continue // Continue
		}

		address, err := common.StringToAddress(remoteChain) // Get address value

		if err != nil { // Check for errors
			return err // Return found error
		}

		chain, err := types.ReadChainFromMemory(address) // Read chain

		if !commonGoP2P.StringInSlice(localChains, remoteChain) || err != nil { // Check remote chain does not exist locally
			common.Logf("== P2P == chain %s does not exist locally, downloading...\n", remoteChain) // Log download chain

			chain, err = client.RequestChain(address, 8) // Request chain

			if err != nil { // Check for errors
				return err // Return found error
			}

			err = chain.WriteToMemory() // Write chain to persistent memory

			if err != nil { // Check for errors
				return err // Return found error
			}

			common.Logf("== P2P == finished downloading chain %s\n", remoteChain) // Log finish download chain
		}

		remoteBestTransaction, err := client.RequestBestTransaction(address, 16) // Request best tx

		if err != nil { // Check for errors
			return err // Return found error
		}

		common.Logf("== P2P == determined must sync up to tx with hash %s\n", remoteBestTransaction.String()) // Log sync up to

		hash := common.NewHash(crypto.Sha3(nil)) // Get nil hash

		localBestTransaction := &types.Transaction{Hash: &hash} // init local best tx buffer

		if len(chain.Transactions) != 0 { // Check chain has txs
			localBestTransaction = chain.Transactions[len(chain.Transactions)-1] // Get best tx

			common.Logf("== P2P == starting tx sync with local best tx %s\n", localBestTransaction.Hash.String()) // Log sync up to

			localBestRemoteTxInstance, err := chain.QueryTransaction(remoteBestTransaction) // Query remote on local chain

			if err == nil || bytes.Equal(remoteBestTransaction.Bytes(), common.NewHash(crypto.Sha3(nil)).Bytes()) { // Check has instance
				common.Logf("== P2P == detected possible necessary rebroadcast for tx with hash %s\n", localBestTransaction.Hash.String()) // Log sync up to

				if localBestRemoteTxInstance.Timestamp.Before(localBestTransaction.Timestamp) || bytes.Equal(remoteBestTransaction.Bytes(), common.NewHash(crypto.Sha3(nil)).Bytes()) { // Check needs re-broadcast
					common.Logf("== P2P == rebroadcasting tx with hash %s\n", localBestTransaction.Hash.String()) // Log sync up to

					rebroadcastCtx, cancel := context.WithCancel(context.Background()) // Get context

					err = client.PublishTransaction(rebroadcastCtx, localBestTransaction) // Re-broadcast

					if err != nil { // Check for errors
						cancel() // Cancel

						return err // Return found error
					}

					cancel() // Cancel

					remoteBestTransaction, err = client.RequestBestTransaction(address, 16) // Request best tx

					if err != nil { // Check for errors
						return err // Return found error
					}
				}
			}
		}

		for !bytes.Equal(localBestTransaction.Hash.Bytes(), remoteBestTransaction.Bytes()) { // Do until synced up to remote best tx
			localBestTransaction, err = client.RequestNextTransaction(*localBestTransaction.Hash, address, 16) //  Request next tx

			if err != nil { // Check for errors
				return err // Return found error
			}

			if bytes.Equal(localBestTransaction.Hash.Bytes(), remoteBestTransaction.Bytes()) { // Check synced
				break // Break
			}

			err = (*client.Validator).ValidateTransaction(localBestTransaction) // Validate tx

			if err != nil { // Check for errors
				return err // Return
			}

			err = chain.AddTransaction(localBestTransaction) // Add transaction

			if err != nil { // Check for errors
				return err // Return
			}

			err = chain.WriteToMemory() // Write to memory

			if err != nil { // Check for errors
				return err // Return found error
			}
		}

		common.Logf("== P2P == finished syncing chain %s\n", remoteChain) // Log sync up to
	}

	common.Logf("== P2P == 👏  sync finished successfully!\n") // Log sync chain

	return nil // No error occurred, return nil
}

// RequestBestTransaction requests the best transaction hash.
func (client *Client) RequestBestTransaction(account common.Address, sampleSize uint) (common.Hash, error) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	responses, err := BroadcastDhtResult(ctx, client.Host, []byte(account.String()), GetStreamHeaderProtocolPath(client.Network, RequestBestTransaction), client.Network, int(sampleSize)) // Broadcast, get result

	if err != nil { // Check for errors
		return common.Hash{}, err // Return found error
	}

	occurrences := make(map[common.Hash]int64) // Init occurrences buffer

	var bestResponse []byte // Init best response buffer

	for _, response := range responses { // Iterate through responses
		if len(response) == 0 || response == nil || bytes.Equal(response, make([]byte, len(response))) { // Check is nil
			continue // Continue
		}

		occurrences[common.NewHash(crypto.Sha3(response))]++ // Increment occurrences

		if occurrences[common.NewHash(crypto.Sha3(response))] > occurrences[common.NewHash(crypto.Sha3(bestResponse))] { // Check is better response
			bestResponse = response // Set best response
		}
	}

	if len(bestResponse) == 0 || bestResponse == nil { // Check no best response
		return common.Hash{}, errors.New("nil response") // Return error
	}

	return common.NewHash(bestResponse), nil // Return hash value
}

// RequestNextTransaction requests the next transaction with a given account chain address, and sample size.
func (client *Client) RequestNextTransaction(lastTransactionHash common.Hash, account common.Address, sampleSize uint) (*types.Transaction, error) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	responses, err := BroadcastDhtResult(ctx, client.Host, []byte(fmt.Sprintf("%s_%s", account.String(), lastTransactionHash.String())), GetStreamHeaderProtocolPath(client.Network, RequestNextTransaction), client.Network, int(sampleSize)) // Broadcast, get result

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	occurrences := make(map[common.Hash]int64) // Init occurrences buffer

	var bestResponse []byte // Init best response buffer

	for _, response := range responses { // Iterate through responses
		if len(response) == 0 || response == nil || bytes.Equal(response, make([]byte, len(response))) { // Check is nil
			continue // Continue
		}

		occurrences[common.NewHash(crypto.Sha3(response))]++ // Increment occurrences

		if occurrences[common.NewHash(crypto.Sha3(response))] > occurrences[common.NewHash(crypto.Sha3(bestResponse))] { // Check is better response
			bestResponse = response // Set best response
		}
	}

	transaction, err := types.TransactionFromBytes(bestResponse) // Get transaction value

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	return transaction, nil // Return transaction
}

// RequestAllChains requests all chain addresses from the working network with a given sample size.
func (client *Client) RequestAllChains(sampleSize uint) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	responses, err := BroadcastDhtResult(ctx, client.Host, []byte("req_all_chains"), GetStreamHeaderProtocolPath(client.Network, RequestAllChains), client.Network, int(sampleSize)) // Broadcast, get result

	if err != nil { // Check for errors
		return []string{}, err // Return found error
	}

	var bestResponse []byte // Init best response buffer

	for _, response := range responses { // Iterate through responses
		if len(response) == 0 || response == nil || bytes.Equal(response, make([]byte, len(response))) { // Check is nil
			continue // Continue
		}

		if len(response) > len(bestResponse) { // Check better response
			bestResponse = response // Set best response
		}
	}

	remoteChains := strings.Split(string(bestResponse), "_") // Split remote chain addresses

	return remoteChains, nil // Return remote chains
}

// RequestChain requests a chain from the working network with a given sample size.
func (client *Client) RequestChain(account common.Address, sampleSize uint) (*types.Chain, error) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	responses, err := BroadcastDhtResult(ctx, client.Host, account.Bytes(), GetStreamHeaderProtocolPath(client.Network, RequestChain), client.Network, int(sampleSize)) // Broadcast, get result

	if err != nil { // Check for errors
		return &types.Chain{}, err // Return found error
	}

	occurrences := make(map[common.Hash]int64) // Init occurrences buffer

	var bestResponse []byte // Init best response buffer

	for _, response := range responses { // Iterate through responses
		if len(response) == 0 || response == nil || bytes.Equal(response, make([]byte, len(response))) { // Check is nil
			continue // Continue
		}

		occurrences[common.NewHash(crypto.Sha3(response))]++ // Increment occurrences

		if occurrences[common.NewHash(crypto.Sha3(response))] > occurrences[common.NewHash(crypto.Sha3(bestResponse))] { // Check is better response
			bestResponse = response // Set best response
		}
	}

	chain, err := types.FromBytes(bestResponse) // Deserialize chain

	if err != nil { // Check for errors
		return &types.Chain{}, err // Return found error
	}

	return chain, nil // Return chain
}

/* END EXPORTED METHODS */
