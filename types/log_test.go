package types

import "testing"

// TestNewLog - test functionality of log init
func TestNewLog(t *testing.T) {
	log := NewLog("error", []byte("test"), Error) // Init log

	t.Log(log) // Log log
}

// TestStringLog - test functionality of log to string method
func TestStringLog(t *testing.T) {
	log := NewLog("error", []byte("test"), Error) // Init log

	t.Log(log.String()) // Log log
}

// TestStringLogs - test functionality of log slice to string method
func TestStringLogs(t *testing.T) {
	log := NewLog("error", []byte("test"), Error) // Init log

	logs := []*Log{log, log, log, log} // Init logs

	t.Log(StringLogs(logs)) // Log logs
}

// TestStringLogKeyType - test functionality of log key type to string method
func TestStringLogKeyType(t *testing.T) {
	err := Error // Init custom

	t.Log(err.String()) // Log custom
}
