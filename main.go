package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/space55/summertech-blockchain/cli"
	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/handler"
	accountsServer "github.com/space55/summertech-blockchain/internal/rpc/accounts"
	cryptoServer "github.com/space55/summertech-blockchain/internal/rpc/crypto"
	accountsProto "github.com/space55/summertech-blockchain/internal/rpc/proto/accounts"
	cryptoProto "github.com/space55/summertech-blockchain/internal/rpc/proto/crypto"
	upnpProto "github.com/space55/summertech-blockchain/internal/rpc/proto/upnp"
	upnpServer "github.com/space55/summertech-blockchain/internal/rpc/upnp"
	"github.com/space55/summertech-blockchain/types"
	"github.com/space55/summertech-blockchain/upnp"
)

var (
	terminalFlag       = flag.Bool("terminal", false, "launch node in terminal mode")                                                                                     // Init term flag
	upnpFlag           = flag.Bool("no-upnp", false, "launch node without automatic UPnP port forwarding")                                                                // Init upnp flag
	rpcPortFlag        = flag.Int("rpc-port", 8080, "launch node with specified RPC port")                                                                                // Init RPC port flag
	forwardRPCFlag     = flag.Bool("forward-rpc", false, "enables forwarding of node RPC terminal ports")                                                                 // Init forward RPC flag
	rpcAddrFlag        = flag.String("rpc-address", fmt.Sprintf("localhost:%s", strconv.Itoa(*rpcPortFlag)), "connects to remote RPC terminal (default: localhost:8080)") // Init remote rpc addr flag
	dataDirFlag        = flag.String("data-dir", common.DataDir, "performs all node i/o operations in given data directory")                                              // Init data dir flag
	nodePortFlag       = flag.Int("node-port", common.DefaultNodePort, "launch node on give port")                                                                        // Init node port flag
	privateNetworkFlag = flag.Bool("private-net", false, "launch node in context of private network")                                                                     // Init private network flag
	archivalNodeFlag   = flag.Bool("archival", false, "launch node in archival mode")                                                                                     // Init archival node flag
)

func main() {
	flag.Parse() // Parse flags

	common.DataDir = *dataDirFlag // Set data-dir

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
	}

	go startNode(*archivalNodeFlag) // Start node

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

	cryptoHandler := cryptoProto.NewCryptoServer(&cryptoServer.Server{}, nil)         // Init handler
	upnpHandler := upnpProto.NewUpnpServer(&upnpServer.Server{}, nil)                 // Init handler
	accountsHandler := accountsProto.NewAccountsServer(&accountsServer.Server{}, nil) // Init handler

	mux := http.NewServeMux() // Init mux

	mux.Handle(cryptoProto.CryptoPathPrefix, cryptoHandler)       // Start mux node handler
	mux.Handle(upnpProto.UpnpPathPrefix, upnpHandler)             // Start mux upnp handler
	mux.Handle(accountsProto.AccountsPathPrefix, accountsHandler) // Start mux accounts handler

	go http.ListenAndServeTLS(":"+strconv.Itoa(*rpcPortFlag), "termCert.pem", "termKey.pem", mux) // Start server
}

// startNode - start necessary services for full node
func startNode(archivalNode bool) {
	if archivalNode { // Check
		err := types.JoinNetwork(common.BootstrapNodes[0], true) // Register node

		if err != nil { // Check for errors
			panic(err) // Panic
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
	- terminal.go accounts service support
	- archival node flag
	- readme
*/
