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

	figure "github.com/common-nighthawk/go-figure"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	accountsProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/accounts"
	chainProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/chain"
	commonProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/common"
	configProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/config"
	coordinationChainProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/coordinationchain"
	cryptoProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/crypto"
	p2pProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/p2p"
	transactionProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/transaction"
	upnpProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/upnp"
)

var (
	// ErrInvalidParams - error definition describing invalid input parameters
	ErrInvalidParams = errors.New("invalid parameters")

	// workingNetwork is the working network.
	workingNetwork = "main_net"
)

// NewTerminal - attempts to start handler for term commands
func NewTerminal(rpcPort uint, rpcAddress string, network string) {
	workingNetwork = network // Set network

	reader := bufio.NewScanner(os.Stdin) // Init reader

	transport := &http.Transport{ // Init transport
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	logHeader() // Log header

	for {
		fmt.Print("\n> ") // Print prompt

		reader.Scan() // Scan

		input := reader.Text() // Fetch string input

		input = strings.TrimSuffix(input, "\n") // Trim newline

		receiver, methodname, params, err := common.ParseStringMethodCall(input) // Attempt to parse as method call
		if err != nil {                                                          // Check for errors
			fmt.Println(err.Error()) // Log found error

			continue // Continue
		}

		handleCommand(receiver, methodname, params, rpcPort, rpcAddress, transport) // Handle command
	}
}

// handleCommand - run handler for given receiver
func handleCommand(receiver string, methodname string, params []string, rpcPort uint, rpcAddress string, transport *http.Transport) {
	cryptoClient := cryptoProto.NewCryptoProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                                  // Init crypto client
	upnpClient := upnpProto.NewUpnpProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                                        // Init upnp client
	accountsClient := accountsProto.NewAccountsProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                            // Init accounts client
	configClient := configProto.NewConfigProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                                  // Init config client
	transactionClient := transactionProto.NewTransactionProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                   // Init transaction client
	chainClient := chainProto.NewChainProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                                     // Init chain client
	coordinationChainClient := coordinationChainProto.NewCoordinationChainProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport}) // Init coordinationChain client
	commonClient := commonProto.NewCommonProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                                  // Init common client
	p2pClient := p2pProto.NewP2PProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                                           // Init p2p client

	switch receiver {
	case "crypto":
		err := handleCrypto(&cryptoClient, methodname, params) // Handle crypto
		if err != nil {                                        // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "upnp":
		err := handleUpnp(&upnpClient, methodname, params) // Handle upnp
		if err != nil {                                    // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "accounts":
		err := handleAccounts(&accountsClient, methodname, params) // Handle accounts
		if err != nil {                                            // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "config":
		err := handleConfig(&configClient, methodname, params) // Handle config
		if err != nil {                                        // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "transaction":
		err := handleTransaction(&transactionClient, methodname, params) // Handle tx
		if err != nil {                                                  // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "chain":
		err := handleChain(&chainClient, methodname, params) // Handle chain
		if err != nil {                                      // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "coordinationChain":
		err := handleCoordinationChain(&coordinationChainClient, methodname, params) // Handle coordination chain
		if err != nil {                                                              // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "common":
		err := handleCommon(&commonClient, methodname, params) // Handle common
		if err != nil {                                        // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "p2p":
		err := handleP2P(&p2pClient, methodname, params) // Handle p2p
		if err != nil {                                  // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	default:
		fmt.Println("\n" + "unrecognized namespace " + `"` + receiver + `"` + ", available namespaces: crypto, upnp, accounts, config, transaction, chain, coordinationChain, common") // Log invalid namespace
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
		} else if methodname == "Sha3d" || methodname == "Sha3dString" {
			methodname = methodname[:4] + "D" + methodname[4+1:] // Correct namespace
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&cryptoProto.GeneralRequest{Input: []byte(params[0])})) // Append params
	case "Sha3n", "Sha3nString":
		if len(params) != 2 { // Check for invalid params
			return ErrInvalidParams // return error
		}

		methodname = methodname[:4] + "N" + methodname[4+1:] // Correct namespace

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
		if err != nil {                      // Check for errors
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

// handleAccounts - handle accounts receiver
func handleAccounts(accountsClient *accountsProto.Accounts, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "NewAccount", "GetAllAccounts":
		reflectParams = append(reflectParams, reflect.ValueOf(&accountsProto.GeneralRequest{})) // Append params
	case "MakeEncodingSafe", "RecoverSafeEncoding", "String", "Bytes", "ReadAccountFromMemory", "GetAllContracts":
		if len(params) != 1 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&accountsProto.GeneralRequest{Address: params[0]})) // Append params
	case "AccountFromKey":
		if len(params) != 1 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&accountsProto.GeneralRequest{PrivateKey: params[0]})) // Append params
	case "NewContractAccount":
		if len(params) != 2 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string, string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&accountsProto.GeneralRequest{Address: params[0], PrivateKey: params[1]})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: NewAccount(), NewContractAccount(), GetAllAccounts(), MakeEncodingSafe(), RecoverSafeEncoding(), String(), Bytes(), ReadAccountFromMemory()") // Return error
	}

	result := reflect.ValueOf(*accountsClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*accountsProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// handleConfig - handle config receiver
func handleConfig(configClient *configProto.Config, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "NewChainConfig":
		if len(params) != 1 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&configProto.GeneralRequest{GenesisPath: params[0]})) // Append params
	case "Bytes", "String", "WriteToMemory", "ReadChainConfigFromMemory", "GetTotalSupply", "GetInflationRate":
		if len(params) != 0 { // Check for invalid parameters
			return errors.New("invalid parameters (accepts 0 params)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&configProto.GeneralRequest{})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: NewChainConfig(), Bytes(), String(), WriteToMemory(), ReadChainConfigFromMemory(), GetTotalSupply()") // Return error
	}

	result := reflect.ValueOf(*configClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*configProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// handleTransaction - handle transaction receiver
func handleTransaction(transactionClient *transactionProto.Transaction, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "NewTransaction":
		if len(params) != 4 { // Check for invalid parameters
			return errors.New("invalid parameters (require string, string, float64, []byte)") // Return error
		}

		floatVal, _ := strconv.ParseFloat(params[2], 64) // Parse float

		reflectParams = append(reflectParams, reflect.ValueOf(&transactionProto.GeneralRequest{Address: params[0], Address2: params[1], Amount: floatVal, Payload: []byte(params[3])})) // Append params
	case "TransactionFromBytes":
		if len(params) != 1 { // Check for invalid parameters
			return errors.New("invalid parameters (require []byte)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&transactionProto.GeneralRequest{Payload: []byte(params[0])})) // Append params
	case "Bytes", "String", "SignTransaction", "VerifyTransactionSignature":
		if len(params) != 1 {
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&transactionProto.GeneralRequest{Address: params[0]})) // Append params
	case "Publish":
		if len(params) != 1 {
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&transactionProto.GeneralRequest{Address: params[0], Address2: workingNetwork})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: NewTransaction(), TransactionFromBytes(), Publish(), Bytes(), String(), SignTransaction(), VerifyTransactionSignature()") // Return error
	}

	result := reflect.ValueOf(*transactionClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*transactionProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// handleChain - handle chain receiver
func handleChain(chainClient *chainProto.Chain, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "GetBalance", "Bytes", "String", "ReadChainFromMemory", "QueryTransaction", "GetNumTransactions":
		if len(params) != 1 {
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&chainProto.GeneralRequest{Address: params[0]})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: GetBalance(), Bytes(), String(), ReadChainFromMemory(), QueryTransaction(), GetNumTransactions()") // Return error
	}

	result := reflect.ValueOf(*chainClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*chainProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// handleCoordinationChain - handle chain receiver
func handleCoordinationChain(chainClient *coordinationChainProto.CoordinationChain, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "SyncNetwork", "GetPeers", "Bytes", "String":
		if len(params) != 0 {
			return errors.New("invalid parameters (accepts no parameters)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&coordinationChainProto.GeneralRequest{})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: SyncNetwork(), GetPeers(), Bytes(), String()") // Return error
	}

	result := reflect.ValueOf(*chainClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*coordinationChainProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// handleCommon - handle common receiver
func handleCommon(commonClient *commonProto.Common, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "Encode", "EncodeString", "Decode":
		if len(params) == 0 { // Check for invalid params
			return errors.New("invalid parameters (requires []byte)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&commonProto.GeneralRequest{Input: []byte(params[0])})) // Append params
	case "DecodeString":
		if len(params) != 1 { // Check for invalid params
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&commonProto.GeneralRequest{S: params[0]}))
	default:
		return errors.New("illegal method: " + methodname + ", available methods: Encode(), EncodeString(), Decode(), DecodeString()") // Return error
	}

	result := reflect.ValueOf(*commonClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*commonProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// handleP2P - handle p2p receiver
func handleP2P(p2pClient *p2pProto.P2P, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "NumConnectedPeers", "ConnectedPeers":
		reflectParams = append(reflectParams, reflect.ValueOf(&p2pProto.GeneralRequest{})) // Empty request
	case "SyncNetwork":
		if len(params) != 1 { // Check not enough params
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&p2pProto.GeneralRequest{Network: params[0]})) // Empty request
	default:
		return errors.New("illegal method: " + methodname + ", available methods: NumConnectedPeers(), ConnectedPeers(), SyncNetwork()") // Return error
	}

	result := reflect.ValueOf(*p2pClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*p2pProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// logHeader - log contents of header file
func logHeader() {
	header := figure.NewFigure("SummerCash v"+config.Version, "slant", true) // Generate header text

	header.Print() // Log

	fmt.Println("") // Spacing
	fmt.Println("") // Spacing
}
