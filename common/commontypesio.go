package common

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

/* BEGIN EXPORTED METHODS */

// WriteToMemory - write given address to persistent memory
func (address *Address) WriteToMemory(privateKey *ecdsa.PrivateKey) error {
	err := CreateDirIfDoesNotExit(fmt.Sprintf("%s/keystore", DataDir)) // Create dir if necessary

	if err != nil { // Check for errors
		return err // Return error
	}

	err = WriteGob(fmt.Sprintf("%s/keystore/address_%s.gob", DataDir, address.String()), *address) // Write gob

	if err != nil { // Check for errors
		return err // Return found error
	}

	marshaledPrivateKey, err := x509.MarshalECPrivateKey(privateKey) // Marshal private key

	if err != nil { // Check for errors
		return err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshaledPrivateKey}) // Encode to memory

	err = ioutil.WriteFile(fmt.Sprintf("%s/keystore/key_%s.pem", DataDir, address.String()), pemEncoded, 0644) // Write encoded private key

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// ReadAccountFromMemory - attempt to retrieve given address, keystore from persistent storage
func ReadAccountFromMemory(address *Address) (*Address, *ecdsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/keystore/key_%s.pem", DataDir, address.String())) // Read encoded PEM

	if err != nil { // Check for errors
		return &Address{}, &ecdsa.PrivateKey{}, err // Return found error
	}

	pemEncoded, _ := pem.Decode(data) // Decode data

	privateKey, err := x509.ParseECPrivateKey(pemEncoded.Bytes) // Decrypt pem

	if err != nil { // Check for errors
		return &Address{}, &ecdsa.PrivateKey{}, err // Return found error
	}

	return address, privateKey, nil // No error occurred, return read addr, private key
}

/* END EXPORTED METHODS */
