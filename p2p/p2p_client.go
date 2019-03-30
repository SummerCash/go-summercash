// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/SummerCash/go-summercash/crypto"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
	"github.com/SummerCash/go-summercash/validator"
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

// SyncNetwork syncs all available chains and state roots.
func (client *Client) SyncNetwork() error {
	localChains, err := types.GetAllLocalizedChains() // Get all local chains

	if err != nil { // Check for errors
		return err // Return found error
	}

	remoteChains, err := client.RequestAllChains(16) // Request remote chains

	if err != nil { // Check for errors
		return err // Return found error
	}

	for _, remoteChain := range remoteChains { // Iterate through remote chains
		address, err := common.StringToAddress(remoteChain) // Get address value

		if err != nil { // Check for errors
			return err // Return found error
		}

		chain, err := types.ReadChainFromMemory(address) // Read chain

		if !commonGoP2P.StringInSlice(localChains, remoteChain) || err != nil { // Check remote chain does not exist locally
			chain, err := client.RequestChain(address, 8) // Request chain

			if err != nil { // Check for errors
				return err // Return found error
			}

			err = chain.WriteToMemory() // Write chain to persistent memory

			if err != nil { // Check for errors
				return err // Return found error
			}
		}

		remoteBestTransaction, err := client.RequestBestTransaction(address, 16) // Request best tx

		if err != nil { // Check for errors
			return err // Return found error
		}

		localBestTransaction := chain.Transactions[len(chain.Transactions)-1] // Get best tx

		for !bytes.Equal(localBestTransaction.Hash.Bytes(), remoteBestTransaction.Bytes()) { // Do until synced up to remote best tx

		}
	}

	return nil // No error occurred, return nil
}

// RequestBestTransaction requests the best transaction hash.
func (client *Client) RequestBestTransaction(account common.Address, sampleSize uint) (common.Hash, error) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	responses, err := BroadcastDhtResult(ctx, client.Host, append([]byte(account.String()), '\f'), GetStreamHeaderProtocolPath(client.Network, RequestNextTransaction), client.Network, int(sampleSize)) // Broadcast, get result

	if err != nil { // Check for errors
		return common.Hash{}, err // Return found error
	}

	occurrences := make(map[common.Hash]int64) // Init occurrences buffer

	var bestResponse []byte // Init best response buffer

	for _, response := range responses { // Iterate through responses
		if len(response) == 0 || response == nil || bytes.Equal(response, make([]byte, len(response))) { // Check is nil
			continue // Continue
		}

		occurrences[common.NewHash(response)]++ // Increment occurrences

		if occurrences[common.NewHash(response)] > occurrences[common.NewHash(bestResponse)] { // Check is better response
			bestResponse = response // Set best response
		}
	}

	return common.NewHash(bestResponse), nil // Return hash value
}

// RequestNextTransaction requests the next transaction with a given account chain address, and sample size.
func (client *Client) RequestNextTransaction(lastTransactionHash common.Hash, account common.Address, sampleSize uint) (*types.Transaction, error) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	responses, err := BroadcastDhtResult(ctx, client.Host, append([]byte(fmt.Sprintf("%s_%s", account.String(), lastTransactionHash.String())), '\f'), GetStreamHeaderProtocolPath(client.Network, RequestNextTransaction), client.Network, int(sampleSize)) // Broadcast, get result

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

	responses, err := BroadcastDhtResult(ctx, client.Host, append([]byte("req_all_chains"), '\f'), GetStreamHeaderProtocolPath(client.Network, RequestChain), client.Network, int(sampleSize)) // Broadcast, get result

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
