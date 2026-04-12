package sqlite

import (
	"database/sql/driver"
	"fmt"
)

/*
#include <sqlite3.h>
*/
import "C"

type sqliteStmt struct {
	handle     *C.sqlite3_stmt
	connection *sqliteConn
}

func (s *sqliteStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &sqliteRows{
		stmt: s,
	}, nil
}

func (s *sqliteStmt) Exec(args []driver.Value) (driver.Result, error) {
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
	return 0
}

func (s *sqliteStmt) Close() error {
	res := C.sqlite3_finalize(s.handle)

	if res != 0 {
		return fmt.Errorf("Failed to finalize sqlite statement: %d", res)
	}

	return nil
}
