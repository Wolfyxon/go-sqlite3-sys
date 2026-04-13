package sqlite

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

/*
#include <sqlite3.h>
#include <stdlib.h>
*/
import "C"

type sqliteStmt struct {
	handle      *C.sqlite3_stmt
	connection  *sqliteConn
	queryString string
}

func (s *sqliteStmt) Query(args []driver.Value) (driver.Rows, error) {
	bindErr := s.bindArgs(args)

	if bindErr != nil {
		return nil, bindErr
	}

	return &sqliteRows{
		stmt: s,
	}, nil
}

func (s *sqliteStmt) Exec(args []driver.Value) (driver.Result, error) {
	bindErr := s.bindArgs(args)

	if bindErr != nil {
		return nil, bindErr
	}

	var stepRes C.int

	for true {
		stepRes = C.sqlite3_step(s.handle)

		if stepRes == C.SQLITE_DONE {
			break
		}

		if stepRes == C.SQLITE_ERROR {
			return nil, fmt.Errorf("sqlite error")
		}

		if stepRes == C.SQLITE_BUSY {
			return nil, fmt.Errorf("unable to lock database")
		}
	}

	C.sqlite3_reset(s.handle)

	return &sqliteResult{
		stmt: s,
	}, nil
}

func (s *sqliteStmt) NumInput() int {
	return strings.Count(s.queryString, "?")
}

func (s *sqliteStmt) Close() error {
	res := C.sqlite3_finalize(s.handle)

	if res != 0 {
		return fmt.Errorf("Failed to finalize sqlite statement: %d", res)
	}

	return nil
}

func (s *sqliteStmt) bindArgs(args []driver.Value) error {
	for i, v := range args {
		cI := C.int(i + 1)
		valType := reflect.TypeOf(v).Kind().String()

		var bindRes any

		switch valType {
		case "int", "int32":
			bindRes = C.sqlite3_bind_int(s.handle, cI, C.int(v.(int)))
		case "int64":
			bindRes = C.sqlite3_bind_int64(s.handle, cI, C.sqlite3_int64(v.(int64)))
		case "string":
			cStr := C.CString(v.(string))
			bindRes = C.sqlite3_bind_text(s.handle, cI, cStr, -1, C.SQLITE_TRANSIENT)
			C.free(unsafe.Pointer(cStr))
		case "nil":
			bindRes = C.sqlite3_bind_null(s.handle, cI)
		case "slice":
			blob := v.([]byte)
			bindRes = C.sqlite3_bind_blob(s.handle, cI, unsafe.Pointer(&blob[0]), C.int(len(blob)), C.SQLITE_TRANSIENT)
		default:
			return fmt.Errorf("Unsupported type '%s'", valType)
		}

		if bindRes != C.int(0) {
			return fmt.Errorf("Failed to bind arg %d: %d", i, bindRes)
		}
	}

	return nil
}
