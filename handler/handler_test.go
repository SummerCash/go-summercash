package handler

import (
	"crypto/tls"
	"net"
	"strconv"
	"testing"

	"github.com/space55/summertech-blockchain/common"
)

// TestStartHandler - test StartHandler() method
func TestStartHandler(t *testing.T) {
	ln, err := tls.Listen("tcp", ":"+strconv.Itoa(7890), common.GeneralTLSConfig) // Listen on port

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	go func(listener *net.Listener) { // Silence timeout
		err = StartHandler(listener) // Start handler

		if err != nil { // Check for errors
			t.Error(err) // Log found error
			t.FailNow()  // Panic
		}
	}(&ln) // Call
}
