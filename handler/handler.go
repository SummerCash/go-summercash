package handler

import (
	"errors"
	"net"
	"strings"

	"github.com/space55/summertech-blockchain/common"
)

var (
	// ErrNilListener - error definition describing a nil listener
	ErrNilListener = errors.New("nil listener")
)

/* BEGIN EXPORTED METHODS */

// StartHandler - attempt to accept and forward requests on given listener
func StartHandler(ln *net.Listener) error {
	if ln == nil { // Check for nil listener
		return ErrNilListener // Return error
	}

	for {
		conn, err := (*ln).Accept() // Accept connection

		if err == nil { // Check for errors
			go handleConnection(conn) // Handle connection
		}
	}
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// handleConnection - attempt to handle given connection
func handleConnection(conn net.Conn) error {
	data, err := common.ReadConnectionWaitAsyncNoTLS(conn) // Read data

	if err != nil { // Check for errors
		return err // Return error
	}

	switch {
	case strings.Contains(string(data), "scope"): // Handle coordinationNode

	}

	return nil // No error occurred, return nil
}

/* END INTERNAL METHODS */
