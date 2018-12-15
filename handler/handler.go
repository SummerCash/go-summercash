package handler

import (
	"errors"
	"fmt"
	"net"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/types"
)

var (
	// ErrNilListener - error definition describing a nil listener
	ErrNilListener = errors.New("nil listener")
)

/* BEGIN EXPORTED METHODS */

// StartHandler - attempt to accept and forward requests on given listener
func StartHandler(ln *net.Listener) error {
	fmt.Println("test")
	if ln == nil { // Check for nil listener
		return ErrNilListener // Return error
	}

	for {
		fmt.Println("DESPACITO DESPACITO")
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

	defer conn.Close() // Close connection

	switch string(data)[0:9] { // Handle signatures
	case "{" + `"` + "scope" + `"` + ":": // Check coordinationNode
		return types.HandleReceivedCoordinationNode(data) // Handle received data
	case "chainRequ":
		chain, err := types.HandleReceivedChainRequest(data) // Handle received chain request

		if err != nil { // Check for errors
			return err // Return found error
		}

		_, err = conn.Write(chain.Bytes()) // Write chain

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	return nil // No error occurred, return nil
}

/* END INTERNAL METHODS */
