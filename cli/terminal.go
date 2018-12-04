package cli

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/space55/summertech-blockchain/common"
	cryptoProto "github.com/space55/summertech-blockchain/internal/rpc/proto/crypto"
	upnpProto "github.com/space55/summertech-blockchain/internal/rpc/proto/upnp"
)

var (
	// ErrInvalidParams - error definition describing invalid input parameters
	ErrInvalidParams = errors.New("invalid parameters")
)

// NewTerminal - attempts to start handler for term commands
func NewTerminal(rpcPort uint, rpcAddress string) {
	reader := bufio.NewScanner(os.Stdin) // Init reader

	transport := &http.Transport{ // Init transport
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	for {
		fmt.Print("\n> ") // Print prompt

		reader.Scan() // Scan

		input := reader.Text() // Fetch string input

		input = strings.TrimSuffix(input, "\n") // Trim newline

		receiver, methodname, params, err := common.ParseStringMethodCall(input) // Attempt to parse as method call

		if err != nil { // Check for errors
			fmt.Println(err.Error()) // Log found error

			continue
		}

		handleCommand(receiver, methodname, params, rpcPort, rpcAddress, transport) // Handle command
	}
}

// handleCommand - run handler for given receiver
func handleCommand(receiver string, methodname string, params []string, rpcPort uint, rpcAddress string, transport *http.Transport) {
	cryptoClient := cryptoProto.NewCryptoProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport}) // Init crypto client
	upnpClient := upnpProto.NewUpnpProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})       // Init upnp client

	switch receiver {
	case "crypto":
		err := handleCrypto(&cryptoClient, methodname, params) // Handle crypto

		if err != nil { // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "upnp":
		err := handleUpnp(&upnpClient, methodname, params) // Handle upnp

		if err != nil { // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	}
}

// handleCrypto - handle crypto receiver
func handleCrypto(cryptoClient *cryptoProto.Crypto, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname { // Handle different methods
	case "Sha3", "Sha3String", "Sha3d", "Sha3dString":
		if len(params) != 1 { // Check for invalid params
			return ErrInvalidParams // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&cryptoProto.GeneralRequest{Input: []byte(params[0])})) // Append params
	case "Sha3n", "Sha3nString":
		if len(params) != 2 { // Check for invalid params
			return ErrInvalidParams // return error
		}

		intVal, _ := strconv.Atoi(params[1]) // Convert to int

		reflectParams = append(reflectParams, reflect.ValueOf(&cryptoProto.GeneralRequest{Input: []byte(params[0]), N: uint32(intVal)})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: Sha3(), Sha3String(), Sha3d(), Sha3dString(), Sha3n(), Sha3nString()") // Return error
	}

	result := reflect.ValueOf(*cryptoClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*cryptoProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// handleUpnp - handle upnp receiver
func handleUpnp(upnpClient *upnpProto.Upnp, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "GetGateway":
		reflectParams = append(reflectParams, reflect.ValueOf(&upnpProto.GeneralRequest{})) // Append params
	case "ForwardPortSilent", "ForwardPort", "RemoveForwarding":
		if len(params) != 1 { // Check for invalid parameters
			return errors.New("invalid parameters (requires uint32)") // Return error
		}

		port, err := strconv.Atoi(params[0]) // Convert to int

		if err != nil { // Check for errors
			return err // Return found error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&upnpProto.GeneralRequest{PortNumber: uint32(port)})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: GetGateway(), ForwardPortSilent(), ForwardPort(), RemoveForwarding()") // Return error
	}

	result := reflect.ValueOf(*upnpClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*upnpProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}
