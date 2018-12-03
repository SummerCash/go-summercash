package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/space55/summertech-blockchain/cli"
	"github.com/space55/summertech-blockchain/common"
	cryptoServer "github.com/space55/summertech-blockchain/internal/rpc/crypto"
	cryptoProto "github.com/space55/summertech-blockchain/internal/rpc/proto/crypto"
	"github.com/space55/summertech-blockchain/upnp"
)

var (
	terminalFlag   = flag.Bool("terminal", false, "launch node in terminal mode")                                                                                     // Init term flag
	upnpFlag       = flag.Bool("no-upnp", false, "launch node without automatic UPnP port forwarding")                                                                // Init upnp flag
	rpcPortFlag    = flag.Int("rpc-port", 8080, "launch node with specified RPC port")                                                                                // Init RPC port flag
	forwardRPCFlag = flag.Bool("forward-rpc", false, "enables forwarding of node RPC terminal ports")                                                                 // Init forward RPC flag
	rpcAddrFlag    = flag.String("rpc-address", fmt.Sprintf("localhost:%s", strconv.Itoa(*rpcPortFlag)), "connects to remote RPC terminal (default: localhost:8080)") // Init remote rpc addr flag
)

func main() {
	flag.Parse() // Parse flags

	if !*upnpFlag { // Check for UPnP
		if *forwardRPCFlag {
			go upnp.ForwardPortSilent(uint(*rpcPortFlag)) // Forward RPC port
		}

		go upnp.ForwardPortSilent(3000) // Forward port 3000
	}

	if strings.Contains(*rpcAddrFlag, "localhost") { // Check for default RPC address
		startRPCServer() // Start RPC server
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

	cryptoHandler := cryptoProto.NewCryptoServer(&cryptoServer.Server{}, nil) // Init handler

	mux := http.NewServeMux() // Init mux

	mux.Handle(cryptoProto.CryptoPathPrefix, cryptoHandler) // Start mux node handler

	go http.ListenAndServeTLS(":"+strconv.Itoa(*rpcPortFlag), "termCert.pem", "termKey.pem", mux) // Start server
}
