package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// ErrNilInput - error definition describing input of 0 char length
	ErrNilInput = errors.New("nil input")

	// ErrVerboseNotAllowed - error definition describing config preventing print call
	ErrVerboseNotAllowed = errors.New("verbose output not allowed")

	// DataDir - global data directory definition
	DataDir = getDataDir()

	// ConfigDir - global config directory definition
	ConfigDir = filepath.FromSlash(fmt.Sprintf("%s/config", DataDir))

	// DisableTimestamps - global declaration for timestamp format config
	DisableTimestamps = false

	// GeneralTLSConfig - general global GoP2P TLS Config
	GeneralTLSConfig = &tls.Config{ // Init TLS config
		Certificates:       []tls.Certificate{getTLSCerts("general")},
		InsecureSkipVerify: true,
		ServerName:         "localhost",
	}

	// BootstrapNodes - global definition for set of bootstrap nodes
	BootstrapNodes = []string{
		"108.41.124.60:3000", // Boot node 0
	}

	// Silent - silent config
	Silent = false

	// NodePort - global port definition
	NodePort = 3000
)

/* BEGIN EXPORTED METHODS */

/*
	BEGIN TERMINAL METHODS
*/

// Log - fmt.Println wrapper
func Log(a ...interface{}) (int, error) {
	if !Silent { // Check verbose allowed
		return fmt.Println(a...) // Log
	}

	return 0, ErrVerboseNotAllowed // Return error
}

// Logf - fmt.Printf wrapper
func Logf(format string, a ...interface{}) (int, error) {
	if !Silent { // Check verbose allowed
		if !DisableTimestamps { // Check timestamps not disabled
			if strings.Contains(format, "==") { // Check room for formatting
				format = fmt.Sprintf("[%s] ", time.Now().UTC().Format("Jan 2 03:04:05PM 2006")) + format // Append current time
			}
		}

		return fmt.Printf(format, a...) // Print
	}

	return 0, ErrVerboseNotAllowed // Return error
}

// ParseStringMethodCall - attempt to parse string as method call, returning receiver, method name, and params
func ParseStringMethodCall(input string) (string, string, []string, error) {
	if input == "" { // Check for errors
		return "", "", []string{}, ErrNilInput // Return found error
	} else if !strings.Contains(input, "(") || !strings.Contains(input, ")") {
		input = input + "()" // Fetch receiver methods
	}

	if !strings.Contains(input, ".") { // Check for nil receiver
		return "", "", []string{}, errors.New("invalid method " + input) // Return found error
	}

	method := strings.Split(strings.Split(input, "(")[0], ".")[1] // Fetch method

	receiver := StringFetchCallReceiver(input) // Fetch receiver

	params := []string{} // Init buffer

	if !strings.Contains(input, "()") { // Check for nil params
		params, _ = ParseStringParams(input) // Fetch params
	}

	return receiver, method, params, nil // No error occurred, return parsed method+params
}

// ParseStringParams - attempt to fetch string parameters from (..., ..., ...) style call
func ParseStringParams(input string) ([]string, error) {
	if input == "" { // Check for errors
		return []string{}, ErrNilInput // Return found error
	}

	parenthesesStripped := StringStripParentheses(StringStripReceiverCall(input)) // Strip parentheses

	params := strings.Split(parenthesesStripped, ", ") // Split by ', '

	return params, nil // No error occurred, return split params
}

// StringStripReceiverCall - strip receiver from string method call
func StringStripReceiverCall(input string) string {
	return "(" + strings.Split(input, "(")[1] // Split
}

// StringStripParentheses - strip parantheses from string
func StringStripParentheses(input string) string {
	leftStripped := strings.Replace(input, "(", "", -1) // Strip left parent

	return strings.Replace(leftStripped, ")", "", -1) // Return right stripped
}

// StringFetchCallReceiver - attempt to fetch receiver from string, as if it were a x.y(..., ..., ...) style method call
func StringFetchCallReceiver(input string) string {
	return strings.Split(strings.Split(input, "(")[0], ".")[0] // Return split string
}

