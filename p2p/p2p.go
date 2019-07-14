// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/SummerCash/go-summercash/common"

	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"

	"github.com/SummerCash/go-summercash/config"
)

// ErrTimedOut defines an error describing a standard timeout.
var ErrTimedOut = errors.New("request timed out")

/* BEGIN EXPORTED METHODS */

// CheckPeerCompatible checks that a given peer is compatible with the working host.
func CheckPeerCompatible(ctx context.Context, host *routed.RoutedHost, peer peer.ID, network string) bool {
	if peer == (*host).ID() { // Check same node
		return false // Not compatible
	}

	stream, err := (*host).NewStream(ctx, peer, protocol.ID(GetStreamHeaderProtocolPath(network, RequestAlive))) // Initialize stream
	if err != nil {                                                                                              // Check for errors
		return false // Not compatible
	}

	reader := bufio.NewReader(stream) // Initialize reader

	networkBytes, err := common.ReadAll(reader) // Read network
	if err != nil {                             // Check for errors
		return false // Not compatible
	}

	networkBytes = bytes.Replace(networkBytes, []byte("\n"), []byte{}, 1) // Remove delimiter

	if string(networkBytes) != fmt.Sprintf("despacito: %s", config.Version) { // Check incompatible
		return false // Not compatible
	}

	return true // Compatible
}

// BroadcastDht attempts to send a given message to all nodes in a dht at a given endpoint.
func BroadcastDht(ctx context.Context, host *routed.RoutedHost, message []byte, streamProtocol string, dagIdentifier string) error {
	if bytes.Contains(message, []byte("\n")) { // Check control char
		return errors.New("message contains a restricted control character") // Return error
	}

	peers := host.Network().Peers() // Get peers

	for _, peer := range peers { // Iterate through peers
		if peer == (*host).ID() || !CheckPeerCompatible(ctx, host, peer, dagIdentifier) { // Check not same node, compatible
			continue // Continue
		}

		stream, err := (*host).NewStream(ctx, peer, protocol.ID(streamProtocol)) // Connect
		if err != nil {                                                          // Check for errors
			continue // Continue
		}

		writer := bufio.NewWriter(stream) // Initialize writer

		_, err = writer.Write(append(message, []byte("\n")...)) // Write message

		if err != nil { // Check for errors
			continue // Continue
		}

		writer.Flush() // Flush
	}

	return nil // No error occurred, return nil
}

// BroadcastDhtResult send a given message to all nodes in a dht, and returns the result from each node.
func BroadcastDhtResult(ctx context.Context, host *routed.RoutedHost, message []byte, streamProtocol string, dagIdentifier string, nPeers int) (responses [][]byte, err error) {
	if bytes.Contains(message, []byte("\n")) { // Check control char
		return nil, errors.New("message contains a restricted control character") // Return error
	}

	peers := host.Network().Peers() // Get peers

	results := [][]byte{} // Init results buffer

	x := 0 // Init x buffer

	for _, currentPeer := range peers { // Iterate through peers
		if x >= int(math.Ceil((float64(nPeers)/100)*float64(len(peers)))) { // Check has sent to enough peers
			break // Break
		}

		go func(peer peer.ID) {
			if peer == (*host).ID() || !CheckPeerCompatible(ctx, host, peer, dagIdentifier) { // Check not same node, compatible
				return // Continue
			}

			stream, err := (*host).NewStream(ctx, peer, protocol.ID(streamProtocol)) // Connect
			if err != nil {                                                          // Check for errors
				return // Continue
			}

			readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer

			_, err = readWriter.Write(append(message, []byte("\n")...)) // Write message

			if err != nil { // Check for errors
				return // Continue
			}

			readWriter.Flush() // Flush

			responseBytes, err := common.ReadAll(readWriter.Reader) // Read up to delimiter
			if err != nil {                                         // Check for errors
				return // Continue
			}

			responseBytes = bytes.Trim(responseBytes, "\n") // Trim delmiter

			results = append(results, responseBytes) // Append response

			readWriter.Flush() // Flush

			x++ // Increment
		}(currentPeer) // Run
	}

	startTime := time.Now() // Get start time

	for x < int(math.Ceil((float64(nPeers)/100)*float64(len(peers)))) { // Wait until enough responses
		if time.Now().Sub(startTime) > 10*time.Second { // Check for timeout
			return [][]byte{}, ErrTimedOut // Return error
		}
	}

	return results, nil // No error occurred, return response
}

/* END EXPORTED METHODS */
