// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"bufio"
	"bytes"
	"context"

	protocol "github.com/libp2p/go-libp2p-protocol"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
)

/* BEGIN EXPORTED METHODS */

// BroadcastDht attempts to send a given message to all nodes in a dht at a given endpoint.
func BroadcastDht(ctx context.Context, host *routed.RoutedHost, message []byte, streamProtocol string, dagIdentifier string) error {
	peers := host.Peerstore().Peers() // Get peers

	for _, peer := range peers { // Iterate through peers
		if peer == (*host).ID() { // Check not same node
			continue // Continue
		}

		stream, err := (*host).NewStream(ctx, peer, protocol.ID(streamProtocol)) // Connect

		if err != nil { // Check for errors
			continue // Continue
		}

		writer := bufio.NewWriter(stream) // Initialize writer

		_, err = writer.Write(append(message, '\r')) // Write message

		if err != nil { // Check for errors
			continue // Continue
		}

		writer.Flush() // Flush
	}

	return nil // No error occurred, return nil
}

// BroadcastDhtResult send a given message to all nodes in a dht, and returns the result from each node.
func BroadcastDhtResult(ctx context.Context, host *routed.RoutedHost, message []byte, streamProtocol string, dagIdentifier string, nPeers int) ([][]byte, error) {
	peers := host.Peerstore().Peers() // Get peers

	results := [][]byte{} // Init results buffer

	for x, peer := range peers { // Iterate through peers
		if x >= nPeers { // Check has sent to enough peers
			break // Break
		}

		if peer == (*host).ID() { // Check not same node
			continue // Continue
		}

		stream, err := (*host).NewStream(ctx, peer, protocol.ID(streamProtocol)) // Connect

		if err != nil { // Check for errors
			continue // Continue
		}

		readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer

		_, err = readWriter.Write(append(message, byte('\r'))) // Write message

		if err != nil { // Check for errors
			continue // Continue
		}

		readWriter.Flush() // Flush

		responseBytes, err := readWriter.ReadBytes('\r') // Read up to delimiter

		if err != nil { // Check for errors
			continue // Continue
		}

		responseBytes = bytes.Trim(responseBytes, "\v") // Trim delmiter

		results = append(results, responseBytes) // Append response

		readWriter.Flush() // Flush
	}

	return results, nil // No error occurred, return response
}

/* END EXPORTED METHODS */
