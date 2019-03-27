// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/polaris-project/go-polaris/common"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

var (
	// ErrIdentityAlreadyExists is an error definition describing a duplicate peer identity value (already an identity in p2p/).
	ErrIdentityAlreadyExists = errors.New("identity already exists")

	// ErrNoExistingIdentity is an error definition describing the lack of a p2p identity.
	ErrNoExistingIdentity = errors.New("no existing identity")
)

/* BEGIN EXPORTED METHODS */

// GetPeerIdentity gets the peer identity. If no identity exists persistently, it creates one.
func GetPeerIdentity() (*crypto.PrivKey, error) {
	identity, err := GetExistingPeerIdentity() // Get existing peer identity

	if err != nil { // Check for errors
		if err != ErrNoExistingIdentity { // Check for errors
			return nil, err // Return found error
		}

		identity, err = NewPeerIdentity() // Initialize identity

		if err != nil { // Check for errors
			return nil, err // Return found error
		}
	}

	return identity, nil // Return identity
}

// NewPeerIdentity creates a new p2p identity, and writes it to memory.
func NewPeerIdentity() (*crypto.PrivKey, error) {
	privateKey, _, err := crypto.GenerateRSAKeyPair(2048, rand.Reader) // Generate RSA key pair

	err = WritePeerIdentity(&privateKey) // Write identity

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &privateKey, nil // Return identity
}

// WritePeerIdentity writes a given p2p identity to persistent memory.
func WritePeerIdentity(identity *crypto.PrivKey) error {
	if _, err := os.Stat(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir))); err == nil { // Check existing p2p identity
		return ErrIdentityAlreadyExists // Return error
	}

	encoded, err := crypto.MarshalPrivateKey(*identity) // Marshal identity

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = common.CreateDirIfDoesNotExit(common.PeerIdentityDir) // Create identity dir if it doesn't already exist

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir)), encoded, 0644) // Write identity

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// GetExistingPeerIdentity attempts to read an existing p2p identity.
func GetExistingPeerIdentity() (*crypto.PrivKey, error) {
	if _, err := os.Stat(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir))); err != nil {
		return nil, ErrNoExistingIdentity // Return error
	}

	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir))) // Read identity

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	peerIdentity, err := crypto.UnmarshalPrivateKey(data) // Unmarshal data

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &peerIdentity, nil // Return read identity
}

/* END EXPORTED METHODS */
