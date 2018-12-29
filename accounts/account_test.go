package accounts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/config"
	"github.com/space55/summertech-blockchain/types"
)

// TestNewAccount - test functionality of account generation
func TestNewAccount(t *testing.T) {
	address, err := common.StringToAddress("0x040028d536d5351e83fbbec320c194629ace") // Get addr value

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = makeChainConfig(address) // Make config

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := NewAccount() // Generate account

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(account.String()) // Log success
}

// TestAccountFromKey - test functionality of account generation given a privateKey x
func TestAccountFromKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(account.String()) // Log success
}

// TestGetAllAccounts - test functionality of keystore walk method
func TestGetAllAccounts(t *testing.T) {
	account, err := NewAccount() // Generate account

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.WriteToMemory() // Make sure we have at least one account to walk

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	addresses, err := GetAllAccounts() // Walk

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(addresses) // Log success
}

// TestMakeEncodingSafe - test functionality of safe account encoding
func TestMakeEncodingSafe(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(account.String()) // Log success
}

// TestRecoverSafeEncoding - test functionality of safe account encoding recovery
func TestRecoverSafeEncoding(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = account.RecoverSafeEncoding() // Recover from safe encoding

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(account.String()) // Log success
}

// TestString - test functionality of string account serialization
func TestString(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := account.String() // Get string val

	if stringVal == "" { // Check nil
		t.Error("invalid string val") // Log error
		t.FailNow()                   // Panic
	}

	t.Log(stringVal) // Log success
}

// TestBytes - test byte array serialization of account
func TestBytes(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account, err := AccountFromKey(privateKey) // Generate account from key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := account.Bytes() // Get byte val

	if byteVal == nil { // Check is nil
		t.Error("invalid byte val") // Log error
		t.FailNow()                 // Panic
	}

	t.Log(byteVal) // Log success
}

type genesis struct {
	NetworkID uint `json:"networkID"`

	Alloc map[string]map[string]string `json:"alloc"`
}

// makeChainConfig - generate necessary config files
func makeChainConfig(address common.Address) error {
	alloc := make(map[string]map[string]string) // Init map

	alloc[address.String()] = make(map[string]string) // Init map

	alloc[address.String()]["balance"] = "500000000000000" // Set balance

	genesis := genesis{NetworkID: 0, Alloc: alloc} // Init genesis

	json, err := json.MarshalIndent(genesis, "", "  ") // Marshal genesis

	if err != nil { // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile("genesis.json", json, 0644) // Write genesis to JSON

	if err != nil { // Check for errors
		return err // Return error
	}

	config, err := config.NewChainConfig("genesis.json") // Generate config

	if err != nil { // Check for errors
		return err // Return error
	}

	err = config.WriteToMemory() // Write to memory

	if err != nil { // Check for errors
		return err // Return error
	}

	coordinationChain, err := types.NewCoordinationChain() // Init coordinationChain

	if err != nil { // Check for errors
		return err // Return error
	}

	return coordinationChain.WriteToMemory() // Write to memory
}
