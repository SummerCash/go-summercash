package common

import (
	"crypto/tls"
	"strings"
	"testing"
	"time"
)

/* BEGIN EXTERNAL METHODS */

// TestSendBytes - test functionality of SendBytes() method
func TestSendBytes(t *testing.T) {
	err := SendBytes([]byte("test"), "1.1.1.1:443") // Write to address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("wrote to address 1.1.1.1") // Log success
}

// TestSendBytesResult - test functionality of SendBytesResult() method
func TestSendBytesResult(t *testing.T) {
	var chainBytes []byte // Init buffer
	var err error         // Init error buffer

	for x := 0; x != 5; x++ { // Read 10 times
		chainBytes, err = SendBytesResult([]byte("cChainRequest"), BootstrapNodes[0]) // Get coordination chain

		if err != nil { // Check for errors
			t.Error(err) // Log found error
			t.FailNow()  // Panic
		}
	}

	t.Log(chainBytes) // Log success
}

// TestReadConnectionWaitAsyncNoTLS - test functionality of ReadConnectionWaitAsyncNoTLS() method
func TestReadConnectionWaitAsyncNoTLS(t *testing.T) {
	connection, err := tls.Dial("tcp", "1.1.1.1:443", GeneralTLSConfig) // Connect to given address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	_, err = connection.Write(append([]byte("test"), byte('\''))) // Write test data to connection

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	data, err := ReadConnectionWaitAsyncNoTLS(connection) // Read connection

	if err != nil && !strings.Contains(err.Error(), "EOF") { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(data) // Log success
}

/*
	BEGIN IP ADDR METHODS
*/

// TestGetExtIPAddrWitUPnP - test functionality of GetExtIPAddrWithUPnP() method
func TestGetExtIPAddrWitUPnP(t *testing.T) {
	ip, err := GetExtIPAddrWithUPnP() // Get IP

	if err != nil && !strings.Contains(err.Error(), "no UPnP") { // Check for non-no-upnp errors
		t.Log(ip)    // Log IP for cov
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(ip) // Log success
}

// TestGetExtIPAddrWithoutUPnP - test functionality of GetExtIPAddrWithoutUPnP() method
func TestGetExtIPAddrWithoutUPnP(t *testing.T) {
	ip, err := GetExtIPAddrWithoutUPnP() // Get IP

	if err != nil { // Check for errors
		t.Log(ip)    // Log IP for cov
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(ip) // Log success
}

// TestConnected - test is connected to internet
func TestConnected(t *testing.T) {
	result := Connected() // Get is connected

	t.Log(result) // Log success
}

// TestGetIPFromProvider - test functionality of getIPFromProvider() method
func TestGetIPFromProvider(t *testing.T) {
	ip, err := getIPFromProvider("http://checkip.amazonaws.com/") // Get IP from provider

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(ip) // Log success
}

// TestGetIPFromProviderAsync - test functionality of getIPFromProviderAsync() method
func TestGetIPFromProviderAsync(t *testing.T) {
	buffer := &[]string{} // Init buffer
	finished := false     // Init finished

	startTime := time.Now() // Get start time

	go getIPFromProviderAsync("http://checkip.amazonaws.com/", buffer, &finished) // Asynchronously get IP from provider

	for finished != true { // Wait until finished
		if time.Now().Sub(startTime) > 10*time.Second { // Check timeout
			t.Fatal("timed out...") // Panic
		}
	}

	t.Log((*buffer)[0]) // Log success
}

/*
	END IP ADDR METHODS
*/

/* END EXTERNAL METHODS */
