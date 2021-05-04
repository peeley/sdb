package db

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// DBState is used to track which database the user is currently in, along with
// any data associated with transactions.
type DBState struct {
	CurrentDB   string
	Transaction *Transaction
}

// All SQL statement types implement this interface. The `Execute` function
// contains the core logic of the query, which is executed in the REPL at the
// `sdb/main.go` main function. See each statement's respective file in
// `sdb/statements/<statement>.go`.
type Executable interface {
	Execute(*DBState) error
}

// Transaction structs are responsible for tracking the lock files created
// during transactions, along with storing statements that atomically executed
// when the transaction is committed.
type Transaction struct {
	LockFiles  []string
	Statements []Executable
}

func (state *DBState) BeginTransaction() {
	state.Transaction = &Transaction{}
}

func (state *DBState) IsTransacting() bool {
	return state.Transaction != nil
}

// Checks if the given table is currently locked by any process.
func (state *DBState) TableLockExists(tableName string) bool {
	lockFileName := state.CurrentDB + "/." + tableName + "_lock"
	_, err := os.Stat(lockFileName)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// Creates lock file to signify table is undergoing transaction. Returns tuple
// of (string, error) signifying the name of the lock file created or any errors
// during creation.
func (state *DBState) createTableLock(tableName string) (string, error) {
	lockFileName := state.CurrentDB + "/." + tableName + "_lock"
	_, err := os.Create(lockFileName)
	if err != nil {
		return "", fmt.Errorf("!Failed to create transaction lock: %v", err)
	}

	return lockFileName, nil
}

// Lets state "acquire" a lock on a table file. Checks if a lock exists - if
// not, it creates the lock file. If a lock file does exist, it checks if the
// PID in the lock file matches the current process' PID. An existing lock file
// with a different PID means another process is currently locking that table.
func (state *DBState) AcquireTableLock(tableName string) (string, error) {
	lockFileName := state.CurrentDB + "/." + tableName + "_lock"
	if state.TableLockExists(tableName) {
		file, err := os.OpenFile(lockFileName, os.O_RDWR, 0777)
		if err != nil {
			return "", err
		}

		lockFileReader := bufio.NewReader(file)
		contents, err := lockFileReader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}

		if strings.Contains(contents, fmt.Sprintf("%v", os.Getpid())) {
			return lockFileName, nil
		} else {
			return "", fmt.Errorf("!Table %v is locked.\n", tableName)
		}
	} else {
		lockFileName, _ := state.createTableLock(tableName)
		file, _ := os.OpenFile(lockFileName, os.O_RDWR|os.O_CREATE, 0777)
		_, err := file.WriteString(fmt.Sprintf("%v", os.Getpid()))
		if err != nil {
			return "", err
		}

		return lockFileName, nil
	}
}

// Create a new DBState with no current database. Will not be valid - user must
// execute a `USE` query before actually executing any queries other than
// `CREATE DATABASE`.
func NewState() DBState {
	return DBState{
		CurrentDB: "",
	}
}
