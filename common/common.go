package common

import (
	"errors"
	"strings"
)

/* BEGIN EXPORTED METHODS */

/*
	BEGIN TERMINAL METHODS
*/

// ParseStringMethodCall - attempt to parse string as method call, returning receiver, method name, and params
func ParseStringMethodCall(input string) (string, string, []string, error) {
	if input == "" { // Check for errors
		return "", "", []string{}, errors.New("nil input") // Return found error
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
		return []string{}, errors.New("nil input") // Return found error
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

/* END EXPORTED METHODS */
