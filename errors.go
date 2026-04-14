package sqlite

import (
	"errors"
	"fmt"
)

import "C"

// what am I even doing
var (
	ErrGeneric       = errors.New("SQlite error")
	ErrInternal      = errors.New("Internal sqlite3 library error. This is a bug")
	ErrPermission    = errors.New("Permission denied")
	ErrAbort         = errors.New("Callback aborted")
	ErrBusy          = errors.New("Database file is already locked")
	ErrTableLocked   = errors.New("Table is already locked")
	ErrNoMem         = errors.New("Out of memory")
	ErrReadOnly      = errors.New("Database is read only")
	ErrInterrupt     = errors.New("Operation was interrupted")
	ErrIO            = errors.New("IO error")
	ErrCorrupt       = errors.New("Database is corrupted")
	ErrNotFound      = errors.New("Unknown opcode")
	ErrFull          = errors.New("Databse is full")
	ErrCantOpen      = errors.New("Unable to open database")
	ErrProtocol      = errors.New("Lock protocol error")
	ErrEmpty         = errors.New("Empty. This is a bug")
	ErrSchema        = errors.New("Invalid schema")
	ErrTooBig        = errors.New("BLOB is too big")
	ErrConstraint    = errors.New("Constraint violation")
	ErrMismatch      = errors.New("Mismatched data types")
	ErrMisuse        = errors.New("SQLite3 misuse. This is a bug")
	ErrOsUnsupported = errors.New("Feature not supported by OS")
	ErrAuth          = errors.New("Authorization failed")
	ErrFormat        = errors.New("Format error")
	ErrRange         = errors.New("Out of range")
	ErrNotDb         = errors.New("Not a valid database")
)

var errorCodes = []error{
	nil,
	ErrGeneric,
	ErrInternal,
	ErrPermission,
	ErrAbort,
	ErrBusy,
	ErrTableLocked,
	ErrNoMem,
	ErrReadOnly,
	ErrInterrupt,
	ErrIO,
	ErrCorrupt,
	ErrNotFound,
	ErrFull,
	ErrCantOpen,
	ErrProtocol,
	ErrEmpty,
	ErrSchema,
	ErrTooBig,
	ErrConstraint,
	ErrMismatch,
	ErrMisuse,
	ErrOsUnsupported,
	ErrAuth,
	ErrFormat,
	ErrRange,
	ErrNotDb,
}

func getErrorFromCode(result C.int) error {
	if result == 0 {
		return nil
	}

	if result < 0 || int(result) > len(errorCodes) {
		return fmt.Errorf("Unknown error: %d", result)
	}

	return errorCodes[result]
}