/*
	END TERMINAL METHODS
*/

/*
	BEGIN TLS METHODS
*/

// GenerateTLSCertificates - generate necessary TLS certificates, keys
func GenerateTLSCertificates(namePrefix string) error {
	_, certErr := os.Stat(fmt.Sprintf("%sCert.pem", namePrefix)) // Check for error reading file
	_, keyErr := os.Stat(fmt.Sprintf("%sKey.pem", namePrefix))   // Check for error reading file

	if os.IsNotExist(certErr) || os.IsNotExist(keyErr) { // Check for does not exist error
		privateKey, err := generateTLSKey(namePrefix) // Generate key

		if err != nil { // Check for errors
			return err // Return found error
		}

		err = generateTLSCert(privateKey, namePrefix) // Generate cert

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	return nil // No error occurred, return nil
}

/*
	END TLS METHODS
*/

/*
	BEGIN MAIN METHODS
*/

// Forever - prevent thread from closing
func Forever() {
	for {
		time.Sleep(time.Second)
	}
}

/*
	END MAIN METHODS
*/

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

/*
	BEGIN TLS METHODS
*/

// generateTLSKey - generates necessary TLS key
func generateTLSKey(namePrefix string) (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	marshaledPrivateKey, err := x509.MarshalECPrivateKey(privateKey) // Marshal private key

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshaledPrivateKey}) // Encode to memory

	err = ioutil.WriteFile(fmt.Sprintf("%sKey.pem", namePrefix), pemEncoded, 0644) // Write pem

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return privateKey, nil // No error occurred, return nil
}

// generateTLSCert - generates necessary TLS cert
func generateTLSCert(privateKey *ecdsa.PrivateKey, namePrefix string) error {
	notBefore := time.Now() // Fetch current time

	notAfter := notBefore.Add(292 * (365 * (24 * time.Hour))) // Fetch 'deadline'

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)     // Init limit
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit) // Init serial number

	if err != nil { // Check for errors
		return err // Return found error
	}

	template := x509.Certificate{ // Init template
		SerialNumber: serialNumber, // Generate w/serial number
		Subject: pkix.Name{ // Generate w/subject
			Organization: []string{"localhost"}, // Generate w/org
		},
		NotBefore: notBefore, // Generate w/not before
		NotAfter:  notAfter,  // Generate w/not after

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, // Generate w/key usage
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},               // Generate w/ext key
		BasicConstraintsValid: true,                                                         // Generate w/basic constraints
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privateKey), privateKey) // Generate certificate

	if err != nil { // Check for errors
		return err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert}) // Encode pem

	err = ioutil.WriteFile(fmt.Sprintf("%sCert.pem", namePrefix), pemEncoded, 0644) // Write cert file

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// getTLSCert - attempt to read TLS cert from current dir
func getTLSCerts(certPrefix string) tls.Certificate {
	GenerateTLSCertificates(certPrefix) // Generate certs

	cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%sCert.pem", certPrefix), fmt.Sprintf("%sKey.pem", certPrefix)) // Load key pair

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	return cert // Return read certificates
}

// publicKey - cast to public key
func publicKey(privateKey interface{}) interface{} {
	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

/*
	END TLS METHODS
*/

/*
	BEGIN MISC METHODS
*/

// getNonNilInStringSlice - get non nil string in slice
func getNonNilInStringSlice(slice []string) (string, error) {
	for _, entry := range slice { // Iterate through entries
		if entry != "" { // Check for non-nil entry
			return entry, nil // Return valid entry
		}
	}

	return "", fmt.Errorf("couldn't find non-nil element in slice %v", slice) // Couldn't find valid address, return error
}

// getDataDir - get absolute data dir
func getDataDir() string {
	abs, _ := filepath.Abs("./data") // Get absolute dir

	return filepath.FromSlash(abs) // Match slashes
}

/*
	END MISC METHODS
*/

/* END INTERNAL METHODS */
