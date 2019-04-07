package handler

import (
	"errors"
	"net"
	"strings"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/types"
)

var (
	// ErrNilListener - error definition describing a nil listener
	ErrNilListener = errors.New("nil listener")

	// ErrInvalidConnectionHeader - error definition describing a nil connection header
	ErrInvalidConnectionHeader = errors.New("invalid connection header")
)

/* BEGIN EXPORTED METHODS */

// StartHandler - attempt to accept and forward requests on given listener
func StartHandler(ln *net.Listener, isArchival bool) error {
	if ln == nil { // Check for nil listener
		return ErrNilListener // Return error
	}

	for {
		conn, err := (*ln).Accept() // Accept connection

		if err == nil { // Check for errors
			go handleConnection(conn, isArchival) // Handle connection
		}
	}
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// handleConnection - attempt to handle given connection
func handleConnection(conn net.Conn, isArchival bool) error {
	common.Logf("== CONNECTION == incoming connection from peer %s\n", conn.RemoteAddr().String()) // Log conn

	data, err := common.ReadConnectionWaitAsyncNoTLS(conn) // Read data

	if err != nil { // Check for errors
		return err // Return error
	}

	if len(string(data)) < 9 { // Check invalid input
		common.Logf("== ERROR == connection data %s invalid (does not contain connection header, or is nil)", string(data)) // Log error

		return ErrInvalidConnectionHeader // Return found error
	}

	switch string(data)[0:9] { // Handle signatures
	case "{" + `"` + "address": // Check coordinationNode
		common.Logf("== NETWORK == received peer coordination node info %s\n", string(data)) // Log node

		err = types.HandleReceivedCoordinationNode(data, isArchival) // Handle received data

		if err != nil { // Check for errors
			return err // Return found error
		}

		return conn.Close() // Close connection
	case "chainRequ":
		common.Logf("== NETWORK == received chain request from peer %s\n", conn.RemoteAddr().String()) // Log request

		chain, err := types.HandleReceivedChainRequest(data) // Handle received chain request

		if err != nil { // Check for errors
			common.Logf("== ERROR == error handling chain request from peer %s %s\n", conn.RemoteAddr().String(), err.Error()) // Log request

			return err // Return found error
		}

		common.Logf("== NETWORK == responding to chain request from peer %s\n", conn.RemoteAddr().String()) // Log request

		_, err = conn.Write(append(chain.Bytes(), byte('\b'))) // Write chain

		if err != nil { // Check for errors
			return err // Return found error
		}

		return conn.Close() // Close connection
	case "{" + `"` + "nonce" + `"` + ":": // Check transaction
		if strings.Contains(string(data), `"`+"is-init-contract"+`"`+":true") { // Check is contract creation
			common.Logf("== NETWORK == received contract creation from peer %s\n", conn.RemoteAddr().String()) // Log tx

			err = types.HandleReceivedContractCreation(data) // Handle received data

			if err != nil { // Check for errors
				return err // Return found error
			}

			return conn.Close() // Close connection
		}

		common.Logf("== NETWORK == received transaction from peer %s\n", conn.RemoteAddr().String()) // Log tx

		err = types.HandleReceivedTransaction(data) // Handle received data

		if err != nil { // Check for errors
			return err // Return found error
		}

		return conn.Close() // Close connection
	case "cChainReq": // Check coordinationChain request
		common.Logf("== NETWORK == received coordination chain request from peer %s\n", conn.RemoteAddr().String()) // Log request

		chainBytes, err := types.HandleReceivedCoordinationChainRequest() // Handle chain request

		if err != nil { // Check for errors
			common.Logf("== ERROR == error handling coordination chain request from peer %s %s\n", conn.RemoteAddr().String(), err.Error()) // Log request

			return err // Return found error
		}

		common.Logf("== NETWORK == responding to coordination chain request from peer %s\n", conn.RemoteAddr().String()) // Log request

		_, err = conn.Write(append(chainBytes, byte('\b'))) // Write chain

		if err != nil { // Check for errors
			return err // Return found error
		}

		return conn.Close() // Close connection
	case "configReq":
		common.Logf("== NETWORK == received chain config request from peer %s\n", conn.RemoteAddr().String()) // Log request

		configBytes, err := config.HandleReceivedConfigRequest() // Handle config request

		if err != nil { // Check for errors
			common.Logf("== ERROR == error handling chain config request from peer %s %s\n", conn.RemoteAddr().String(), err.Error()) // Log request

			return err // Return found error
		}

		common.Logf("== NETWORK == responding to chain config request from peer %s\n", conn.RemoteAddr().String()) // Log request

		_, err = conn.Write(append(configBytes, byte('\b'))) // Write chain

		if err != nil { // Check for errors
			return err // Return found error
		}

		return conn.Close() // Close connection
	case "stateRequ":
		common.Logf("== NETWORK == received state request from peer %s\n", conn.RemoteAddr().String()) // Log request

		state, err := types.HandleReceivedStateRequest(data) // Handle received state request

		if err != nil { // Check for errors
			common.Logf("== ERROR == error handling state request from peer %s %s\n", conn.RemoteAddr().String(), err.Error()) // Log request

			return err // Return found error
		}

		common.Logf("== NETWORK == responding to state request from peer %s\n", conn.RemoteAddr().String()) // Log request

		_, err = conn.Write(append(state, byte('\b'))) // Write chain

		if err != nil { // Check for errors
			return err // Return found error
		}

		return conn.Close() // Close connection
	case "{" + `"` + "account":
		common.Logf("== NETWORK == received account chain from peer %s\n", conn.RemoteAddr().String()) // Log post

		err = types.HandleReceivedChain(data) // Handle received data

		if err != nil { // Check for errors
			return err // Return found error
		}

		return conn.Close() // Close connection
	}

	return conn.Close() // No error occurred, return closeConn()
}

/* END INTERNAL METHODS */
