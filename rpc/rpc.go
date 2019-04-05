package rpc

import (
	"net/http"
	"strconv"

	"github.com/SummerCash/go-summercash/common"

	accountsServer "github.com/SummerCash/go-summercash/intrnl/rpc/accounts"
	chainServer "github.com/SummerCash/go-summercash/intrnl/rpc/chain"
	commonServer "github.com/SummerCash/go-summercash/intrnl/rpc/common"
	configServer "github.com/SummerCash/go-summercash/intrnl/rpc/config"
	coordinationChainServer "github.com/SummerCash/go-summercash/intrnl/rpc/coordinationchain"
	cryptoServer "github.com/SummerCash/go-summercash/intrnl/rpc/crypto"
	accountsProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/accounts"
	chainProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/chain"
	commonProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/common"
	configProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/config"
	coordinationChainProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/coordinationchain"
	cryptoProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/crypto"
	transactionProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/transaction"
	upnpProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/upnp"
	transactionServer "github.com/SummerCash/go-summercash/intrnl/rpc/transaction"
	upnpServer "github.com/SummerCash/go-summercash/intrnl/rpc/upnp"
)

/* BEGIN EXPORTED METHODS */

// StartRPCServer - start RPC server
func StartRPCServer(port int) error {
	err := common.GenerateTLSCertificates("term") // Generate certs

	if err != nil { // Check for errors
		return err // Return found error
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

	go http.ListenAndServeTLS(":"+strconv.Itoa(port), "termCert.pem", "termKey.pem", mux) // Start server
	go http.ListenAndServe(":"+strconv.Itoa(port+1), mux)                                 // Start server

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */
