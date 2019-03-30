package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SummerCash/go-summercash/validator"

	"github.com/SummerCash/go-summercash/p2p"

	"github.com/SummerCash/go-summercash/cli"
	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	accountsServer "github.com/SummerCash/go-summercash/internal/rpc/accounts"
	chainServer "github.com/SummerCash/go-summercash/internal/rpc/chain"
	commonServer "github.com/SummerCash/go-summercash/internal/rpc/common"
	configServer "github.com/SummerCash/go-summercash/internal/rpc/config"
	coordinationChainServer "github.com/SummerCash/go-summercash/internal/rpc/coordinationchain"
	cryptoServer "github.com/SummerCash/go-summercash/internal/rpc/crypto"
	accountsProto "github.com/SummerCash/go-summercash/internal/rpc/proto/accounts"
	chainProto "github.com/SummerCash/go-summercash/internal/rpc/proto/chain"
	commonProto "github.com/SummerCash/go-summercash/internal/rpc/proto/common"
	configProto "github.com/SummerCash/go-summercash/internal/rpc/proto/config"
	coordinationChainProto "github.com/SummerCash/go-summercash/internal/rpc/proto/coordinationchain"
	cryptoProto "github.com/SummerCash/go-summercash/internal/rpc/proto/crypto"
	transactionProto "github.com/SummerCash/go-summercash/internal/rpc/proto/transaction"
	upnpProto "github.com/SummerCash/go-summercash/internal/rpc/proto/upnp"
	transactionServer "github.com/SummerCash/go-summercash/internal/rpc/transaction"
	upnpServer "github.com/SummerCash/go-summercash/internal/rpc/upnp"
	"github.com/SummerCash/go-summercash/upnp"
)

var (
	terminalFlag        = flag.Bool("terminal", false, "launch node in terminal mode")                                                                                     // Init term flag
	upnpFlag            = flag.Bool("no-upnp", false, "launch node without automatic UPnP port forwarding")                                                                // Init upnp flag
	rpcPortFlag         = flag.Int("rpc-port", 8080, "launch node with specified RPC port")                                                                                // Init RPC port flag
	forwardRPCFlag      = flag.Bool("forward-rpc", false, "enables forwarding of node RPC terminal ports")                                                                 // Init forward RPC flag
	rpcAddrFlag         = flag.String("rpc-address", fmt.Sprintf("localhost:%s", strconv.Itoa(*rpcPortFlag)), "connects to remote RPC terminal (default: localhost:8080)") // Init remote rpc addr flag
	dataDirFlag         = flag.String("data-dir", common.DataDir, "performs all node i/o operations in given data directory")                                              // Init data dir flag
	nodePortFlag        = flag.Int("node-port", common.NodePort, "launch node on give port")                                                                               // Init node port flag
	privateNetworkFlag  = flag.Bool("private-net", false, "launch node in context of private network")                                                                     // Init private network flag
	archivalNodeFlag    = flag.Bool("archival", false, "launch node in archival mode")                                                                                     // Init archival node flag
	silent              = flag.Bool("silent", false, "silence all fmt.Print calls")                                                                                        // Init silent flag
	exitOnJoin          = flag.Bool("exit-on-join", false, "exit node on network join")                                                                                    // Init exit on join flag
	version             = flag.Bool("version", false, "get node software version")                                                                                         // Init version flag
	bootstrapNode       = flag.String("bootstrap-node", "", "launch node with provided bootstrap node")                                                                    // Init bootstrap node flag
	bootstrapHost       = flag.Bool("bootstrap", false, "launch node as a genesis boostrap node")                                                                          // Init bootstrap host flag
	disableLogTimeStamp = flag.Bool("silence-timestamps", false, "launch node without terminal timestamp output")                                                          // Init disable log timestamp flag
	networkFlag         = flag.String("network", "main_net", "launch with a given network")                                                                                // Init network flag
)

