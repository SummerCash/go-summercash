// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
	multiaddr "github.com/multiformats/go-multiaddr"

	"github.com/SummerCash/go-summercash/config"
)

// BootstrapNodes represents all default bootstrap nodes on the given network.
var BootstrapNodes = []string{
	"/ip4/108.41.124.60/tcp/3000/ipfs/QmWRdp5HQ1SfENPnLFrXviY9gr6Y5BcWYbkWuqDkZSizAj",
	"/ip4/108.41.124.60/tcp/3003/ipfs/QmYEaiPPsqtPes1AtPZA1zpqozmDzfgFJv3ygjTjLVXedp",
	"/ip4/174.129.191.246/tcp/3000/ipfs/QmUXThFht8qoZGdKLMVmr8Bk34VJ9oS3WWw3a25jeZucYd",
	"/ip4/54.234.2.165/tcp/3000/ipfs/Qma3QTswnKK48gzsaVhzakPApY4kPuBGsZLS1i837gns2s",
}

/* BEGIN EXPORTED METHODS */

// GetBestBootstrapAddress attempts to fetch the best bootstrap node.
func GetBestBootstrapAddress(ctx context.Context, host *routed.RoutedHost, network string) string {
	for _, bootstrapAddress := range BootstrapNodes { // Iterate through bootstrap nodes
		multiaddr, err := multiaddr.NewMultiaddr(bootstrapAddress) // Parse address
		if err != nil {                                            // Check for errors
			continue // Continue
		}

		peerID, err := peer.IDB58Decode(strings.Split(bootstrapAddress, "ipfs/")[1]) // Get peer ID
		if err != nil {                                                              // Check for errors
			continue // Continue
		}

		peerInfo, err := peerstore.InfoFromP2pAddr(multiaddr) // Get peer info
		if err != nil {                                       // Check for errors
			continue // Continue
		}

		bootstrapCheckCtx, cancel := context.WithCancel(ctx) // Get context

		err = host.Connect(bootstrapCheckCtx, *peerInfo) // Connect to peer

		if err != nil { // Check for errors
			cancel() // Cancel

			continue // Continue
		}

		stream, err := (*host).NewStream(ctx, peerInfo.ID, protocol.ID(GetStreamHeaderProtocolPath(network, RequestAlive))) // Get stream
		if err != nil {                                                                                                     // Check for errors
			cancel() // Cancel

			continue // Continue
		}

		host.Peerstore().AddAddr(peerID, multiaddr, peerstore.PermanentAddrTTL) // Add bootstrap peer

		reader := bufio.NewReader(stream) // Get reader

		errChan := make(chan error) // Init error buffer
		doneChan := make(chan bool) // Init done buffer

		timer := time.NewTimer(time.Second * time.Duration(15)) // Init timer

		go func() {
			network, err := common.ReadAll(reader) // Read
			if err != nil {                        // Check for errors
				err = host.Network().ClosePeer(peerInfo.ID) // Disconnect from peer

				if err != nil { // Check for errors
					errChan <- err // Write err

					return // Return
				}

				errChan <- err // Write err

				return // Return
			}

			network = bytes.Replace(network, []byte("\n\r"), []byte{}, 1) // Remove delimiter

			if string(network) != fmt.Sprintf("despacito: %s", config.Version) { // Check networks not matching
				errChan <- fmt.Errorf("network not matching for peer with multi-addr: %s", peerInfo.ID.Pretty()) // Write err
			}

			doneChan <- true // Done
		}()

		select {
		case <-doneChan:
			cancel() // Cancel

			return bootstrapAddress // Done!
		case <-errChan:
			cancel() // Cancel

			continue // Continue
		case <-timer.C:
			cancel() // Cancel

			continue // Continue
		}
	}

	return "localhost" // Return localhost
}

/* END EXPORTED METHODS */
