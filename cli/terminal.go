package cli

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mitsukomegumi/GoP2P/common"
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
