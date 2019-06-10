// Package mobile represents the standard main() client.
package mobile

import (
	"context"
	"net/http"
	"time"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	accountsServer "github.com/SummerCash/go-summercash/intrnl/rpc/accounts"
	chainServer "github.com/SummerCash/go-summercash/intrnl/rpc/chain"
	commonServer "github.com/SummerCash/go-summercash/intrnl/rpc/common"
	configServer "github.com/SummerCash/go-summercash/intrnl/rpc/config"
	coordinationChainServer "github.com/SummerCash/go-summercash/intrnl/rpc/coordinationchain"
	cryptoServer "github.com/SummerCash/go-summercash/intrnl/rpc/crypto"
	p2pServer "github.com/SummerCash/go-summercash/intrnl/rpc/p2p"
	accountsProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/accounts"
	chainProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/chain"
	commonProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/common"
	configProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/config"
	coordinationChainProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/coordinationchain"
	cryptoProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/crypto"
	p2pProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/p2p"
	transactionProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/transaction"
	upnpProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/upnp"
	transactionServer "github.com/SummerCash/go-summercash/intrnl/rpc/transaction"
	upnpServer "github.com/SummerCash/go-summercash/intrnl/rpc/upnp"
	"github.com/SummerCash/go-summercash/p2p"
	"github.com/SummerCash/go-summercash/validator"
)

// Run starts a new mobile client, with RPC enabled.
func Run() {
	common.NodePort = 3000 // Set node port

	startRPCServer() // Start RPC server

	startNode(true) // Start node
}

// startRPCServer - start RPC server
func startRPCServer() {
	err := common.GenerateTLSCertificates("term") // Generate certs
	if err != nil {                               // Check for errors
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
	p2pHandler := p2pProto.NewP2PServer(&p2pServer.Server{}, nil)                                                         // Init handler

	mux := http.NewServeMux() // Init mux

	mux.Handle(cryptoProto.CryptoPathPrefix, cryptoHandler)                                  // Start mux node handler
	mux.Handle(upnpProto.UpnpPathPrefix, upnpHandler)                                        // Start mux upnp handler
	mux.Handle(accountsProto.AccountsPathPrefix, accountsHandler)                            // Start mux accounts handler
	mux.Handle(configProto.ConfigPathPrefix, configHandler)                                  // Start mux config handler
	mux.Handle(transactionProto.TransactionPathPrefix, transactionHandler)                   // Start mux config handler
	mux.Handle(chainProto.ChainPathPrefix, chainHandler)                                     // Start mux chain handler
	mux.Handle(coordinationChainProto.CoordinationChainPathPrefix, coordinationChainHandler) // Start mux coordinationChain handler
	mux.Handle(commonProto.CommonPathPrefix, commonHandler)                                  // Start mux common handler
	mux.Handle(p2pProto.P2PPathPrefix, p2pHandler)                                           // Start mux p2p handler

	go http.ListenAndServeTLS(":8080", "termCert.pem", "termKey.pem", mux) // Start server
	go http.ListenAndServe(":8081", mux)                                   // Start server
}

// startNode - start necessary services for full node
func startNode(archivalNode bool) {
	ctx, cancel := context.WithCancel(context.Background()) // Get node context

	defer cancel() // Cancel

	host, err := p2p.NewHost(ctx, 3000, "main_net") // Initialize libp2p host with context and nat manager
	if err != nil {                                 // Check for errors
		panic(err) // Panic
	}

	config, err := config.ReadChainConfigFromMemory() // Read chain config
	if err != nil {                                   // Check for errors
		config, err = p2p.BootstrapConfig(ctx, host, p2p.GetBestBootstrapAddress(ctx, host, "main_net"), "main_net") // Bootstrap config

		if err != nil { // Check for errors
			panic(err) // panic
		}

		err = config.WriteToMemory() // Write config to memory

		if err != nil { // Check for errors
			panic(err) // Panic
		}
	}

	err = config.UpdateChainVersion() // Update chain version

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	validator := validator.Validator(validator.NewStandardValidator(config)) // Initialize validator

	client := p2p.NewClient(host, &validator, "main_net") // Initialize client

	err = client.StartServingStreams() // Start serving

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	if p2p.GetBestBootstrapAddress(ctx, host, "main_net") != "localhost" { // Check can sync
		err = client.SyncNetwork() // Sync network

		if err != nil { // Check for errors
			panic(err) // Panic
		}
	}

	client.StartIntermittentSync(60 * time.Second) // Start intermittent sync
}
