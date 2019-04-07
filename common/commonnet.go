package common

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
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

	writer := bufio.NewWriter(connection) // Initialize writer

	_, err = writer.Write(append(b, byte('\a'))) // Write data

	if err != nil { // Check for errors
		return err // Return found errors
	}

	writer.Flush() // Flush writer

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

	d := net.Dialer{Timeout: 5 * time.Second} // Init dialer with timeout

	connection, err := tls.DialWithDialer(&d, "tcp", address, GeneralTLSConfig) // Connect to given address

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	readWriter := bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection)) // Initialize read writer

	_, err = readWriter.Write(append(b, byte('\a'))) // Write data to connection

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	readWriter.Flush() // Flush

	response, err := readWriter.ReadBytes('\a') // Read up to delimiter

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return bytes.Trim(response, "\a"), nil // Return response
}

// ReadConnectionWaitAsyncNoTLS - attempt to read from connection in an asynchronous fashion, after waiting for peer to write
func ReadConnectionWaitAsyncNoTLS(conn net.Conn) ([]byte, error) {
	reader := bufio.NewReader(conn) // Initialize reader

	readBytes, err := reader.ReadBytes('\a') // Read up to delimiter

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return bytes.Trim(readBytes, "\a"), nil // Return read bytes w/trimmed delimiter
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

	startTime := time.Now() // Get start time

	addresses := []string{} // Init address buffer

	finished := false // Init finished

	for _, provider := range ExtIPProviders { // Iterate through providers
		go getIPFromProviderAsync(provider, &addresses, &finished) // Fetch IP
	}

	for finished != true { // Wait until finished
		if time.Now().Sub(startTime) > 10*time.Second { // Check timeout
			return "localhost", errors.New("timed out") // Return timed out error
		}
	}

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
func getIPFromProviderAsync(provider string, buffer *[]string, finished *bool) {
	if Connected() { // Check connected to internet
		if len(*buffer) == 0 { // Check IP not already determined
			resp, err := http.Get(provider) // Attempt to check IP via provider

			if err != nil { // Check for errors
				if len(*buffer) == 0 { // Double check IP not already determined
					*buffer = append(*buffer, "") // Set IP
					*finished = true              // Set finished
				}
			} else {
				defer resp.Body.Close() // Close connection

				ip, _ := ioutil.ReadAll(resp.Body) // Read address

				stringVal := string(ip[:]) // Fetch string value

				if len(*buffer) == 0 { // Double check IP not already determined
					*buffer = append(*buffer, strings.TrimSpace(stringVal)) // Set ip

					*finished = true // Set finished
				}
			}
		}
	} else {
		if len(*buffer) == 0 {
			*buffer = append(*buffer, "localhost") // Set localhost

			*finished = true // Set finished
		}
	}
}

/* END INTERNAL METHODS */
