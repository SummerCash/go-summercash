package common

import "testing"

/* BEGIN EXPORTED METHODS */

/*
	BEGIN TERMINAL METHODS
*/

// TestParseStringMethodCall - test functionality of ParseStringMethodCall() function
func TestParseStringMethodCall(t *testing.T) {
	input := "common.Sha3(test)" // Init input

	receiver, methodName, params, err := ParseStringMethodCall(input) // Parse string method call

	if err != nil { // Check for errors
		t.Errorf(err.Error()) // Log error
		t.FailNow()           // Panic
	}

	t.Logf("found parsed method call %s, %s, %s", receiver, methodName, params[0]) // Log success
}

// TestParseStringParams - test functionality of ParseStringParams() function
func TestParseStringParams(t *testing.T) {
	input := "common.Sha3(test)" // Init input

	params, err := ParseStringParams(input) // Parse string params

	if err != nil { // Check for errors
		t.Errorf(err.Error()) // Log error
		t.FailNow()           // Panic
	}

	t.Logf("found parsed params %s", params[0]) // Log success
}

// TestStringStripReceiverCall - test functionality of StripReceiverCall() function
func TestStringStripReceiverCall(t *testing.T) {
	input := "common.Sha3(test)" // Init input

	stripped := StringStripReceiverCall(input) // Parse string params

	t.Logf("found stripped %s", stripped) // Log success
}

// TestStringStripParentheses - test functionality of StringStripParentheses() function
func TestStringStripParentheses(t *testing.T) {
	input := "common.Sha3(test)" // Init input

	stripped := StringStripParentheses(input) // Strip parentheses

	t.Logf("found value %s", stripped) // Log success
}

// TestStringFetchCallReceiver - test functionality of StringFetchCallReceiver() method
func TestStringFetchCallReceiver(t *testing.T) {
	input := "common.Sha3(test)" // Init input

	receiver := StringFetchCallReceiver(input) // Fetch receiver

	t.Logf("found receiver %s", receiver) // Log success
}

/*
	END TERMINAL METHODS
*/

/* END EXPORTED METHODS */
