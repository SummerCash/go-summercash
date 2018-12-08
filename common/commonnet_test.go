package common

import (
	"strings"
	"testing"
)

/* BEGIN EXTERNAL METHODS */

// TestSendBytes - test functionality of SendBytes() method
func TestSendBytes(t *testing.T) {
	err := SendBytes([]byte("test"), "1.1.1.1:443") // Write to address

	if err != nil { // Check for errors
		t.Errorf(err.Error()) // Log found error
		t.FailNow()           // Panic
	}

	t.Logf("wrote to address 1.1.1.1") // Log success
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
	buffer := &[]string{}       // Init buffer
	finished := make(chan bool) // Init finished

	go getIPFromProviderAsync("http://checkip.amazonaws.com/", buffer, finished) // Asynchronously get IP from provider

	<-finished // Wait for finished

	t.Log((*buffer)[0]) // Log success
}

/*
	END IP ADDR METHODS
*/

/* END EXTERNAL METHODS */
