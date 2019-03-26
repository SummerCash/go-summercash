// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/polaris-project/go-polaris/common"
)

var (
	// ErrIdentityAlreadyExists is an error definition describing a duplicate peer identity value (already an identity in p2p/).
	ErrIdentityAlreadyExists = errors.New("identity already exists")

	// ErrNoExistingIdentity is an error definition describing the lack of a p2p identity.
	ErrNoExistingIdentity = errors.New("no existing identity")
)

/* BEGIN EXPORTED METHODS */

// GetPeerIdentity gets the peer identity. If no identity exists persistently, it creates one.
func GetPeerIdentity() (*ecdsa.PrivateKey, error) {
	identity, err := GetExistingPeerIdentity() // Get existing peer identity

	if err != nil { // Check for errors
		if err != ErrNoExistingIdentity { // Check for errors
			return &ecdsa.PrivateKey{}, err // Return found error
		}

		identity, err = NewPeerIdentity() // Initialize identity

		if err != nil { // Check for errors
			return &ecdsa.PrivateKey{}, err // Return found error
		}
	}

	return identity, nil // Return identity
}

// GetLibp2pPeerIdentity gets the peer identity as a libp2p PrivKey. If no identity exists persistently, it creates one.
func GetLibp2pPeerIdentity() (*crypto.PrivKey, error) {
	identity, err := GetPeerIdentity() // Get peer identity

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	libp2pPrivateKey, _, err := crypto.ECDSAKeyPairFromKey(identity) // Get libp2p private key

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &libp2pPrivateKey, nil // Return libp2p private key
}

// NewPeerIdentity creates a new p2p identity, and writes it to memory.
func NewPeerIdentity() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Initialize private key

	if err != nil { // Check for errors
		return &ecdsa.PrivateKey{}, err // Return found error
	}

	err = WritePeerIdentity(privateKey) // Write identity

	if err != nil { // Check for errors
		return &ecdsa.PrivateKey{}, err // Return found error
	}

	return privateKey, nil // Return identity
}

// WritePeerIdentity writes a given p2p identity to persistent memory.
func WritePeerIdentity(identity *ecdsa.PrivateKey) error {
	if _, err := os.Stat(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir))); err == nil { // Check existing p2p identity
		return ErrIdentityAlreadyExists // Return error
	}

	x509Encoded, err := x509.MarshalECPrivateKey(identity) // Marshal identity

	if err != nil { // Check for errors
		return err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded}) // Encode to pem

	err = common.CreateDirIfDoesNotExit(common.PeerIdentityDir) // Create identity dir if it doesn't already exist

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir)), pemEncoded, 0644) // Write identity

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// GetExistingPeerIdentity attempts to read an existing p2p identity.
func GetExistingPeerIdentity() (*ecdsa.PrivateKey, error) {
	if _, err := os.Stat(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir))); err != nil {
		return &ecdsa.PrivateKey{}, ErrNoExistingIdentity // Return error
	}

	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir))) // Read identity

	if err != nil { // Check for errors
		return &ecdsa.PrivateKey{}, err // Return found error
	}

	block, _ := pem.Decode(data) // Decode pem

	peerIdentity, err := x509.ParseECPrivateKey(block.Bytes) // Parse private key pem block

	if err != nil { // Check for errors
		return &ecdsa.PrivateKey{}, err // Return found error
	}

	return peerIdentity, nil // Return read identity
}

/* END EXPORTED METHODS */
