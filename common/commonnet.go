package common

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	upnp "github.com/NebulousLabs/go-upnp"
)

var (
	// ExtIPProviders - preset macro defining list of available external IP checking services
	ExtIPProviders = []string{"http://checkip.amazonaws.com/", "http://icanhazip.com/", "http://www.trackip.net/ip", "http://bot.whatismyipaddress.com/", "https://ipecho.net/plain", "http://myexternalip.com/raw"}
)

/* BEGIN EXPORTED METHODS */

// SendBytes - attempt to send specified bytes to given address
func SendBytes(b []byte, address string) error {
	if strings.Count(address, ":") > 1 { // Check IPv6
		address = "[" + address[:strings.LastIndex(address, ":")] + "]" + address[strings.LastIndex(address, ":"):] // Set address
	}

	connection, err := tls.Dial("tcp", address, GeneralTLSConfig) // Connect to given address

	if err != nil { // Check for errors
		return err // Return found error
	}

	_, err = connection.Write(b) // Write data to connection

	if err != nil { // Check for errors
		return err // Return found errors
	}

	err = connection.Close() // Close connection

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// ReadConnectionWaitAsyncNoTLS - attempt to read from connection in an asynchronous fashion, after waiting for peer to write
func ReadConnectionWaitAsyncNoTLS(conn net.Conn) ([]byte, error) {
	data := make(chan []byte) // Init buffer
	err := make(chan error)   // Init error buffer

	conn.SetReadDeadline(time.Now().Add(2 * time.Second)) // Set read deadline

	go func(data chan []byte, err chan error) {
		reads := 0 // Init reads buffer

		for {
			reads++ // Increment read

			var buffer bytes.Buffer // Init buffer

			readData, readErr := io.Copy(&buffer, conn) // Read connection

			if readErr != nil && readErr != io.EOF && reads > 3 { // Check for errors
				err <- readErr // Write read error
			} else if readData == 0 { // Check for nil readData
				continue // Continue
			}

			data <- buffer.Bytes() // Write read data
		}
	}(data, err)

	ticker := time.Tick(3 * time.Second) // Init ticker

	for { // Continuously read from connection
		select {
		case readData := <-data: // Read data from connection
			return readData, nil // Return read data
		case readErr := <-err: // Error on read
			return []byte{}, readErr // Return error
		case <-ticker: // Timed out
			return []byte{}, errors.New("timed out") // Return timed out error
		}
	}
}

/*
	BEGIN IP ADDR METHODS
*/

// GetExtIPAddrWithUPnP - retrieve the external IP address of the current machine via upnp
func GetExtIPAddrWithUPnP() (string, error) {
	// connect to router
	d, err := upnp.Discover()
	if err != nil { // Check for errors
		return "", err // return error
	}

	// discover external IP
	ip, err := d.ExternalIP()
	if err != nil { // Check for errors
		return "", err // return error
	}
	return ip, nil
}

// GetExtIPAddrWithoutUPnP - retrieve the external IP address of the current machine w/o upnp
func GetExtIPAddrWithoutUPnP() (string, error) {
	addresses := []string{} // Init address buffer

	finished := make(chan bool) // Init finished

	for _, provider := range ExtIPProviders { // Iterate through providers
		go getIPFromProviderAsync(provider, &addresses, finished) // Fetch IP
	}

	<-finished // Wait until finished

	close(finished) // Close channel

	return getNonNilInStringSlice(addresses) // Return valid address
}

/*
	END IP ADDR METHODS
*/

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// getIPFromProvider - get IP address from given IP provider
func getIPFromProvider(provider string) (string, error) {
	resp, err := http.Get(provider) // Attempt to check IP via provider

	if err != nil { // Check for errors
		return "", err // Return error
	}

	defer resp.Body.Close() // Close connection

	ip, err := ioutil.ReadAll(resp.Body) // Read address

	if err != nil { // Check for errors
		return "", err // Return error
	}

	stringVal := string(ip[:]) // Fetch string value

	return strings.TrimSpace(stringVal), nil // Return ip
}

// getIPFromProviderAsync - asynchronously get IP address from given IP provider
func getIPFromProviderAsync(provider string, buffer *[]string, finished chan bool) {
	if len(*buffer) == 0 { // Check IP not already determined
		resp, err := http.Get(provider) // Attempt to check IP via provider

		if err != nil { // Check for errors
			if len(*buffer) == 0 { // Double check IP not already determined
				*buffer = append(*buffer, "") // Set IP
				finished <- true              // Set finished
			}
		} else {
			defer resp.Body.Close() // Close connection

			ip, _ := ioutil.ReadAll(resp.Body) // Read address

			stringVal := string(ip[:]) // Fetch string value

			if len(*buffer) == 0 { // Double check IP not already determined
				*buffer = append(*buffer, strings.TrimSpace(stringVal)) // Set ip

				finished <- true // Set finished
			}
		}
	}
}

/* END INTERNAL METHODS */