func main() {
	flag.Parse() // Parse flags

	common.DataDir = *dataDirFlag   // Set data-dir
	common.Silent = *silent         // Set is silent
	common.NodePort = *nodePortFlag // Set node port

	if *disableLogTimeStamp { // Check must disable timestamps
		common.DisableTimestamps = true // Set timestamps disabled
	}

	if *version { // Check needs version
		fmt.Println(config.Version) // Log version

		os.Exit(0) // Stop execution
	}

	if *privateNetworkFlag { // Check private network
		common.ExtIPProviders = []string{} // Set nil providers
	}

	if *bootstrapNode != "" { // Check needs bootstrap node
		common.BootstrapNodes = []string{*bootstrapNode} // Set bootstrap node
	}

	if *bootstrapHost { // Check is bootstrap host
		ipAddr, err := common.GetExtIPAddrWithoutUPnP() // Get IP

		if err != nil { // Check for errors
			panic(err) // Panic
		}

		common.BootstrapNodes = []string{ipAddr + ":" + strconv.Itoa(*nodePortFlag)} // Set bootstrap nodes to local host
	}

	if !*upnpFlag { // Check for UPnP
		if *forwardRPCFlag {
			go upnp.ForwardPortSilent(uint(*rpcPortFlag)) // Forward RPC port
		}

		go upnp.ForwardPortSilent(uint(*nodePortFlag)) // Forward port 3000
	}

	if strings.Contains(*rpcAddrFlag, "localhost") { // Check for default RPC address
		startRPCServer() // Start RPC server

		if !*terminalFlag { // Check only daemon
			startNode(*archivalNodeFlag) // Start node
		} else { // Check with terminal
			go startNode(*archivalNodeFlag) // Start node
		}
	}

	if *terminalFlag { // Check for terminal
		*rpcAddrFlag = strings.Split(*rpcAddrFlag, ":")[0] // Remove port

		cli.NewTerminal(uint(*rpcPortFlag), *rpcAddrFlag) // Initialize terminal
	}
}

// startRPCServer - start RPC server
func startRPCServer() {
	err := common.GenerateTLSCertificates("term") // Generate certs

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	cryptoHandler := cryptoProto.NewCryptoServer(&cryptoServer.Server{}, nil)                                             // Init handler
	upnpHandler := upnpProto.NewUpnpServer(&upnpServer.Server{}, nil)                                                     // Init handler
	accountsHandler := accountsProto.NewAccountsServer(&accountsServer.Server{}, nil)                                     // Init handler
	configHandler := configProto.NewConfigServer(&configServer.Server{}, nil)                                             // Init handler
	transactionHandler := transactionProto.NewTransactionServer(&transactionServer.Server{}, nil)                         // Init handler
	chainHandler := chainProto.NewChainServer(&chainServer.Server{}, nil)                                                 // Init handler
	coordinationChainHandler := coordinationChainProto.NewCoordinationChainServer(&coordinationChainServer.Server{}, nil) // Init handler
	commonHandler := commonProto.NewCommonServer(&commonServer.Server{}, nil)                                             // Init handler

	mux := http.NewServeMux() // Init mux

	mux.Handle(cryptoProto.CryptoPathPrefix, cryptoHandler)                                  // Start mux node handler
	mux.Handle(upnpProto.UpnpPathPrefix, upnpHandler)                                        // Start mux upnp handler
	mux.Handle(accountsProto.AccountsPathPrefix, accountsHandler)                            // Start mux accounts handler
	mux.Handle(configProto.ConfigPathPrefix, configHandler)                                  // Start mux config handler
	mux.Handle(transactionProto.TransactionPathPrefix, transactionHandler)                   // Start mux config handler
	mux.Handle(chainProto.ChainPathPrefix, chainHandler)                                     // Start mux chain handler
	mux.Handle(coordinationChainProto.CoordinationChainPathPrefix, coordinationChainHandler) // Start mux coordinationChain handler
	mux.Handle(commonProto.CommonPathPrefix, commonHandler)                                  // Start mux common handler

	go http.ListenAndServeTLS(":"+strconv.Itoa(*rpcPortFlag), "termCert.pem", "termKey.pem", mux) // Start server
	go http.ListenAndServe(":"+strconv.Itoa(*rpcPortFlag+1), mux)                                 // Start server
}

// startNode - start necessary services for full node
func startNode(archivalNode bool) {
	ctx, cancel := context.WithCancel(context.Background()) // Get node context

	defer cancel() // Cancel

	host, err := p2p.NewHost(ctx, *nodePortFlag) // Initialize libp2p host with context and nat manager

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	config, err := config.ReadChainConfigFromMemory() // Read chain config

	if err != nil { // Check for errors
		config, err = p2p.BootstrapConfig(ctx, host, p2p.GetBestBootstrapAddress(ctx, host), *networkFlag) // Bootstrap config

		if err != nil { // Check for errors
			panic(err) // panic
		}
	}

	validator := validator.Validator(validator.NewStandardValidator(config)) // Initialize validator

	client := p2p.NewClient(host, &validator, *networkFlag) // Initialize client

	err = client.StartServingStreams() // Start serving

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	if p2p.GetBestBootstrapAddress(ctx, host) != "localhost" { // Check can sync
		err = client.SyncNetwork() // Sync network

		if err != nil { // Check for errors
			panic(err) // Panic
		}
	}

	if !*terminalFlag { // Check is not locally running terminal
		client.StartIntermittentSync(120 * time.Second) // Start intermittent sync
	} else { // Check local term
		go client.StartIntermittentSync(120 * time.Second) // Start intermittent sync
	}
}

/*
	TODO:
	- readme
*/
