package common

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	upnp "github.com/NebulousLabs/go-upnp"
)

var (
	// ExtIPProviders - preset macro defining list of available external IP checking services
	ExtIPProviders = []string{
		"http://checkip.amazonaws.com/",
		"http://icanhazip.com/",
		"http://www.trackip.net/ip",
		"http://bot.whatismyipaddress.com/",
		"https://ipecho.net/plain",
		"http://myexternalip.com/raw",
	}
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

// SendBytesResult - gop2pCommon.SendBytesResult performance-optimized wrapper
func SendBytesResult(b []byte, address string) ([]byte, error) {
	if strings.Count(address, ":") > 1 { // Check IPv6
		address = "[" + address[:strings.LastIndex(address, ":")] + "]" + address[strings.LastIndex(address, ":"):] // Set address
	}

	connection, err := tls.Dial("tcp", address, GeneralTLSConfig) // Connect to given address

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	_, err = connection.Write(b) // Write data to connection

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	requestStartTime := time.Now() // Get start time

	var response []byte // Init buffer

	for response == nil || len(response) == 0 { // Keep reading until buffer isn't nil
		response, err = ioutil.ReadAll(connection) // Read connection

		if err != nil && time.Now().Sub(requestStartTime) > 10*time.Second { // Check for errors after timeout
			return nil, err // Return found error
		}
	}

	return response, nil // Return response
}

// ReadConnectionWaitAsyncNoTLS - attempt to read from connection in an asynchronous fashion, after waiting for peer to write
func ReadConnectionWaitAsyncNoTLS(conn net.Conn) ([]byte, error) {
	requestStartTime := time.Now() // Get start time

	var response []byte // Init buffer
	var err error       // Init error buffer

	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second)) // Set read deadline

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	for response == nil || len(response) == 0 { // Keep reading until buffer isn't nil
		response, err = ioutil.ReadAll(conn) // Read connection

		if err != nil && time.Now().Sub(requestStartTime) > 10*time.Second { // Check for errors after timeout
			return nil, err // Return found error
		}
	}

	return response, nil // Return response
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
	if len(ExtIPProviders) == 0 { // Check no providers
		host, err := os.Hostname() // Get host

		if err != nil { // Check for errors
			return "localhost", err // Return found error
		}

		addrs, err := net.LookupIP(host) // Get IP

		if err != nil { // Check for errors
			return "localhost", err // Return found error
		}

		for _, addr := range addrs { // Iterate through addresses
			if ipv4 := addr.To4(); ipv4 != nil { // Get IPv4 representation
				return ipv4.String(), nil // Return IP
			}

			return addr.String(), nil // Return raw IP
		}

		return "localhost", nil // Return localhost
	}

	addresses := []string{} // Init address buffer

	finished := make(chan bool) // Init finished

	for _, provider := range ExtIPProviders { // Iterate through providers
		go getIPFromProviderAsync(provider, &addresses, finished) // Fetch IP
	}

	<-finished // Wait until finished

	close(finished) // Close channel

	return getNonNilInStringSlice(addresses) // Return valid address
}

// Connected - check if connected to internet
func Connected() (connected bool) {
	resp, err := http.Get("http://clients3.google.com/generate_204") // Perform request

	if err != nil { // Check for errors
		return false // Return result
	}

	_, err = ioutil.ReadAll(resp.Body) // Read response

	if err != nil { // Check for errors
		return false // Return result
	}

	resp.Body.Close() // Close

	return true // Return connected
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
	if Connected() { // Check connected to internet
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
	} else {
		if len(*buffer) == 0 {
			*buffer = append(*buffer, "localhost") // Set localhost

			finished <- true // Set finished
		}
	}
}

/* END INTERNAL METHODS */
