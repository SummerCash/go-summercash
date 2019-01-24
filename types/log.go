package types

import (
	"encoding/json"
	"strings"
)

const (
	// Return - return log key type
	Return LogKeyType = iota

	// Error - error log key type
	Error

	// Custom - custom log key type
	Custom
)

// LogKeyType - log key type
type LogKeyType int

// Log - log meta container
type Log struct {
	Type  LogKeyType `json:"type"`  // Log type
	Key   string     `json:"key"`   // Log key
	Value []byte     `json:"value"` // Log val
}

/* BEGIN EXPORTED METHODS */

// NewLog - init log with given key value pair
func NewLog(key string, value []byte, logType LogKeyType) *Log {
	return &Log{ // Return log
		Type:  logType, // Set type
		Key:   key,     // Set key
		Value: value,   // Set val
	}
}

// String - get string representation of log
func (log *Log) String() string {
	marshaledVal, _ := json.MarshalIndent(*log, "", "  ") // Marshal

	var marshaledString map[string]interface{} // Init json buffer

	json.Unmarshal(marshaledVal, &marshaledString) // Unmarshal JSON

	marshaledString["value"] = string(log.Value) // Get string representation
	marshaledString["type"] = log.Type.String()  // Get type string representation

	marshaledVal, _ = json.MarshalIndent(marshaledString, "", "  ") // Marshal

	return string(marshaledVal) // Return success
}

// StringLogs - get string representation of log slice
func StringLogs(logs []*Log) string {
	logStrings := []string{} // init log buffer

	for _, log := range logs { // Iterate through logs
		logStrings = append(logStrings, log.String()) // Append string val
	}

	return strings.Join(logStrings, ", ") // Return success
}

// String - get string representation of key type
func (l LogKeyType) String() string {
	return [...]string{"return", "error", "custom"}[l]
}

/* END EXPORTED METHODS */
