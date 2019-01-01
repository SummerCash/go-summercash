package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/SummerCash/go-summercash/cli"
	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/handler"
	accountsServer "github.com/SummerCash/go-summercash/internal/rpc/accounts"
	configServer "github.com/SummerCash/go-summercash/internal/rpc/config"
	cryptoServer "github.com/SummerCash/go-summercash/internal/rpc/crypto"
	accountsProto "github.com/SummerCash/go-summercash/internal/rpc/proto/accounts"
	configProto "github.com/SummerCash/go-summercash/internal/rpc/proto/config"
	cryptoProto "github.com/SummerCash/go-summercash/internal/rpc/proto/crypto"
	transactionProto "github.com/SummerCash/go-summercash/internal/rpc/proto/transaction"
	upnpProto "github.com/SummerCash/go-summercash/internal/rpc/proto/upnp"
	transactionServer "github.com/SummerCash/go-summercash/internal/rpc/transaction"
	upnpServer "github.com/SummerCash/go-summercash/internal/rpc/upnp"
	"github.com/SummerCash/go-summercash/types"
	"github.com/SummerCash/go-summercash/upnp"
	commonGoP2P "github.com/dowlandaiello/GoP2P/common"
)

var (
	terminalFlag       = flag.Bool("terminal", false, "launch node in terminal mode")                                                                                     // Init term flag
	upnpFlag           = flag.Bool("no-upnp", false, "launch node without automatic UPnP port forwarding")                                                                // Init upnp flag
	rpcPortFlag        = flag.Int("rpc-port", 8080, "launch node with specified RPC port")                                                                                // Init RPC port flag
	forwardRPCFlag     = flag.Bool("forward-rpc", false, "enables forwarding of node RPC terminal ports")                                                                 // Init forward RPC flag
	rpcAddrFlag        = flag.String("rpc-address", fmt.Sprintf("localhost:%s", strconv.Itoa(*rpcPortFlag)), "connects to remote RPC terminal (default: localhost:8080)") // Init remote rpc addr flag
	dataDirFlag        = flag.String("data-dir", common.DataDir, "performs all node i/o operations in given data directory")                                              // Init data dir flag
	nodePortFlag       = flag.Int("node-port", common.NodePort, "launch node on give port")                                                                               // Init node port flag
	privateNetworkFlag = flag.Bool("private-net", false, "launch node in context of private network")                                                                     // Init private network flag
	archivalNodeFlag   = flag.Bool("archival", false, "launch node in archival mode")                                                                                     // Init archival node flag
	silent             = flag.Bool("silent", false, "silence all fmt.Print calls")                                                                                        // Init silent flag
)

func main() {
	flag.Parse() // Parse flags

	common.DataDir = *dataDirFlag   // Set data-dir
	common.Silent = *silent         // Set is silent
	common.NodePort = *nodePortFlag // Set node port

	if *privateNetworkFlag {
		common.ExtIPProviders = []string{} // Set nil providers
	}

	if !*upnpFlag { // Check for UPnP
		if *forwardRPCFlag {
			go upnp.ForwardPortSilent(uint(*rpcPortFlag)) // Forward RPC port
		}

		go upnp.ForwardPortSilent(uint(*nodePortFlag)) // Forward port 3000
	}

	if strings.Contains(*rpcAddrFlag, "localhost") { // Check for default RPC address
		startRPCServer() // Start RPC server

		go startNode(*archivalNodeFlag) // Start node
	}

	if *terminalFlag { // Check for terminal
		*rpcAddrFlag = strings.Split(*rpcAddrFlag, ":")[0] // Remove port

		cli.NewTerminal(uint(*rpcPortFlag), *rpcAddrFlag) // Initialize terminal
	}

	go common.Forever() // Prevent main from closing
	select {}           // Prevent main from closing
}

// startRPCServer - start RPC server
func startRPCServer() {
	err := common.GenerateTLSCertificates("term") // Generate certs

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	cryptoHandler := cryptoProto.NewCryptoServer(&cryptoServer.Server{}, nil)                     // Init handler
	upnpHandler := upnpProto.NewUpnpServer(&upnpServer.Server{}, nil)                             // Init handler
	accountsHandler := accountsProto.NewAccountsServer(&accountsServer.Server{}, nil)             // Init handler
	configHandler := configProto.NewConfigServer(&configServer.Server{}, nil)                     // Init handler
	transactionHandler := transactionProto.NewTransactionServer(&transactionServer.Server{}, nil) // Init handler

	mux := http.NewServeMux() // Init mux

	mux.Handle(cryptoProto.CryptoPathPrefix, cryptoHandler)                // Start mux node handler
	mux.Handle(upnpProto.UpnpPathPrefix, upnpHandler)                      // Start mux upnp handler
	mux.Handle(accountsProto.AccountsPathPrefix, accountsHandler)          // Start mux accounts handler
	mux.Handle(configProto.ConfigPathPrefix, configHandler)                // Start mux config handler
	mux.Handle(transactionProto.TransactionPathPrefix, transactionHandler) // Start mux config handler

	go http.ListenAndServeTLS(":"+strconv.Itoa(*rpcPortFlag), "termCert.pem", "termKey.pem", mux) // Start server
}

// startNode - start necessary services for full node
func startNode(archivalNode bool) {
	ip, _ := common.GetExtIPAddrWithoutUPnP() // Get IP

	common.Logf("== NODE == starting on port %d with external IP %s\n", *nodePortFlag, ip) // Log init

	coordinationChain, err := types.ReadCoordinationChainFromMemory() // Read coordination chain

	if err == nil { // Check no error
		_, err = coordinationChain.QueryArchivalNode(ip) // Set error
	}

	if strings.Contains(ip, ":") { // Check is IPv6
		ip = "[" + ip + "]" + ":" + strconv.Itoa(common.NodePort) // Add port
	} else {
		ip = ip + ":" + strconv.Itoa(common.NodePort) // Add port
	}

	if err != nil { // Check for errors
		if archivalNode && !commonGoP2P.StringInSlice(common.BootstrapNodes, ip) && !*privateNetworkFlag { // Check is not bootstrap node
			common.Logf("== NETWORK == joining with bootstrap node %s\n", common.BootstrapNodes[0]) // Log join

			err := types.JoinNetwork(common.BootstrapNodes[0], true) // Register node

			if err != nil { // Check for errors
				panic(err) // Panic
			}
		} else if !commonGoP2P.StringInSlice(common.BootstrapNodes, ip) { // Plz, no recursion TODO: fix ipv6
			err := types.SyncNetwork() // Sync network

			if err != nil { // Check for errors
				panic(err) // Panic
			}
		}
	}

	ln, err := tls.Listen("tcp", ":"+strconv.Itoa(*nodePortFlag), common.GeneralTLSConfig) // Listen on port

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	err = handler.StartHandler(&ln) // Start handler

	if err != nil { // Check for errors
		panic(err) // Panic
	}
}

/*
	TODO:
	- readme
*/
