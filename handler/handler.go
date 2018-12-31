package handler

import (
	"errors"
	"net"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/config"
	"github.com/space55/summertech-blockchain/types"
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
	common.Logf("== CONNECTION == incoming connection from peer %s\n", conn.RemoteAddr().String()) // Log conn

	data, err := common.ReadConnectionWaitAsyncNoTLS(conn) // Read data

	if err != nil { // Check for errors
		return err // Return error
	}

	defer conn.Close() // Close connection

	switch string(data)[0:9] { // Handle signatures
	case "{" + `"` + "address": // Check coordinationNode
		common.Logf("== NETWORK == received peer coordination node info %s\n", string(data)[:175]) // Log node

		return types.HandleReceivedCoordinationNode(data) // Handle received data
	case "chainRequ":
		common.Logf("== NETWORK == received chain request from peer %s\n", conn.RemoteAddr().String()) // Log request

		chain, err := types.HandleReceivedChainRequest(data) // Handle received chain request

		if err != nil { // Check for errors
			return err // Return found error
		}

		common.Logf("== NETWORK == responding to chain request from peer %s\n", conn.RemoteAddr().String()) // Log request

		_, err = conn.Write(chain.Bytes()) // Write chain

		if err != nil { // Check for errors
			return err // Return found error
		}
	case "{" + `"` + "nonce" + `"` + ":": // Check transaction
		common.Logf("== NETWORK == received transaction from peer %s\n", conn.RemoteAddr().String()) // Log tx

		return types.HandleReceivedTransaction(data) // Handle received data
	case "cChainReq": // Check coordinationChain request
		common.Logf("== NETWORK == received coordination chain request from peer %s\n", conn.RemoteAddr().String()) // Log request

		chainBytes, err := types.HandleReceivedCoordinationChainRequest() // Handle chain request

		if err != nil { // Check for errors
			return err // Return found error
		}

		common.Logf("== NETWORK == responding to coordination chain request from peer %s\n", conn.RemoteAddr().String()) // Log request

		_, err = conn.Write(chainBytes) // Write chain

		if err != nil { // Check for errors
			return err // Return found error
		}
	case "configReq":
		common.Logf("== NETWORK == received chain config request from peer %s\n", conn.RemoteAddr().String()) // Log request

		configBytes, err := config.HandleReceivedConfigRequest() // Handle config request

		if err != nil { // Check for errors
			return err // Return found error
		}

		common.Logf("== NETWORK == responding to chain config request from peer %s\n", conn.RemoteAddr().String()) // Log request

		_, err = conn.Write(configBytes) // Write chain

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	return nil // No error occurred, return nil
}

/* END INTERNAL METHODS */
